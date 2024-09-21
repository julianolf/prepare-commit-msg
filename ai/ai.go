package ai

import (
	"github.com/julianolf/prepare-commit-msg/ai/anthropic"
	"github.com/julianolf/prepare-commit-msg/ai/config"
	"github.com/julianolf/prepare-commit-msg/ai/openai"
)

type AI interface {
	CommitMessage(string) (string, error)
	RefineText(string) (string, error)
}

func New(cfg *config.Config) AI {
	switch cfg.AI {
	case "openai":
		return openai.New(cfg)
	default:
		return anthropic.New(cfg)
	}
}
