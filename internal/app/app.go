package app

import (
    "context"
    "fmt"
    "log"

    "github.com/w1zzzle/ai-newsbot/internal/bot"
    "github.com/w1zzzle/ai-newsbot/internal/scraper"
    "github.com/w1zzzle/ai-newsbot/internal/storage"
    "github.com/w1zzzle/ai-newsbot/internal/translation"
)

type App struct {
    store      storage.Store
    scraper    scraper.Scraper
    translator translation.Translator
    bot        bot.Bot
}

func New(store storage.Store, scraper scraper.Scraper, translator translation.Translator, bot bot.Bot) *App {
    return &App{
        store:      store,
        scraper:    scraper,
        translator: translator,
        bot:        bot,
    }
}

func (a *App) RunPipeline(ctx context.Context) error {
    log.Println("Starting AI NewsBot pipeline...")

    // Step 1: Fetch posts from Reddit
    log.Println("Fetching posts from Reddit...")
    posts, err := a.scraper.FetchPosts(ctx)
    if err != nil {
        return fmt.Errorf("failed to fetch posts: %w", err)
    }
    log.Printf("Fetched %d posts", len(posts))

    // Step 2: Filter new posts and translate them
    newPosts := 0
    for _, post := range posts {
        // Check if we've seen this post before
        seen, err := a.store.IsPostSeen(ctx, post.RedditID)
        if err != nil {
            log.Printf("Error checking if post %s is seen: %v", post.RedditID, err)
            continue
        }

        if seen {
            continue
        }

        // Translate the post
        log.Printf("Translating post: %s", post.Title)
        translatedBody, err := a.translator.TranslateToRussian(ctx, post.Body)
        if err != nil {
            log.Printf("Failed to translate post %s: %v", post.RedditID, err)
            continue
        }

        post.TranslatedBody = translatedBody

        // Save the post
        if err := a.store.SavePost(ctx, post); err != nil {
            log.Printf("Failed to save post %s: %v", post.RedditID, err)
            continue
        }

        newPosts++
    }

    log.Printf("Processed %d new posts", newPosts)

    // Step 3: Publish unpublished posts
    log.Println("Publishing unpublished posts...")
    unpublishedPosts, err := a.store.ListUnpublishedPosts(ctx)
    if err != nil {
        return fmt.Errorf("failed to list unpublished posts: %w", err)
    }

    published := 0
    for _, post := range unpublishedPosts {
        if err := a.bot.SendPost(ctx, post); err != nil {
            log.Printf("Failed to send post %s: %v", post.RedditID, err)
            continue
        }

        if err := a.store.MarkPublished(ctx, post.RedditID); err != nil {
            log.Printf("Failed to mark post %s as published: %v", post.RedditID, err)
            continue
        }

        published++
        log.Printf("Published post: %s", post.Title)
    }

    log.Printf("Pipeline completed. Published %d posts", published)
    return nil
}