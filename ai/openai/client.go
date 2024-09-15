package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	URL       = "https://api.openai.com/v1/chat/completions"
	Model     = "gpt-4o-mini"
	MaxTokens = 1024
	System    = "You will receive a Git diff output. Based on the diff, generate a commit message. The message should include a short description on the first line, followed by a more detailed explanation of the changes made. Do not add comments or descriptions about the generated text."
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
	APIKey string
	System string
}

func New(key, system string) *Client {
	if key == "" {
		key = os.Getenv("OPENAI_API_KEY")
	}
	if system == "" {
		system = System
	}
	return &Client{APIKey: key, System: system}
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
	req.Header.Set("Authorization", "Bearer "+cli.APIKey)

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
	msgs := []Message{{Role: "system", Content: cli.System}, {Role: "user", Content: diff}}
	return cli.Chat(msgs)
}

func (cli *Client) RefineText(text string) (string, error) {
	// TODO needs refactoring.
	sys := "You are a writing assistant specialized in spelling and grammar correction. You will receive a Git commit message describing changes made to source code. Your task is to fix any spelling or grammatical errors while keeping changes minimal. Do not include explanations or comments about the corrections."
	msgs := []Message{{Role: "system", Content: sys}, {Role: "user", Content: text}}
	return cli.Chat(msgs)
}
