package codeforces

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/togatoga/goforces"
)

type contestStandingsResponse struct {
	Status  string             `json:"status"`
	Comment string             `json:"comment,omitempty"`
	Result  goforces.Standings `json:"result"`
}

type standingsRequestMode int

const (
	standingsRequestPublic standingsRequestMode = iota
	standingsRequestAuthenticated
	standingsRequestManager
)

func (c *Service) GetContestStandings(ctx context.Context, contestID int) (*goforces.Standings, error) {
	if c.apiKey != "" && c.apiSecret != "" {
		managerStandings, managerErr := c.getContestStandings(ctx, contestID, standingsRequestManager)
		if managerErr == nil {
			return managerStandings, nil
		}

		authStandings, authErr := c.getContestStandings(ctx, contestID, standingsRequestAuthenticated)
		if authErr == nil {
			return authStandings, nil
		}

		publicStandings, publicErr := c.getContestStandings(ctx, contestID, standingsRequestPublic)
		if publicErr == nil {
			return publicStandings, nil
		}

		return nil, fmt.Errorf(
			"error getting contest standings: manager request failed: %v; authenticated request failed: %v; public request failed: %w",
			managerErr,
			authErr,
			publicErr,
		)
	}

	standings, err := c.getContestStandings(ctx, contestID, standingsRequestPublic)
	if err != nil {
		return nil, fmt.Errorf("error getting contest standings: %w", err)
	}
	return standings, nil
}

func (c *Service) getContestStandings(ctx context.Context, contestID int, mode standingsRequestMode) (*goforces.Standings, error) {
	values := url.Values{}
	values.Set("contestId", strconv.Itoa(contestID))

	if mode != standingsRequestPublic {
		values.Set("apiKey", c.apiKey)
		values.Set("time", strconv.FormatInt(time.Now().Unix(), 10))
		if mode == standingsRequestManager {
			values.Set("asManager", "true")
		}
		values.Set("apiSig", generateCodeforcesAPISig("contest.standings", c.apiSecret, values))
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		codeforcesAPIBaseURL+"/contest.standings?"+values.Encode(),
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

	var body contestStandingsResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK || body.Status != "OK" {
		if body.Comment != "" {
			return nil, fmt.Errorf("codeforces contest.standings failed: %s", body.Comment)
		}
		return nil, fmt.Errorf("codeforces contest.standings failed: %s", res.Status)
	}

	return &body.Result, nil
}
