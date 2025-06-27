package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ClaudeClient handles communication with Claude API
type ClaudeClient struct {
	APIKey     string
	BaseURL    string
	Model      string
	HTTPClient *http.Client
}

// ClaudeRequest represents the request structure for Claude API
type ClaudeRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []Message `json:"messages"`
	System    string    `json:"system,omitempty"`
}

// Message represents a message in the conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ClaudeResponse represents the response from Claude API
type ClaudeResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Model string `json:"model"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	StopReason string `json:"stop_reason"`
}

// ClaudeError represents an error response from Claude API
type ClaudeError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (e ClaudeError) Error() string {
	return fmt.Sprintf("Claude API error (%s): %s", e.Type, e.Message)
}

// NewClaudeClient creates a new Claude API client
func NewClaudeClient(apiKey, baseURL, model string) *ClaudeClient {
	return &ClaudeClient{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Model:   model,
		HTTPClient: &http.Client{
			Timeout: 120 * time.Second, // 2 minutes timeout for longer requests
		},
	}
}

// ProcessContent sends content to Claude API with the given system prompt
func (c *ClaudeClient) ProcessContent(content, systemPrompt string) (*ClaudeResponse, error) {
	if c.APIKey == "" {
		return nil, fmt.Errorf("API key is required. Set ANTHROPIC_API_KEY environment variable or add it to config.json")
	}

	request := ClaudeRequest{
		Model:     c.Model,
		MaxTokens: 4096, // Adjust based on your needs
		System:    systemPrompt,
		Messages: []Message{
			{
				Role:    "user",
				Content: content,
			},
		},
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/v1/messages", c.BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	// Make the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle non-200 status codes
	if resp.StatusCode != http.StatusOK {
		var claudeErr ClaudeError
		if err := json.Unmarshal(body, &claudeErr); err != nil {
			return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
		}
		return nil, claudeErr
	}

	// Parse successful response
	var response ClaudeResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// ExtractTextFromResponse extracts the text content from Claude's response
func (c *ClaudeClient) ExtractTextFromResponse(response *ClaudeResponse) string {
	if len(response.Content) == 0 {
		return ""
	}

	var result string
	for _, content := range response.Content {
		if content.Type == "text" {
			result += content.Text
		}
	}
	return result
}

// GetTokenUsage returns the total tokens used in the response
func (c *ClaudeClient) GetTokenUsage(response *ClaudeResponse) int {
	return response.Usage.InputTokens + response.Usage.OutputTokens
}