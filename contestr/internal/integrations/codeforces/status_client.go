package codeforces

import (
	"context"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

const (
	codeforcesAPIBaseURL     = "https://codeforces.com/api"
	codeforcesStatusPageSize = 1000
	codeforcesMinRequestGap  = 2 * time.Second
)

type StatusSubmission struct {
	ID                  int           `json:"id"`
	RelativeTimeSeconds int           `json:"relativeTimeSeconds"`
	Problem             StatusProblem `json:"problem"`
	Author              StatusParty   `json:"author"`
	Verdict             string        `json:"verdict"`
	Points              *float64      `json:"points,omitempty"`
}

type StatusProblem struct {
	Index string `json:"index"`
}

type StatusParty struct {
	Members []StatusMember `json:"members"`
}

type StatusMember struct {
	Handle string `json:"handle"`
}

type contestStatusResponse struct {
	Status  string             `json:"status"`
	Comment string             `json:"comment,omitempty"`
	Result  []StatusSubmission `json:"result"`
}

func (c *Service) GetContestStatus(ctx context.Context, contestID int) ([]StatusSubmission, error) {
	var result []StatusSubmission

	for from := 1; ; from += codeforcesStatusPageSize {
		page, err := c.getContestStatusPage(ctx, contestID, from, codeforcesStatusPageSize)
		if err != nil {
			return nil, err
		}
		result = append(result, page...)
		if len(page) < codeforcesStatusPageSize {
			break
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(codeforcesMinRequestGap):
		}
	}

	return result, nil
}

func (c *Service) getContestStatusPage(ctx context.Context, contestID int, from int, count int) ([]StatusSubmission, error) {
	if c.apiKey != "" && c.apiSecret != "" {
		managerPage, managerErr := c.getContestStatusPageWithMode(ctx, contestID, from, count, statusRequestManager)
		if managerErr == nil {
			return managerPage, nil
		}

		authPage, authErr := c.getContestStatusPageWithMode(ctx, contestID, from, count, statusRequestAuthenticated)
		if authErr == nil {
			return authPage, nil
		}

		publicPage, publicErr := c.getContestStatusPageWithMode(ctx, contestID, from, count, statusRequestPublic)
		if publicErr == nil {
			return publicPage, nil
		}

		return nil, fmt.Errorf(
			"codeforces contest.status failed: manager request failed: %v; authenticated request failed: %v; public request failed: %w",
			managerErr,
			authErr,
			publicErr,
		)
	}

	return c.getContestStatusPageWithMode(ctx, contestID, from, count, statusRequestPublic)
}

type statusRequestMode int

const (
	statusRequestPublic statusRequestMode = iota
	statusRequestAuthenticated
	statusRequestManager
)

func (c *Service) getContestStatusPageWithMode(ctx context.Context, contestID int, from int, count int, mode statusRequestMode) ([]StatusSubmission, error) {
	values := url.Values{}
	values.Set("contestId", strconv.Itoa(contestID))
	values.Set("from", strconv.Itoa(from))
	values.Set("count", strconv.Itoa(count))
	if mode != statusRequestPublic {
		values.Set("apiKey", c.apiKey)
		values.Set("time", strconv.FormatInt(time.Now().Unix(), 10))
		if mode == statusRequestManager {
			values.Set("asManager", "true")
		}
		values.Set("apiSig", generateCodeforcesAPISig("contest.status", c.apiSecret, values))
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		codeforcesAPIBaseURL+"/contest.status?"+values.Encode(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var body contestStatusResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK || body.Status != "OK" {
		if body.Comment != "" {
			return nil, fmt.Errorf("codeforces contest.status failed: %s", body.Comment)
		}
		return nil, fmt.Errorf("codeforces contest.status failed: %s", res.Status)
	}

	return body.Result, nil
}

func generateCodeforcesAPISig(method string, apiSecret string, values url.Values) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	prefix := ""
	for i := 0; i < 6; i++ {
		prefix += strconv.Itoa(r.Intn(10))
	}

	type param struct {
		name  string
		value string
	}
	params := make([]param, 0)
	for name, vals := range values {
		if name == "apiSig" {
			continue
		}
		for _, val := range vals {
			params = append(params, param{name: name, value: val})
		}
	}
	sort.Slice(params, func(i, j int) bool {
		if params[i].name != params[j].name {
			return params[i].name < params[j].name
		}
		return params[i].value < params[j].value
	})

	encoded := ""
	for i, p := range params {
		if i > 0 {
			encoded += "&"
		}
		encoded += p.name + "=" + p.value
	}

	hash := sha512.Sum512([]byte(prefix + "/" + method + "?" + encoded + "#" + apiSecret))
	return prefix + fmt.Sprintf("%x", hash)
}
