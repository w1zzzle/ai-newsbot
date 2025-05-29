package scraper

import (
    "context"
    "fmt"
    "net/http"
    "regexp"
    "strconv"
    "strings"
    "time"

    "github.com/PuerkitoBio/goquery"
    "github.com/youruser/ai-newsbot/internal/storage"
)

type Scraper interface {
    FetchPosts(ctx context.Context) ([]storage.Post, error)
}

type RedditScraper struct {
    urls            []string
    upvoteThreshold int
    client          *http.Client
}

func New(urls []string, upvoteThreshold int) *RedditScraper {
    return &RedditScraper{
        urls:            urls,
        upvoteThreshold: upvoteThreshold,
        client: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (s *RedditScraper) FetchPosts(ctx context.Context) ([]storage.Post, error) {
    var allPosts []storage.Post

    for _, url := range s.urls {
        posts, err := s.fetchFromURL(ctx, url)
        if err != nil {
            return nil, fmt.Errorf("failed to fetch from %s: %w", url, err)
        }
        allPosts = append(allPosts, posts...)
    }

    return allPosts, nil
}

func (s *RedditScraper) fetchFromURL(ctx context.Context, url string) ([]storage.Post, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }

    // Set User-Agent to avoid being blocked
    req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; AI-NewsBot/1.0)")

    resp, err := s.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
    }

    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return nil, err
    }

    var posts []storage.Post

    // Find post elements (this selector might need adjustment based on Reddit's current HTML structure)
    doc.Find("div[data-testid='post-container']").Each(func(i int, postEl *goquery.Selection) {
        post := s.extractPost(postEl)
        if post != nil && s.meetsThreshold(*post) {
            posts = append(posts, *post)
        }
    })

    return posts, nil
}

func (s *RedditScraper) extractPost(postEl *goquery.Selection) *storage.Post {
    // Extract Reddit ID from data attributes or URL
    redditID := s.extractRedditID(postEl)
    if redditID == "" {
        return nil
    }

    // Extract title
    titleEl := postEl.Find("h3[data-testid='post-title']")
    title := strings.TrimSpace(titleEl.Text())
    if title == "" {
        return nil
    }

    // Extract body text
    bodyEl := postEl.Find("div[data-testid='post-content'] p")
    body := strings.TrimSpace(bodyEl.Text())

    // Extract media URLs
    mediaURLs := s.extractMediaURLs(postEl)

    return &storage.Post{
        RedditID:  redditID,
        Title:     title,
        Body:      body,
        MediaURLs: mediaURLs,
        CreatedAt: time.Now(),
    }
}

func (s *RedditScraper) extractRedditID(postEl *goquery.Selection) string {
    // Try to get from data-post-id attribute
    if id, exists := postEl.Attr("data-post-id"); exists {
        return id
    }

    // Try to extract from permalink URL
    linkEl := postEl.Find("a[data-testid='post-title']")
    if href, exists := linkEl.Attr("href"); exists {
        re := regexp.MustCompile(`/comments/([a-zA-Z0-9]+)/`)
        matches := re.FindStringSubmatch(href)
        if len(matches) > 1 {
            return matches[1]
        }
    }

    return ""
}

func (s *RedditScraper) extractMediaURLs(postEl *goquery.Selection) []string {
    var urls []string

    // Extract image URLs
    postEl.Find("img").Each(func(i int, imgEl *goquery.Selection) {
        if src, exists := imgEl.Attr("src"); exists && strings.HasPrefix(src, "http") {
            urls = append(urls, src)
        }
    })

    // Extract video URLs
    postEl.Find("video").Each(func(i int, videoEl *goquery.Selection) {
        if src, exists := videoEl.Attr("src"); exists && strings.HasPrefix(src, "http") {
            urls = append(urls, src)
        }
    })

    return urls
}

func (s *RedditScraper) meetsThreshold(post storage.Post) bool {
    // For now, we'll assume all posts meet the threshold
    // In a real implementation, you'd extract upvote count and compare
    return true
}

func (s *RedditScraper) extractUpvotes(postEl *goquery.Selection) int {
    upvoteEl := postEl.Find("span[data-testid='upvote-count']")
    upvoteText := strings.TrimSpace(upvoteEl.Text())
    
    // Clean up the text (remove 'k', 'M' suffixes, etc.)
    upvoteText = strings.ReplaceAll(upvoteText, "k", "000")
    upvoteText = strings.ReplaceAll(upvoteText, "M", "000000")
    upvoteText = regexp.MustCompile(`[^\d]`).ReplaceAllString(upvoteText, "")
    
    upvotes, _ := strconv.Atoi(upvoteText)
    return upvotes
}