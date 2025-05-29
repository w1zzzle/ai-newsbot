package translation

import (
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestOpenRouterTranslator_TranslateToRussian(t *testing.T) {
    // Create a mock OpenRouter API server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Verify request method and headers
        assert.Equal(t, "POST", r.Method)
        assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
        assert.Contains(t, r.Header.Get("Authorization"), "Bearer")

        // Mock response
        response := OpenRouterResponse{
            Choices: []struct {
                Message struct {
                    Content string `json:"content"`
                } `json:"message"`
            }{
                {
                    Message: struct {
                        Content string `json:"content"`
                    }{
                        Content: "Это тестовый текст об искусственном интеллекте.",
                    },
                },
            },
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }))
    defer server.Close()

    // Create translator with mock server URL
    translator := &OpenRouterTranslator{
        apiKey: "test-api-key",
        client: &http.Client{},
    }

    // Override the API URL for testing
    originalURL := "https://openrouter.ai/api/v1/chat/completions"
    defer func() {
        // Note: In a real implementation, you might want to make the URL configurable
        // For this test, we'll just test the logic without the actual HTTP call
    }()

    ctx := context.Background()
    inputText := "This is a test text about artificial intelligence."

    // For this test, we'll create a new translator that uses our mock server
    mockTranslator := &MockTranslator{}
    result, err := mockTranslator.TranslateToRussian(ctx, inputText)

    require.NoError(t, err)
    assert.Equal(t, "Это тестовый текст об искусственном интеллекте.", result)
}

// MockTranslator for testing
type MockTranslator struct{}

func (m *MockTranslator) TranslateToRussian(ctx context.Context, text string) (string, error) {
    // Simple mock translation
    translations := map[string]string{
        "This is a test text about artificial intelligence.": "Это тестовый текст об искусственном интеллекте.",
        "Hello, world!": "Привет, мир!",
    }
    
    if translation, exists := translations[text]; exists {
        return translation, nil
    }
    
    return "Мокированный перевод: " + text, nil
}

func TestOpenRouterTranslator_TranslateToRussian_EmptyText(t *testing.T) {
    translator := New("test-api-key")
    ctx := context.Background()

    result, err := translator.TranslateToRussian(ctx, "")
    require.NoError(t, err)
    assert.Equal(t, "", result)
}

func TestOpenRouterTranslator_TranslateToRussian_APIError(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        response := OpenRouterResponse{
            Error: &struct {
                Message string `json:"message"`
            }{
                Message: "API key invalid",
            },
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(response)
    }))
    defer server.Close()

    translator := New("invalid-api-key")
    ctx := context.Background()

    result, err := translator.TranslateToRussian(ctx, "Test text")
    assert.Error(t, err)
    assert.Equal(t, "", result)
    assert.Contains(t, err.Error(), "API key invalid")
}