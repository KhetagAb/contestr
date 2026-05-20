package integrations

import (
	"context"

	"contestr/pkg/regatta"
)

type FetchContestOptions struct {
	ScoringSettings regatta.ScoringSettings
	TourSettings    regatta.TourSettings
}

type ContestAdapter interface {
	FetchContest(ctx context.Context, contestID int, opts FetchContestOptions) (*regatta.Contest, error)
	GetSystem() string // "ejudge" or "codeforces"
}
