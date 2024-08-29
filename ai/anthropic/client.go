package anthropic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	URL          = "https://api.anthropic.com/v1/messages"
	Model        = "claude-3-5-sonnet-20240620"
	ModelVersion = "2023-06-01"
	MaxTokens    = 1024
	Role         = "user"
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
	APIKey string
	System string
}

func New(key, system string) *Client {
	if key == "" {
		key = os.Getenv("ANTHROPIC_API_KEY")
	}
	if system == "" {
		system = System
	}
	return &Client{APIKey: key, System: system}
}

func (cli *Client) CommitMessage(diff string) (string, error) {
	body := Body{
		Model:     Model,
		Messages:  []Message{{Role: Role, Content: diff}},
		MaxTokens: MaxTokens,
		System:    cli.System,
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
	req.Header.Set("x-api-key", cli.APIKey)

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
