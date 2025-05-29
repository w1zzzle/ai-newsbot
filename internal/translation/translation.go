package translation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Translator handles AI-powered translation using OpenRouter
type Translator struct {
	apiKey     string
	httpClient *http.Client
}

// OpenRouterRequest represents the request structure for OpenRouter API
// This matches the OpenAI SDK structure used in the Python example
type OpenRouterRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream,omitempty"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenRouterResponse represents the response from OpenRouter API
type OpenRouterResponse struct {
	Choices []Choice  `json:"choices"`
	Error   *APIError `json:"error,omitempty"`
}

// Choice represents a response choice
type Choice struct {
	Message Message `json:"message"`
}

// APIError represents an API error
type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

// New creates a new Translator instance
func New(apiKey string) *Translator {
	return &Translator{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 60 * time.Second, // DeepSeek R1 can be slower due to reasoning
		},
	}
}

// TranslateToRussian translates text to Russian using DeepSeek R1
func (t *Translator) TranslateToRussian(ctx context.Context, text string) (string, error) {
	if strings.TrimSpace(text) == "" {
		return "", fmt.Errorf("text cannot be empty")
	}

	// Prepare the request payload - matches Python SDK structure
	request := OpenRouterRequest{
		Model: "deepseek/deepseek-r1-0528:free", // Correct model name from your docs
		Messages: []Message{
			{
				Role:    "user",
				Content: fmt.Sprintf("Переведи следующий текст на русский язык. Сохрани оригинальное форматирование и структуру. Переводи только содержание, не добавляй никаких комментариев или пояснений:\n\n%s", text),
			},
		},
		Stream: false, // We want complete response, not streaming
	}

	// Convert to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers exactly as in Python example
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+t.apiKey)
	// Optional headers for OpenRouter rankings (as shown in Python docs)
	req.Header.Set("HTTP-Referer", "https://github.com/w1zzzle/ai-newsbot")
	req.Header.Set("X-Title", "AI News Bot")

	// Make the request
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var apiResponse OpenRouterResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check for API errors
	if apiResponse.Error != nil {
		return "", fmt.Errorf("API error: %s", apiResponse.Error.Message)
	}

	// Validate response structure
	if len(apiResponse.Choices) == 0 {
		return "", fmt.Errorf("no translation choices returned")
	}

	translatedText := strings.TrimSpace(apiResponse.Choices[0].Message.Content)
	if translatedText == "" {
		return "", fmt.Errorf("empty translation returned")
	}

	return translatedText, nil
}

// TranslateBatch translates multiple texts to Russian
func (t *Translator) TranslateBatch(ctx context.Context, texts []string) ([]string, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("no texts to translate")
	}

	results := make([]string, len(texts))
	for i, text := range texts {
		// Add delay between requests to avoid rate limiting
		if i > 0 {
			select {
			case <-time.After(1 * time.Second): // More conservative delay for free tier
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		translated, err := t.TranslateToRussian(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("failed to translate text %d: %w", i, err)
		}
		results[i] = translated
	}

	return results, nil
}

// IsHealthy checks if the translation service is working
func (t *Translator) IsHealthy(ctx context.Context) error {
	testText := "Hello, world!"
	_, err := t.TranslateToRussian(ctx, testText)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	return nil
}

// TranslateRedditPost is a convenience method for translating Reddit posts
// It combines title and content appropriately for better translation
func (t *Translator) TranslateRedditPost(ctx context.Context, title, content string) (translatedTitle, translatedContent string, err error) {
	// If content is empty, just translate title
	if strings.TrimSpace(content) == "" {
		translatedTitle, err = t.TranslateToRussian(ctx, title)
		return translatedTitle, "", err
	}

	// Combine title and content for better context
	combined := fmt.Sprintf("ЗАГОЛОВОК: %s\n\nСОДЕРЖАНИЕ: %s", title, content)
	
	translated, err := t.TranslateToRussian(ctx, combined)
	if err != nil {
		return "", "", err
	}

	// Try to split back into title and content
	parts := strings.Split(translated, "\n\n")
	if len(parts) >= 2 {
		// Extract title (remove "ЗАГОЛОВОК:" prefix if present)
		translatedTitle = strings.TrimSpace(strings.TrimPrefix(parts[0], "ЗАГОЛОВОК:"))
		translatedTitle = strings.TrimSpace(strings.TrimPrefix(translatedTitle, "Заголовок:"))
		
		// Extract content (remove "СОДЕРЖАНИЕ:" prefix if present)
		translatedContent = strings.TrimSpace(strings.TrimPrefix(strings.Join(parts[1:], "\n\n"), "СОДЕРЖАНИЕ:"))
		translatedContent = strings.TrimSpace(strings.TrimPrefix(translatedContent, "Содержание:"))
	} else {
		// Fallback: treat entire translation as content
		translatedTitle = title // Keep original title
		translatedContent = translated
	}

	return translatedTitle, translatedContent, nil
}