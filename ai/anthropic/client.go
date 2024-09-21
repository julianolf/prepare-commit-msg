package anthropic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/julianolf/prepare-commit-msg/ai/config"
)

const (
	URL          = "https://api.anthropic.com/v1/messages"
	Model        = "claude-3-5-sonnet-20240620"
	ModelVersion = "2023-06-01"
	MaxTokens    = 1024
	System       = "You will receive a Git diff output. Based on the diff, generate a commit message. The message should include a short description on the first line, followed by a more detailed explanation of the changes made. Do not add comments or descriptions about the generated text."
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Body struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
	Stream    bool      `json:"stream"`
	System    string    `json:"system"`
}

type Response struct {
	Id           string              `json:"id"`
	Type         string              `json:"type"`
	Role         string              `json:"role"`
	Content      []map[string]string `json:"content"`
	Model        string              `json:"model"`
	StopReason   *string             `json:"stop_reason"`
	StopSequence *string             `json:"stop_sequence"`
	Usage        map[string]*int     `json:"usage"`
}

type Client struct {
	http.Client
	Config *config.Config
}

func New(cfg *config.Config) *Client {
	if cfg.APIKey == "" {
		cfg.APIKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	if cfg.System == "" {
		cfg.System = System
	}
	return &Client{Config: cfg}
}

func (cli *Client) Chat(messages []Message, system string) (string, error) {
	body := Body{
		Model:     Model,
		Messages:  messages,
		MaxTokens: MaxTokens,
		System:    system,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return "", nil
	}

	reader := bytes.NewReader(data)
	req, err := http.NewRequest(http.MethodPost, URL, reader)
	if err != nil {
		return "", nil
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Set("anthropic-version", ModelVersion)
	req.Header.Set("x-api-key", cli.Config.APIKey)

	res, err := cli.Do(req)
	if err != nil {
		return "", nil
	}

	data, err = io.ReadAll(res.Body)
	if err != nil {
		return "", nil
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Request failed [%d]: %s\n", res.StatusCode, string(data))
	}

	response := new(Response)
	err = json.Unmarshal(data, response)
	if err != nil {
		return "", err
	}

	return response.Content[0]["text"], nil
}

func (cli *Client) CommitMessage(diff string) (string, error) {
	msgs := []Message{{Role: "user", Content: diff}}
	return cli.Chat(msgs, cli.Config.System)
}

func (cli *Client) RefineText(text string) (string, error) {
	// TODO needs refactoring.
	sys := "You are a writing assistant specialized in spelling and grammar correction. You will receive a Git commit message describing changes made to source code. Your task is to fix any spelling or grammatical errors while keeping changes minimal. Do not include explanations or comments about the corrections."
	msgs := []Message{{Role: "user", Content: text}}
	return cli.Chat(msgs, sys)
}
