package gpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

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

	promptBytes, err := os.ReadFile("internal/gpt/prompt.txt")
	if err != nil {
		return "", fmt.Errorf("❌ cannot read prompt.txt: %w", err)
	}

	promptTemplate := string(promptBytes)

	if !strings.Contains(promptTemplate, "{{context}}") {
		return "", fmt.Errorf("❌ prompt.txt missing {{context}} placeholder")
	}

	prompt := strings.ReplaceAll(promptTemplate, "{{context}}", context)

	req := GPTRequest{
		Model: "gpt-4o", // bạn có thể dùng gpt-3.5-turbo nếu muốn tiết kiệm chi phí
		Messages: []Message{
			{Role: "system", Content: "You generate API documentation in markdown."},
			{Role: "user", Content: prompt},
		},
	}

	reqBody, _ := json.Marshal(req)

	reqHTTP, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqBody))
	reqHTTP.Header.Set("Authorization", "Bearer "+apiKey)
	reqHTTP.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(reqHTTP)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("GPT API error: %s", string(body))
	}

	var gptResp GPTResponse
	err = json.NewDecoder(resp.Body).Decode(&gptResp)
	if err != nil {
		return "", err
	}

	return gptResp.Choices[0].Message.Content, nil
}
