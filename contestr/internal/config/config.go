package config

import (
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App        AppConfig
		HTTP       HTTPConfig
		MongoDB    MongoDBConfig
		Telegram   TelegramConfig
		Codeforces CodeforcesConfig
	}

	AppConfig struct {
		Name    string
		Version string
	}

	HTTPConfig struct {
		Port            string
		ReadTimeout     time.Duration
		WriteTimeout    time.Duration
		ShutdownTimeout time.Duration
	}

	MongoDBConfig struct {
		URI         string
		Database    string
		MaxPoolSize uint64
		MinPoolSize uint64
		MaxIdleTime time.Duration
	}

	TelegramConfig struct {
		Token string
	}

	CodeforcesConfig struct {
		APIKey    string
		APISecret string
	}
)

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
