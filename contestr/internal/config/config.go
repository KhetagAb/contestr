package config

import (
	"github.com/joho/godotenv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App        AppConfig        `mapstructure:"app"`
		HTTP       HTTPConfig       `mapstructure:"http"`
		MongoDB    MongoDBConfig    `mapstructure:"mongodb"`
		Telegram   TelegramConfig   `mapstructure:"telegram"`
		Codeforces CodeforcesConfig `mapstructure:"codeforces"`
	}

	AppConfig struct {
		Name    string `mapstructure:"name"`
		Version string `mapstructure:"version"`
	}

	HTTPConfig struct {
		Port            string        `mapstructure:"port"`
		ReadTimeout     time.Duration `mapstructure:"read_timeout"`
		WriteTimeout    time.Duration `mapstructure:"write_timeout"`
		ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	}

	MongoDBConfig struct {
		URI         string `mapstructure:"uri"`
		Database    string `mapstructure:"database"`
		MaxPoolSize uint64 `mapstructure:"max_pool_size"`
		MinPoolSize uint64 `mapstructure:"min_pool_size"`
	}

	TelegramConfig struct {
		Token string `mapstructure:"token"`
	}

	CodeforcesConfig struct {
		APIKey    string `mapstructure:"api_key"`
		APISecret string `mapstructure:"api_secret"`
	}
)

func LoadConfig(path string) (*Config, error) {
	_ = godotenv.Load()

	viper.SetConfigFile(path)

	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
