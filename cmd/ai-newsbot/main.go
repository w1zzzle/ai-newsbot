package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/joho/godotenv"
    "github.com/robfig/cron/v3"
    "github.com/w1zzzle/ai-newsbot/internal/app"
    "github.com/w1zzzle/ai-newsbot/internal/bot"
    "github.com/w1zzzle/ai-newsbot/internal/config"
    "github.com/w1zzzle/ai-newsbot/internal/scraper"
    "github.com/w1zzzle/ai-newsbot/internal/storage"
    "github.com/w1zzzle/ai-newsbot/internal/translation"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }

    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Initialize storage
    store, err := storage.NewPostgresStore(cfg.PostgresDSN)
    if err != nil {
        log.Fatalf("Failed to initialize storage: %v", err)
    }
    defer store.Close()

    // Initialize components
    scraperSvc := scraper.New(cfg.RedditURLs, cfg.UpvoteThreshold)
    translatorSvc := translation.New(cfg.OpenRouterAPIKey)
    
    // Fix: bot.New now returns (*TelegramBot, error)
    botSvc, err := bot.New(cfg.TelegramBotToken, cfg.TelegramChatID)
    if err != nil {
        log.Fatalf("Failed to initialize bot: %v", err)
    }

    // Initialize application
    application := app.New(store, scraperSvc, translatorSvc, botSvc)

    // Setup cron scheduler
    c := cron.New()
    _, err = c.AddFunc("@hourly", func() {
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
        defer cancel()
        
        if err := application.RunPipeline(ctx); err != nil {
            log.Printf("Pipeline execution failed: %v", err)
        }
    })
    if err != nil {
        log.Fatalf("Failed to setup cron job: %v", err)
    }

    c.Start()
    log.Println("AI NewsBot started successfully. Running hourly...")

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down AI NewsBot...")
    ctx := c.Stop()
    <-ctx.Done()
    log.Println("AI NewsBot stopped")
}