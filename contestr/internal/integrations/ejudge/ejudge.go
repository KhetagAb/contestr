package ejudge

import (
	"contestr/internal/configs"
	"contestr/pkg/regatta"
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

func (s *ContestXMLFetcher) FetchAndParseXML(contestId int) (*regatta.RunLog, error) {
	client := http.Client{Timeout: s.cfg.Ejudge.RequestTimeout}

	url := fmt.Sprintf("%s/%s.xml", s.cfg.Ejudge.XMLUrl, contestId)
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return s.parseXMLToRunLog(err, bodyBytes)
}

func (s *ContestXMLFetcher) parseXMLToRunLog(err error, bodyBytes []byte) (*regatta.RunLog, error) {
	var runLog regatta.RunLog
	err = xml.Unmarshal(bodyBytes, &runLog)
	if err != nil {
		return nil, fmt.Errorf("error parsing XML: %v", err)
	}
	return &runLog, nil
}
