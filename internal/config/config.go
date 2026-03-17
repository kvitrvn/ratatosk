package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBPath      string
	HTTPTimeout time.Duration
}

func Load() (Config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return Config{}, fmt.Errorf("user config dir: %w", err)
	}

	appDir := filepath.Join(configDir, "ratatosk")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return Config{}, fmt.Errorf("create config dir: %w", err)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(appDir)

	viper.SetDefault("db_path", filepath.Join(appDir, "ratatosk.db"))
	viper.SetDefault("http_timeout", "30s")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, fmt.Errorf("read config: %w", err)
		}
	}

	timeout, err := time.ParseDuration(viper.GetString("http_timeout"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid http_timeout: %w", err)
	}

	return Config{
		DBPath:      viper.GetString("db_path"),
		HTTPTimeout: timeout,
	}, nil
}
