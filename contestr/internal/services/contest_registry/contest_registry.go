package contest_registry

import (
	"contestr/internal/configs"
	"fmt"
)

type ContestRegistry interface {
	GetSystem(contestID int) (string, error)
	GetAllContests() map[string][]int
}

type MongoContestRegistry struct {
	registry map[string][]int
}

func NewContestRegistry(cfg *configs.Config) ContestRegistry {
	registry := make(map[string][]int)
	if cfg.Contests.Registry != nil {
		registry = cfg.Contests.Registry
	}
	return &MongoContestRegistry{
		registry: registry,
	}
}

func (r *MongoContestRegistry) GetSystem(contestID int) (string, error) {
	for system, contestIDs := range r.registry {
		for _, id := range contestIDs {
			if id == contestID {
				return system, nil
			}
		}
	}
	return "", fmt.Errorf("contest %d not found in registry", contestID)
}

func (r *MongoContestRegistry) GetAllContests() map[string][]int {
	return r.registry
}
