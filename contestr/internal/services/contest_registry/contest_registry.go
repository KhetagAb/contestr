package contest_registry

import (
	"context"
	"fmt"

	"contestr/internal/repository"
)

type ContestRegistry interface {
	GetSystem(contestID int) (string, error)
	GetAllContests() map[string][]int
}

type MongoContestRegistry struct {
	repo repository.RegisteredContestRepository
}

func NewContestRegistry(repo repository.RegisteredContestRepository) ContestRegistry {
	return &MongoContestRegistry{repo: repo}
}

func (r *MongoContestRegistry) GetSystem(contestID int) (string, error) {
	contest, err := r.repo.GetByContestID(context.Background(), contestID)
	if err != nil {
		return "", err
	}
	if contest == nil {
		return "", fmt.Errorf("contest %d not found in registry", contestID)
	}
	return contest.System, nil
}

func (r *MongoContestRegistry) GetAllContests() map[string][]int {
	contests, err := r.repo.List(context.Background())
	if err != nil {
		return map[string][]int{}
	}

	result := make(map[string][]int)
	for _, c := range contests {
		result[c.System] = append(result[c.System], c.ContestID)
	}
	return result
}
