package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	baseURL     = "https://api.openai.com/v1"
	endpoint    = "/chat/completions"
	model       = "gpt-3.5-turbo"
	role        = "user"
	temperature = 0.9
	max_tokens  = 1000
)

type Client struct {
	apiKey string
	client *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

// レスポンス
type Response struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []Choices `json:"choices"`
}

type Choices struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
}

// リクエスト
type Request struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (c *Client) GetChatResponse(prompt string) ([]string, error) {
	requestBody := Request{
		Model:     model,
		MaxTokens: max_tokens,
		Messages: []Message{{
			Role:    role,
			Content: prompt,
		}},
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	request, err := http.NewRequest("POST", baseURL+endpoint, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create API request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.apiKey)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}

	defer response.Body.Close()

	var responseBody Response
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if len(responseBody.Choices) == 0 {
		return nil, errors.New("ChatGPT API returned empty response")
	}

	var res []string
	for _, r := range responseBody.Choices {
		res = append(res, r.Message.Content)
	}

	return res, nil
}
