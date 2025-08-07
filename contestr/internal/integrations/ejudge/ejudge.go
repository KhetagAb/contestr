package ejudge

import (
	"contestr/internal/configs"
	"contestr/pkg/logger"
	"contestr/pkg/regatta"
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ContestXMLFetcher struct {
	cfg *configs.Config
}

func NewContestXMLFetcher(cfg *configs.Config) *ContestXMLFetcher {
	return &ContestXMLFetcher{
		cfg: cfg,
	}
}

func (s *ContestXMLFetcher) FetchAndParseXML(ctx context.Context, contestId int) (*regatta.RunLog, error) {
	client := http.Client{Timeout: s.cfg.Ejudge.RequestTimeout}

	url := fmt.Sprintf("%s/%v.xml", s.cfg.Ejudge.XMLUrl, contestId)
	logger.Infof(ctx, "fetching ejudge standings in url: %v", url)
	resp, err := client.Get(url)
	if err != nil {
		logger.Errorf(ctx, "error during get request: %v", err)
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Errorf(ctx, "not 200 status: %v", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return s.parseXMLToRunLog(ctx, err, bodyBytes)
}

func (s *ContestXMLFetcher) parseXMLToRunLog(ctx context.Context, err error, bodyBytes []byte) (*regatta.RunLog, error) {
	var runLog regatta.RunLog
	err = xml.Unmarshal(bodyBytes, &runLog)
	if err != nil {
		return nil, fmt.Errorf("error parsing XML: %v", err)
	}
	logger.Infof(ctx, "success fetched and parsed ejudge contest")
	return &runLog, nil
}
