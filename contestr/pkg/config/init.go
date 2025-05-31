package config

import (
	"contestr/internal/configs"
	"contestr/pkg/logger"
	"context"
	"fmt"
	"os"
)

const defaultConfigPath = "configs/config.yaml"

func InitConfigAndContext() (*configs.Config, context.Context) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = defaultConfigPath
	}

	cfg, err := configs.LoadConfig(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load configuration: %v", err))
	}
	logger.Init(cfg.App.Name, "info")
	return cfg, context.Background()
}
