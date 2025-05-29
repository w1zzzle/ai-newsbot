package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Config struct {
    RedditURLs        []string
    UpvoteThreshold   int
    PostgresDSN       string
    OpenRouterAPIKey  string
    TelegramBotToken  string
    TelegramChatID    int64
}

func Load() (*Config, error) {
    cfg := &Config{}

    // Reddit URLs
    redditURLsStr := os.Getenv("REDDIT_URLS")
    if redditURLsStr == "" {
        redditURLsStr = "https://www.reddit.com/r/ArtificialIntelligence/top/"
    }
    cfg.RedditURLs = strings.Split(redditURLsStr, ",")

    // Upvote threshold
    thresholdStr := os.Getenv("UPVOTE_THRESHOLD")
    if thresholdStr == "" {
        cfg.UpvoteThreshold = 100
    } else {
        threshold, err := strconv.Atoi(thresholdStr)
        if err != nil {
            return nil, fmt.Errorf("invalid upvote threshold: %w", err)
        }
        cfg.UpvoteThreshold = threshold
    }

    // Database
    cfg.PostgresDSN = os.Getenv("POSTGRES_DSN")
    if cfg.PostgresDSN == "" {
        return nil, fmt.Errorf("POSTGRES_DSN is required")
    }

    // OpenRouter API Key
    cfg.OpenRouterAPIKey = os.Getenv("OPENROUTER_API_KEY")
    if cfg.OpenRouterAPIKey == "" {
        return nil, fmt.Errorf("OPENROUTER_API_KEY is required")
    }

    // Telegram Bot Token
    cfg.TelegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
    if cfg.TelegramBotToken == "" {
        return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
    }

    // Telegram Chat ID
    chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")
    if chatIDStr == "" {
        return nil, fmt.Errorf("TELEGRAM_CHAT_ID is required")
    }
    chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
    if err != nil {
        return nil, fmt.Errorf("invalid telegram chat ID: %w", err)
    }
    cfg.TelegramChatID = chatID

    return cfg, nil
}