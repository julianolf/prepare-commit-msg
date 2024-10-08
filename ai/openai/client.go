package openai

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
	URL       = "https://api.openai.com/v1/chat/completions"
	Model     = "gpt-4o-mini"
	MaxTokens = 1024
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Body struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type Choice struct {
	Index        int          `json:"index"`
	Message      Message      `json:"message"`
	Logprobs     *interface{} `json:"logprobs"`
	FinishReason string       `json:"finish_reason"`
}

type Response struct {
	Id                string         `json:"id"`
	Object            string         `json:"object"`
	Created           int            `json:"created"`
	Model             string         `json:"model"`
	SystemFingerprint string         `json:"system_fingerprint"`
	Choices           []Choice       `json:"choices"`
	Usage             map[string]any `json:"usage"`
}

type Client struct {
	http.Client
	Config *config.Config
}

func New(cfg *config.Config) *Client {
	if cfg.APIKey == "" {
		cfg.APIKey = os.Getenv("OPENAI_API_KEY")
	}
	return &Client{Config: cfg}
}

func (cli *Client) Chat(messages []Message) (string, error) {
	body := Body{
		Model:     Model,
		Messages:  messages,
		MaxTokens: MaxTokens,
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

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cli.Config.APIKey)

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

	return response.Choices[0].Message.Content, nil
}

func (cli *Client) CommitMessage(diff string) (string, error) {
	msgs := []Message{{Role: "system", Content: cli.Config.System.GenMsg}, {Role: "user", Content: diff}}
	return cli.Chat(msgs)
}

func (cli *Client) RefineText(text string) (string, error) {
	msgs := []Message{{Role: "system", Content: cli.Config.System.FixMsg}, {Role: "user", Content: text}}
	return cli.Chat(msgs)
}
