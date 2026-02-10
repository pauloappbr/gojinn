package gojinn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AIRequest struct {
	Model    string      `json:"model"`
	Messages []AIMessage `json:"messages"`
	Stream   bool        `json:"stream"`
}

type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AIResponse struct {
	Choices []struct {
		Message AIMessage `json:"message"`
	} `json:"choices"`
}

func (g *Gojinn) askAI(prompt string) (string, error) {
	provider := g.AIProvider
	if provider == "" {
		provider = "openai"
	}
	model := g.AIModel
	if model == "" {
		model = "gpt-3.5-turbo"
	}

	cacheKey := fmt.Sprintf("%s:%s", model, hashString(prompt))
	if cachedVal, ok := g.aiCache.Load(cacheKey); ok {
		return cachedVal.(string), nil
	}

	endpoint := g.AIEndpoint
	if endpoint == "" {
		if provider == "ollama" {
			endpoint = "http://localhost:11434/v1/chat/completions"
		} else {
			endpoint = "https://api.openai.com/v1/chat/completions"
		}
	}

	if len(g.AllowedHosts) > 0 {
		u, err := url.Parse(endpoint)
		if err == nil {
			allowed := false
			hostname := u.Hostname()
			if provider == "ollama" && (hostname == "localhost" || hostname == "127.0.0.1") {
				allowed = true
			} else {
				for _, host := range g.AllowedHosts {
					if strings.Contains(hostname, host) {
						allowed = true
						break
					}
				}
			}
			if !allowed {
				return "", fmt.Errorf("egress denied to %s", hostname)
			}
		}
	}

	reqBody := AIRequest{
		Model:  model,
		Stream: false,
		Messages: []AIMessage{
			{Role: "system", Content: "You are a helpful assistant running inside Gojinn Serverless."},
			{Role: "user", Content: prompt},
		},
	}
	jsonData, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if g.AIToken != "" {
		req.Header.Set("Authorization", "Bearer "+g.AIToken)
	}

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("AI connect error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("AI API error (%d): %s", resp.StatusCode, string(body))
	}

	var aiResp AIResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		return "", fmt.Errorf("json decode error: %w", err)
	}

	if len(aiResp.Choices) > 0 {
		responseContent := aiResp.Choices[0].Message.Content
		g.aiCache.Store(cacheKey, responseContent)
		return responseContent, nil
	}
	return "", fmt.Errorf("AI returned no response")
}
