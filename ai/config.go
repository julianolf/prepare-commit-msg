package ai

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

func (c *Config) Update(config *Config) {
	if config.AI != "" {
		c.AI = config.AI
	}
	if config.APIKey != "" {
		c.APIKey = config.APIKey
	}
	if config.System != "" {
		c.System = config.System
	}
}

func ConfigFromFile(filename string) (*Config, error) {
	config := &Config{}
	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func ConfigFromEnv() *Config {
	config := &Config{
		AI:     os.Getenv("PREPARE_COMMIT_MSG_AI"),
		APIKey: os.Getenv("PREPARE_COMMIT_MSG_APIKEY"),
		System: os.Getenv("PREPARE_COMMIT_MSG_SYSTEM"),
	}
	return config
}
