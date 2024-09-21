package ai

import (
	"github.com/julianolf/prepare-commit-msg/ai/anthropic"
	"github.com/julianolf/prepare-commit-msg/ai/openai"
)

type AI interface {
	CommitMessage(string) (string, error)
	RefineText(string) (string, error)
}

func New(config *Config) AI {
	switch config.AI {
	case "openai":
		return openai.New(config.APIKey, config.System)
	default:
		return anthropic.New(config.APIKey, config.System)
	}
}
