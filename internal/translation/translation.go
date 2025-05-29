package translation

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type Translator interface {
    TranslateToRussian(ctx context.Context, text string) (string, error)
}

type OpenRouterTranslator struct {
    apiKey string
    client *http.Client
}

type OpenRouterRequest struct {
    Model    string    `json:"model"`
    Messages []Message `json:"messages"`
}

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type OpenRouterResponse struct {
    Choices []struct {
        Message struct {
            Content string `json:"content"`
        } `json:"message"`
    } `json:"choices"`
    Error *struct {
        Message string `json:"message"`
    } `json:"error,omitempty"`
}

func New(apiKey string) *OpenRouterTranslator {
    return &OpenRouterTranslator{
        apiKey: apiKey,
        client: &http.Client{
            Timeout: 60 * time.Second,
        },
    }
}

func (t *OpenRouterTranslator) TranslateToRussian(ctx context.Context, text string) (string, error) {
    if text == "" {
        return "", nil
    }

    prompt := fmt.Sprintf("Переведи следующий текст на русский язык, сохраняя тон и формат:\n«%s»", text)

    reqBody := OpenRouterRequest{
        Model: "deepseek/deepseek-r1-distill-llama-70b",
        Messages: []Message{
            {
                Role:    "user",
                Content: prompt,
            },
        },
    }

    jsonData, err := json.Marshal(reqBody)
    if err != nil {
        return "", fmt.Errorf("failed to marshal request: %w", err)
    }

    req, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
    if err != nil {
        return "", fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+t.apiKey)
    req.Header.Set("HTTP-Referer", "https://github.com/youruser/ai-newsbot")

    resp, err := t.client.Do(req)
    if err != nil {
        return "", fmt.Errorf("failed to make request: %w", err)
    }
    defer resp.Body.Close()

    var apiResp OpenRouterResponse
    if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
        return "", fmt.Errorf("failed to decode response: %w", err)
    }

    if apiResp.Error != nil {
        return "", fmt.Errorf("API error: %s", apiResp.Error.Message)
    }

    if len(apiResp.Choices) == 0 {
        return "", fmt.Errorf("no translation choices returned")
    }

    translation := apiResp.Choices[0].Message.Content
    
    // Clean up the translation (remove quotes if they wrap the entire response)
    if len(translation) >= 2 && translation[0] == '«' && translation[len(translation)-1] == '»' {
        translation = translation[1 : len(translation)-1]
    }

    return translation, nil
}