package configs

import (
	"strings"
	"time"

	"github.com/joho/godotenv"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App           AppConfig             `mapstructure:"app"`
		HTTP          HTTPConfig            `mapstructure:"http"`
		MongoDB       MongoDBConfig         `mapstructure:"mongodb"`
		Telegram      TelegramConfig        `mapstructure:"telegram"`
		Codeforces    CodeforcesConfig      `mapstructure:"codeforces"`
		Ejudge        EjudgeConfig          `mapstructure:"ejudge"`
		ContestSync   ContestSyncConfig   `mapstructure:"contest_sync"`
		TimetableSync TimetableSyncConfig `mapstructure:"timetable_sync"`
		Admin         AdminConfig         `mapstructure:"admin"`
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
		MaxPoolSize uint64 `mapstructure:"max_pool_size"`
		MinPoolSize uint64 `mapstructure:"min_pool_size"`

		Database                string `mapstructure:"database"`
		TourCollection          string `mapstructure:"tour_collection"`
		TourTimetableCollection string `mapstructure:"tour_timetable_collection"`
	}

	TelegramConfig struct {
		Token  string `mapstructure:"token"`
		Admins []int  `mapstructure:"admins_tg_id"`
	}

	CodeforcesConfig struct {
		APIKey    string `mapstructure:"api_key"`
		APISecret string `mapstructure:"api_secret"`
	}

	EjudgeConfig struct {
		XMLUrl         string        `mapstructure:"xml_url"`
		RequestTimeout time.Duration `mapstructure:"request_timeout"`
	}

	ContestSyncConfig struct {
		Interval time.Duration `mapstructure:"interval"`
	}

	TimetableSyncConfig struct {
		Interval time.Duration `mapstructure:"interval"`
	}

	AdminConfig struct {
		JWTSecret string        `mapstructure:"jwt_secret"`
		Username  string        `mapstructure:"username"`
		Password  string        `mapstructure:"password"`
		JWTTTL    time.Duration `mapstructure:"jwt_ttl"`
	}
)

func (c AdminConfig) Enabled() bool {
	return strings.TrimSpace(c.JWTSecret) != "" &&
		strings.TrimSpace(c.Username) != "" &&
		strings.TrimSpace(c.Password) != ""
}

func (c TelegramConfig) Enabled() bool {
	token := strings.TrimSpace(c.Token)
	if token == "" || token == "YOUR_TELEGRAM_BOT_TOKEN" || token == "YOUR_BOT_API_KEY" {
		return false
	}
	return true
}

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
