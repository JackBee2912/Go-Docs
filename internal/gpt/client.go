package gpt

import (
	"bytes"
	_ "embed" // dùng để embed file vào binary
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

//go:embed prompt.txt
var promptTemplate string

type GPTRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GPTResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func GenerateMarkdownDocumentation(context, apiKey string) (string, error) {
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY is not set")
	}

	// Check nếu promptTemplate không chứa {{context}}
	if !strings.Contains(promptTemplate, "{{context}}") {
		return "", fmt.Errorf("❌ embedded prompt missing {{context}} placeholder")
	}

	// Thay thế {{context}} trong prompt
	prompt := strings.ReplaceAll(promptTemplate, "{{context}}", context)

	req := GPTRequest{
		Model: "gpt-4o", // hoặc bạn có thể dùng "gpt-3.5-turbo" để tiết kiệm chi phí
		Messages: []Message{
			{Role: "system", Content: "You generate API documentation in markdown."},
			{Role: "user", Content: prompt},
		},
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	reqHTTP, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	reqHTTP.Header.Set("Authorization", "Bearer "+apiKey)
	reqHTTP.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(reqHTTP)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("GPT API error: %s", string(body))
	}

	var gptResp GPTResponse
	err = json.NewDecoder(resp.Body).Decode(&gptResp)
	if err != nil {
		return "", fmt.Errorf("failed to decode GPT response: %w", err)
	}

	if len(gptResp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from GPT")
	}

	return gptResp.Choices[0].Message.Content, nil
}
