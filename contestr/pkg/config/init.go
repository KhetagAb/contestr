package config

import (
	"contestr/internal/configs"
	"contestr/pkg/logger"
	"context"
	"fmt"
)

const configPath = "configs/config.yaml"

func InitConfigAndContext() (*configs.Config, context.Context) {
	cfg, err := configs.LoadConfig(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load configuration: %v", err))
	}
	logger.Init(cfg.App.Name, "info")
	return cfg, context.Background()
}
