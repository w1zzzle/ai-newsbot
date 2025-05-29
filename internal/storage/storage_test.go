package storage

import (
    "context"
    //"database/sql"
    "testing"
    "time"

    _ "github.com/lib/pq"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// Note: This test requires a test database. In a real project, you might use:
// - Docker test containers
// - SQLite for testing
// - In-memory database

func TestPostgresStore_SavePost(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping database integration test")
    }

    // This would require a test database connection
    // For demonstration purposes, we'll create a mock test
    t.Skip("Integration test - requires test database setup")

    // Example of how the test would look:
    /*
    store := setupTestStore(t)
    defer store.Close()

    post := Post{
        RedditID:       "test123",
        Title:          "Test Post",
        Body:           "This is a test post",
        MediaURLs:      []string{"https://example.com/image.jpg"},
        TranslatedBody: "Это тестовый пост",
    }

    err := store.SavePost(context.Background(), post)
    require.NoError(t, err)

    // Verify the post was saved
    seen, err := store.IsPostSeen(context.Background(), "test123")
    require.NoError(t, err)
    assert.True(t, seen)
    */
}

func TestPostgresStore_IsPostSeen(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping database integration test")
    }

    t.Skip("Integration test - requires test database setup")
}

func TestPostgresStore_ListUnpublishedPosts(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping database integration test")
    }

    t.Skip("Integration test - requires test database setup")
}

func TestPostgresStore_MarkPublished(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping database integration test")
    }

    t.Skip("Integration test - requires test database setup")
}

// MockStore for unit testing other components
type MockStore struct {
    posts     map[string]Post
    published map[string]bool
}

func NewMockStore() *MockStore {
    return &MockStore{
        posts:     make(map[string]Post),
        published: make(map[string]bool),
    }
}

func (m *MockStore) SavePost(ctx context.Context, p Post) error {
    m.posts[p.RedditID] = p
    return nil
}

func (m *MockStore) IsPostSeen(ctx context.Context, redditID string) (bool, error) {
    _, exists := m.posts[redditID]
    return exists, nil
}

func (m *MockStore) ListUnpublishedPosts(ctx context.Context) ([]Post, error) {
    var unpublished []Post
    for _, post := range m.posts {
        if !m.published[post.RedditID] && post.TranslatedBody != "" {
            unpublished = append(unpublished, post)
        }
    }
    return unpublished, nil
}

func (m *MockStore) MarkPublished(ctx context.Context, redditID string) error {
    m.published[redditID] = true
    return nil
}

func (m *MockStore) Close() error {
    return nil
}

func TestMockStore(t *testing.T) {
    store := NewMockStore()
    ctx := context.Background()

    post := Post{
        RedditID:       "test123",
        Title:          "Test Post",
        Body:           "This is a test post",
        TranslatedBody: "Это тестовый пост",
        CreatedAt:      time.Now(),
    }

    // Test SavePost
    err := store.SavePost(ctx, post)
    require.NoError(t, err)

    // Test IsPostSeen
    seen, err := store.IsPostSeen(ctx, "test123")
    require.NoError(t, err)
    assert.True(t, seen)

    // Test ListUnpublishedPosts
    unpublished, err := store.ListUnpublishedPosts(ctx)
    require.NoError(t, err)
    assert.Len(t, unpublished, 1)
    assert.Equal(t, "test123", unpublished[0].RedditID)

    // Test MarkPublished
    err = store.MarkPublished(ctx, "test123")
    require.NoError(t, err)

    // Verify post is no longer unpublished
    unpublished, err = store.ListUnpublishedPosts(ctx)
    require.NoError(t, err)
    assert.Len(t, unpublished, 0)
}