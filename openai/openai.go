package openai

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	APIKey      string
	Instruction string
	Input       string
	Temperature float64
	Model       string
	Debug       bool
}

type OpenAIResponse struct {
	Choices []struct {
		Text    string `json:"text,omitempty"`
		Message struct {
			Role    string `json:"role,omitempty"`
			Content string `json:"content,omitempty"`
		} `json:"message,omitempty"`
	} `json:"choices"`
}

var ErrInvalidModel = errors.New("unsupported or invalid model")
var ErrNoChoices = errors.New("no choices returned from the API")
var ErrNoAssistantMessage = errors.New("no assistant message found in the API response")

func ValidateModel(model string) bool {
	acceptableModels := map[string]bool{
		"gpt-4":          true,
		"gpt-4-0314":     true,
		"gpt-4-32k":      true,
		"gpt-4-32k-0314": true,
		"gpt-3.5-turbo":  true,
	}

	_, exists := acceptableModels[model]
	return exists
}

func CallOpenAI(cfg Config, httpClient *http.Client) (string, error) {
	if !ValidateModel(cfg.Model) {
		return "", ErrInvalidModel
	}

	var jsonData []byte
	var err error

	messages := []map[string]string{
		{"role": "system", "content": cfg.Instruction},
		{"role": "user", "content": cfg.Input},
	}
	jsonData, err = json.Marshal(map[string]interface{}{
		"model":       cfg.Model,
		"messages":    messages,
		"temperature": cfg.Temperature,
		"max_tokens":  100,
		"stop":        []string{"\\n"},
	})

	if err != nil {
		return "", err
	}

	data := strings.NewReader(string(jsonData))
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", data)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)

	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var openAIResponse OpenAIResponse
	err = json.Unmarshal(body, &openAIResponse)
	if err != nil {
		return "", err
	}

	if len(openAIResponse.Choices) == 0 {
		return "", ErrNoChoices
	}

	assistantMessage := ""
	for _, choice := range openAIResponse.Choices {
		if choice.Message.Role == "assistant" {
			assistantMessage = strings.TrimSpace(choice.Message.Content)
			break
		}
		if choice.Text != "" {
			assistantMessage = strings.TrimSpace(choice.Text)
			break
		}
	}

	if assistantMessage == "" {
		return "", ErrNoAssistantMessage
	}

	return assistantMessage, nil
}
