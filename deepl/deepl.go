package deepl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	baseURL = "https://api-free.deepl.com/v2"
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

type Translation struct {
	DetectedSourceLanguage string `json:"detected_source_language"`
	Text                   string `json:"text"`
}

type TranslationResponse struct {
	Translations []Translation `json:"translations"`
}

func (c *Client) Translate(text []string, targetLang string) (string, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"text":        text,
		"target_lang": targetLang,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	request, err := http.NewRequest("POST", baseURL+"/translate", bytes.NewReader(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	q := request.URL.Query()
	q.Add("auth_key", c.apiKey)
	request.URL.RawQuery = q.Encode()

	request.Header.Set("Content-Type", "application/json")

	response, err := c.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read HTTP response: %w", err)
	}

	var responseBody TranslationResponse
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return "", fmt.Errorf("failed to decode response body: %w", err)
	}

	if len(responseBody.Translations) == 0 {
		return "", errors.New("DeepL APIから回答がありませんでした")
	}

	return responseBody.Translations[0].Text, nil
}
