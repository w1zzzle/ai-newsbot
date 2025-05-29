package bot

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/youruser/ai-newsbot/internal/storage"
)

func TestTelegramBot_FormatMessage(t *testing.T) {
    bot := &TelegramBot{chatID: 123}

    post := storage.Post{
        Title:          "AI News: Machine Learning Breakthrough",
        TranslatedBody: "Исследователи достигли нового прорыва в машинном обучении.",
    }

    message := bot.formatMessage(post)
    
    assert.Contains(t, message, "📰 *AI News: Machine Learning Breakthrough*")
    assert.Contains(t, message, "Исследователи достигли нового прорыва в машинном обучении.")
}

func TestTelegramBot_EscapeMarkdown(t *testing.T) {
    bot := &TelegramBot{}

    testCases := []struct {
        input    string
        expected string
    }{
        {"Hello *world*", "Hello \\*world\\*"},
        {"Test_underscore", "Test\\_underscore"},
        {"Code `block`", "Code \\`block\\`"},
        {"Link [text](url)", "Link \\[text\\]\\(url\\)"},
        {"Normal text", "Normal text"},
    }

    for _, tc := range testCases {
        result := bot.escapeMarkdown(tc.input)
        assert.Equal(t, tc.expected, result)
    }
}

func TestTelegramBot_IsImageURL(t *testing.T) {
    bot := &TelegramBot{}

    testCases := []struct {
        url      string
        expected bool
    }{
        {"https://example.com/image.jpg", true},
        {"https://example.com/image.jpeg", true},
        {"https://example.com/image.png", true},
        {"https://example.com/image.webp", true},
        {"https://example.com/video.mp4", false},
        {"https://example.com/animation.gif", false},
        {"https://example.com/document.pdf", false},
    }

    for _, tc := range testCases {
        result := bot.isImageURL(tc.url)
        assert.Equal(t, tc.expected, result, "URL: %s", tc.url)
    }
}

func TestTelegramBot_IsVideoURL(t *testing.T) {
    bot := &TelegramBot{}

    testCases := []struct {
        url      string
        expected bool
    }{
        {"https://example.com/video.mp4", true},
        {"https://example.com/video.mov", true},
        {"https://example.com/video.avi", true},
        {"https://example.com/image.jpg", false},
        {"https://example.com/animation.gif", false},
    }

    for _, tc := range testCases {
        result := bot.isVideoURL(tc.url)
        assert.Equal(t, tc.expected, result, "URL: %s", tc.url)
    }
}

func TestTelegramBot_IsGifURL(t *testing.T) {
    bot := &TelegramBot{}

    assert.True(t, bot.isGifURL("https://example.com/animation.gif"))
    assert.False(t, bot.isGifURL("https://example.com/image.jpg"))
    assert.False(t, bot.isGifURL("https://example.com/video.mp4"))
}

// MockBot for testing other components
type MockBot struct {
    SentPosts []storage.Post
}

func (m *MockBot) SendPost(ctx context.Context, post storage.Post) error {
    m.SentPosts = append(m.SentPosts, post)
    return nil
}

func TestMockBot(t *testing.T) {
    bot := &MockBot{}
    ctx := context.Background()

    post := storage.Post{
        RedditID:       "test123",
        Title:          "Test Post",
        TranslatedBody: "Тестовый пост",
        CreatedAt:      time.Now(),
    }

    err := bot.SendPost(ctx, post)
    assert.NoError(t, err)
    assert.Len(t, bot.SentPosts, 1)
    assert.Equal(t, "test123", bot.SentPosts[0].RedditID)
}