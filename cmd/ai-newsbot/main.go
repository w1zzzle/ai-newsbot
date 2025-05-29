package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/w1zzzle/ai-newsbot/internal/translation"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENROUTER_API_KEY environment variable is required")
	}

	// Create translator
	translator := translation.New(apiKey)

	// Test translation service health
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := translator.IsHealthy(ctx); err != nil {
		log.Printf("Translation service health check failed: %v", err)
	} else {
		log.Println("Translation service is healthy")
	}

	// Example usage in your scraping workflow
	redditPost := `
Title: New AI Model Breaks Performance Records

Content: Researchers at OpenAI have announced a breakthrough in large language model performance. The new model, called GPT-5, shows significant improvements in reasoning, coding, and multilingual capabilities.

Key features:
- 50% better performance on coding benchmarks
- Improved reasoning abilities
- Better multilingual support
- More efficient training process

The model is expected to be released later this year after safety testing.
	`

	// Translate the post
	translatedPost, err := translator.TranslateToRussian(ctx, redditPost)
	if err != nil {
		log.Fatalf("Translation failed: %v", err)
	}

	fmt.Println("Original post:")
	fmt.Println(redditPost)
	fmt.Println("\nTranslated post:")
	fmt.Println(translatedPost)

	// Example batch translation
	posts := []string{
		"Breaking: New programming language announced",
		"Tech stocks surge after AI breakthrough",
		"Open source project gains massive community support",
	}

	translatedPosts, err := translator.TranslateBatch(ctx, posts)
	if err != nil {
		log.Fatalf("Batch translation failed: %v", err)
	}

	fmt.Println("\nBatch translation results:")
	for i, original := range posts {
		fmt.Printf("Original: %s\n", original)
		fmt.Printf("Russian: %s\n\n", translatedPosts[i])
	}
}

// Example integration with your scraping workflow
func processRedditPost(translator translation.Service, post RedditPost) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Combine title and content for translation
	fullText := fmt.Sprintf("Заголовок: %s\n\nСодержание: %s", post.Title, post.Content)

	// Translate to Russian
	translatedText, err := translator.TranslateToRussian(ctx, fullText)
	if err != nil {
		return fmt.Errorf("failed to translate post: %w", err)
	}

	// Parse translated text back to title and content if needed
	// This is a simple example - you might want more sophisticated parsing
	post.TranslatedTitle = extractTitleFromTranslation(translatedText)
	post.TranslatedContent = extractContentFromTranslation(translatedText)

	return nil
}

// Mock structures for example
type RedditPost struct {
	Title             string
	Content           string
	URL               string
	Upvotes          int
	TranslatedTitle   string
	TranslatedContent string
}

func extractTitleFromTranslation(translatedText string) string {
	// Implementation would parse "Заголовок: ..." from translated text
	// This is a placeholder
	return "Translated Title"
}

func extractContentFromTranslation(translatedText string) string {
	// Implementation would parse "Содержание: ..." from translated text
	// This is a placeholder
	return "Translated Content"
}