package codeforces

import (
	"context"
	"fmt"
	"github.com/togatoga/goforces"
)

func (c *Service) GetContest(ctx context.Context, contestID int) (*goforces.Standings, error) {
	standings, err := c.client.GetContestStandings(ctx, contestID, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting contest standings: %w", err)
	}
	return standings, nil
}
