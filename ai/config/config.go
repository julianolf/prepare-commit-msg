package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	defaultAI    = "anthropic"
	genMsgPrompt = "You will receive a Git diff output. Based on the diff, generate a commit message. The message should include a short description on the first line, followed by a more detailed explanation of the changes made. Do not add comments or descriptions about the generated text."
	fixMsgPrompt = "You are a writing assistant specialized in spelling and grammar correction. You will receive a Git commit message describing changes made to source code. Your task is to fix any spelling or grammatical errors while keeping changes minimal. Do not include explanations or comments about the corrections."
)

var DefaultConfigFile string

func init() {
	dir, err := os.UserConfigDir()
	if err == nil {
		DefaultConfigFile = filepath.Join(dir, "prepare-commit-msg", "config.json")
	}
}

type SystemPrompt struct {
	GenMsg string
	FixMsg string
}

type Config struct {
	AI     string
	APIKey string
	System *SystemPrompt
}

func (c *Config) Update(cfg *Config) {
	if cfg.AI != "" {
		c.AI = cfg.AI
	}
	if cfg.APIKey != "" {
		c.APIKey = cfg.APIKey
	}
	if cfg.System != nil {
		if cfg.System.GenMsg != "" {
			c.System.GenMsg = cfg.System.GenMsg
		}
		if cfg.System.FixMsg != "" {
			c.System.FixMsg = cfg.System.FixMsg
		}
	}
}

func Default() *Config {
	return &Config{
		AI: defaultAI,
		System: &SystemPrompt{
			GenMsg: genMsgPrompt,
			FixMsg: fixMsgPrompt,
		},
	}
}

func FromFile(filename string) (*Config, error) {
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

func FromEnv() *Config {
	cfg := &Config{
		AI:     os.Getenv("PREPARE_COMMIT_MSG_AI"),
		APIKey: os.Getenv("PREPARE_COMMIT_MSG_APIKEY"),
		System: &SystemPrompt{
			GenMsg: os.Getenv("PREPARE_COMMIT_MSG_SYSTEM_GENMSG"),
			FixMsg: os.Getenv("PREPARE_COMMIT_MSG_SYSTEM_FIXMSG"),
		},
	}
	return cfg
}
