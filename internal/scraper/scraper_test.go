package scraper

import (
    "context"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestRedditScraper_FetchPosts(t *testing.T) {
    // Create a mock HTTP server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        mockHTML := `
        <html>
            <body>
                <div data-testid="post-container" data-post-id="test123">
                    <h3 data-testid="post-title">Test AI News Title</h3>
                    <div data-testid="post-content">
                        <p>This is a test post about artificial intelligence.</p>
                    </div>
                    <img src="https://example.com/image.jpg" alt="test image">
                </div>
                <div data-testid="post-container" data-post-id="test456">
                    <h3 data-testid="post-title">Another AI Post</h3>
                    <div data-testid="post-content">
                        <p>Another interesting AI development.</p>
                    </div>
                </div>
            </body>
        </html>`
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(mockHTML))
    }))
    defer server.Close()

    scraper := New([]string{server.URL}, 100)
    ctx := context.Background()

    posts, err := scraper.FetchPosts(ctx)
    require.NoError(t, err)
    assert.Len(t, posts, 2)

    // Verify first post
    assert.Equal(t, "test123", posts[0].RedditID)
    assert.Equal(t, "Test AI News Title", posts[0].Title)
    assert.Equal(t, "This is a test post about artificial intelligence.", posts[0].Body)
    assert.Contains(t, posts[0].MediaURLs, "https://example.com/image.jpg")

    // Verify second post
    assert.Equal(t, "test456", posts[1].RedditID)
    assert.Equal(t, "Another AI Post", posts[1].Title)
    assert.Equal(t, "Another interesting AI development.", posts[1].Body)
}

func TestRedditScraper_FetchPosts_InvalidURL(t *testing.T) {
    scraper := New([]string{"invalid-url"}, 100)
    ctx := context.Background()

    posts, err := scraper.FetchPosts(ctx)
    assert.Error(t, err)
    assert.Nil(t, posts)
}

func TestRedditScraper_FetchPosts_HTTPError(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusInternalServerError)
    }))
    defer server.Close()

    scraper := New([]string{server.URL}, 100)
    ctx := context.Background()

    posts, err := scraper.FetchPosts(ctx)
    assert.Error(t, err)
    assert.Nil(t, posts)
}