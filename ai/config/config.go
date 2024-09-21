package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var DefaultConfigFile string

func init() {
	dir, err := os.UserConfigDir()
	if err == nil {
		DefaultConfigFile = filepath.Join(dir, "prepare-commit-msg", "config.json")
	}
}

type Config struct {
	AI     string
	APIKey string
	System string
}

func (c *Config) Update(cfg *Config) {
	if cfg.AI != "" {
		c.AI = cfg.AI
	}
	if cfg.APIKey != "" {
		c.APIKey = cfg.APIKey
	}
	if cfg.System != "" {
		c.System = cfg.System
	}
}

func ConfigFromFile(filename string) (*Config, error) {
	cfg := &Config{}
	data, err := os.ReadFile(filename)
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(data, cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func ConfigFromEnv() *Config {
	cfg := &Config{
		AI:     os.Getenv("PREPARE_COMMIT_MSG_AI"),
		APIKey: os.Getenv("PREPARE_COMMIT_MSG_APIKEY"),
		System: os.Getenv("PREPARE_COMMIT_MSG_SYSTEM"),
	}
	return cfg
}
