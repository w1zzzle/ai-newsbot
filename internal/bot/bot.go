package bot

import (
    "context"
    "fmt"
    "strings"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/w1zzzle/ai-newsbot/internal/storage"
)

type Bot interface {
    SendPost(ctx context.Context, post storage.Post) error
}

type TelegramBot struct {
    api    *tgbotapi.BotAPI
    chatID int64
}

func New(token string, chatID int64) (*TelegramBot, error) {
    bot, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        return nil, fmt.Errorf("failed to create telegram bot: %w", err)
    }

    return &TelegramBot{
        api:    bot,
        chatID: chatID,
    }, nil
}

func (b *TelegramBot) SendPost(ctx context.Context, post storage.Post) error {
    // Prepare message text
    messageText := b.formatMessage(post)

    // Send media if available
    if len(post.MediaURLs) > 0 {
        return b.sendMediaWithCaption(ctx, post.MediaURLs[0], messageText)
    }

    // Send text message
    return b.sendTextMessage(ctx, messageText)
}

func (b *TelegramBot) formatMessage(post storage.Post) string {
    var message strings.Builder
    
    if post.Title != "" {
        message.WriteString("📰 *")
        message.WriteString(b.escapeMarkdown(post.Title))
        message.WriteString("*\n\n")
    }

    if post.TranslatedBody != "" {
        message.WriteString(b.escapeMarkdown(post.TranslatedBody))
    }

    return message.String()
}

func (b *TelegramBot) sendTextMessage(ctx context.Context, text string) error {
    msg := tgbotapi.NewMessage(b.chatID, text)
    msg.ParseMode = tgbotapi.ModeMarkdown

    _, err := b.api.Send(msg)
    return err
}

func (b *TelegramBot) sendMediaWithCaption(ctx context.Context, mediaURL, caption string) error {
    // Determine media type based on URL
    if b.isImageURL(mediaURL) {
        return b.sendPhoto(ctx, mediaURL, caption)
    } else if b.isVideoURL(mediaURL) {
        return b.sendVideo(ctx, mediaURL, caption)
    } else if b.isGifURL(mediaURL) {
        return b.sendAnimation(ctx, mediaURL, caption)
    }

    // Fallback to text message if media type is unknown
    return b.sendTextMessage(ctx, caption)
}

func (b *TelegramBot) sendPhoto(ctx context.Context, photoURL, caption string) error {
    msg := tgbotapi.NewPhoto(b.chatID, tgbotapi.FileURL(photoURL))
    msg.Caption = caption
    msg.ParseMode = tgbotapi.ModeMarkdown

    _, err := b.api.Send(msg)
    return err
}

func (b *TelegramBot) sendVideo(ctx context.Context, videoURL, caption string) error {
    msg := tgbotapi.NewVideo(b.chatID, tgbotapi.FileURL(videoURL))
    msg.Caption = caption
    msg.ParseMode = tgbotapi.ModeMarkdown

    _, err := b.api.Send(msg)
    return err
}

func (b *TelegramBot) sendAnimation(ctx context.Context, animationURL, caption string) error {
    msg := tgbotapi.NewAnimation(b.chatID, tgbotapi.FileURL(animationURL))
    msg.Caption = caption
    msg.ParseMode = tgbotapi.ModeMarkdown

    _, err := b.api.Send(msg)
    return err
}

func (b *TelegramBot) isImageURL(url string) bool {
    return strings.HasSuffix(strings.ToLower(url), ".jpg") ||
           strings.HasSuffix(strings.ToLower(url), ".jpeg") ||
           strings.HasSuffix(strings.ToLower(url), ".png") ||
           strings.HasSuffix(strings.ToLower(url), ".webp")
}

func (b *TelegramBot) isVideoURL(url string) bool {
    return strings.HasSuffix(strings.ToLower(url), ".mp4") ||
           strings.HasSuffix(strings.ToLower(url), ".mov") ||
           strings.HasSuffix(strings.ToLower(url), ".avi")
}

func (b *TelegramBot) isGifURL(url string) bool {
    return strings.HasSuffix(strings.ToLower(url), ".gif")
}

func (b *TelegramBot) escapeMarkdown(text string) string {
    replacer := strings.NewReplacer(
        "*", "\\*",
        "_", "\\_",
        "`", "\\`",
        "[", "\\[",
        "]", "\\]",
        "(", "\\(",
        ")", "\\)",
    )
    return replacer.Replace(text)
}