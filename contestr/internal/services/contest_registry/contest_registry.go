package contest_registry

import (
	"context"
	"fmt"

	"contestr/internal/repository"
)

type ContestRegistry interface {
	GetSystem(contestID int) (string, error)
	GetAllContests() map[string][]int
	GetContest(contestID int) (*repository.RegisteredContest, error)
	GetAllRegisteredContests() []repository.RegisteredContest
}

type MongoContestRegistry struct {
	repo repository.RegisteredContestRepository
}

func NewContestRegistry(repo repository.RegisteredContestRepository) ContestRegistry {
	return &MongoContestRegistry{repo: repo}
}

func (r *MongoContestRegistry) GetSystem(contestID int) (string, error) {
	contest, err := r.GetContest(contestID)
	if err != nil {
		return "", err
	}
	return contest.System, nil
}

func (r *MongoContestRegistry) GetContest(contestID int) (*repository.RegisteredContest, error) {
	contest, err := r.repo.GetByContestID(context.Background(), contestID)
	if err != nil {
		return nil, err
	}
	if contest == nil {
		return nil, fmt.Errorf("contest %d not found in registry", contestID)
	}
	return contest, nil
}

func (r *MongoContestRegistry) GetAllContests() map[string][]int {
	result := make(map[string][]int)
	for _, c := range r.GetAllRegisteredContests() {
		result[c.System] = append(result[c.System], c.ContestID)
	}
	return result
}

func (r *MongoContestRegistry) GetAllRegisteredContests() []repository.RegisteredContest {
	contests, err := r.repo.List(context.Background())
	if err != nil {
		return []repository.RegisteredContest{}
	}
	return contests
}
