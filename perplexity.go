package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// PerplexityClient handles communication with Perplexity API
type PerplexityClient struct {
	APIKey     string
	BaseURL    string
	Model      string
	MaxTokens  int
	HTTPClient *http.Client
}

// PerplexityRequest represents the request structure for Perplexity API
type PerplexityRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	Stream      bool      `json:"stream"`
}

// PerplexityResponse represents the response from Perplexity API
type PerplexityResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Citations json.RawMessage `json:"citations,omitempty"`
}

// PerplexityError represents an error response from Perplexity API
type PerplexityError struct {
	ErrorInfo struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

func (e PerplexityError) Error() string {
	return fmt.Sprintf("Perplexity API error (%s): %s", e.ErrorInfo.Type, e.ErrorInfo.Message)
}

// NewPerplexityClient creates a new Perplexity API client
func NewPerplexityClient(apiKey, baseURL, model string, maxTokens int) *PerplexityClient {
	if baseURL == "" {
		baseURL = "https://api.perplexity.ai"
	}
	if model == "" {
		model = "llama-3.1-sonar-large-128k-online"
	}
	if maxTokens == 0 {
		maxTokens = 8192 // Default value for backward compatibility
	}

	return &PerplexityClient{
		APIKey:    apiKey,
		BaseURL:   baseURL,
		Model:     model,
		MaxTokens: maxTokens,
		HTTPClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// ResearchContent performs research using Perplexity API
func (c *PerplexityClient) ResearchContent(query, systemPrompt string) (*PerplexityResponse, error) {
	if c.APIKey == "" {
		return nil, fmt.Errorf("Perplexity API key is required. Set PERPLEXITY_API_KEY environment variable or add it to config.json")
	}

	messages := []Message{}
	
	// Add system prompt if provided
	if systemPrompt != "" {
		messages = append(messages, Message{
			Role:    "system",
			Content: systemPrompt,
		})
	}
	
	// Add user query
	messages = append(messages, Message{
		Role:    "user",
		Content: query,
	})

	request := PerplexityRequest{
		Model:       c.Model,
		Messages:    messages,
		MaxTokens:   c.MaxTokens, // Use configurable max tokens
		Temperature: 0.2, // Lower temperature for more factual research
		Stream:      false,
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/chat/completions", c.BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

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
		var perplexityErr PerplexityError
		if err := json.Unmarshal(body, &perplexityErr); err != nil {
			return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
		}
		return nil, perplexityErr
	}

	// Parse successful response
	var response PerplexityResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// ExtractTextFromResponse extracts the text content from Perplexity's response
func (c *PerplexityClient) ExtractTextFromResponse(response *PerplexityResponse) string {
	if len(response.Choices) == 0 {
		return ""
	}
	return response.Choices[0].Message.Content
}

// Citation represents a single citation from Perplexity
type Citation struct {
	Number int    `json:"number"`
	URL    string `json:"url"`
	Title  string `json:"title"`
}

// ExtractCitationsFromResponse extracts citations from Perplexity's response
func (c *PerplexityClient) ExtractCitationsFromResponse(response *PerplexityResponse) []Citation {
	if len(response.Citations) == 0 {
		return []Citation{}
	}
	
	// Try to parse as array of citation objects first
	var citationsArray []struct {
		Number int    `json:"number"`
		URL    string `json:"url"`
		Title  string `json:"title"`
	}
	
	if err := json.Unmarshal(response.Citations, &citationsArray); err == nil {
		citations := make([]Citation, len(citationsArray))
		for i, citation := range citationsArray {
			citations[i] = Citation{
				Number: citation.Number,
				URL:    citation.URL,
				Title:  citation.Title,
			}
		}
		return citations
	}
	
	// If that fails, try parsing as a string (maybe it's a stringified JSON)
	var citationsString string
	if err := json.Unmarshal(response.Citations, &citationsString); err == nil {
		// Try to parse the string as JSON
		if err := json.Unmarshal([]byte(citationsString), &citationsArray); err == nil {
			citations := make([]Citation, len(citationsArray))
			for i, citation := range citationsArray {
				citations[i] = Citation{
					Number: citation.Number,
					URL:    citation.URL,
					Title:  citation.Title,
				}
			}
			return citations
		}
	}
	
	// If all parsing fails, return empty slice
	return []Citation{}
}

// GetTokenUsage returns the total tokens used in the response
func (c *PerplexityClient) GetTokenUsage(response *PerplexityResponse) int {
	return response.Usage.TotalTokens
}

// FormatResearchOutput formats the research output with citations
func (c *PerplexityClient) FormatResearchOutput(response *PerplexityResponse) string {
	content := c.ExtractTextFromResponse(response)
	citations := c.ExtractCitationsFromResponse(response)
	
	if len(citations) == 0 {
		return content
	}
	
	// Add citations section
	content += "\n\n## Sources\n\n"
	for _, citation := range citations {
		content += fmt.Sprintf("%d. [%s](%s)\n", citation.Number, citation.Title, citation.URL)
	}
	
	return content
}
