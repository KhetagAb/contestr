package integrations

import (
	"contestr/pkg/regatta"
	"context"
)

type ContestAdapter interface {
	FetchContest(ctx context.Context, contestID int) (*regatta.Contest, error)
	GetSystem() string // "ejudge" или "codeforces"
}
