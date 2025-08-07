package codeforces

import (
	"contestr/internal/configs"
	"github.com/togatoga/goforces"
)

type Service struct {
	client *goforces.Client
}

func NewService(cfg *configs.Config) *Service {
	client, _ := goforces.NewClient(nil)
	client.SetAPISecret(cfg.Codeforces.APISecret)
	client.SetAPIKey(cfg.Codeforces.APIKey)
	return &Service{
		client: client,
	}
}
