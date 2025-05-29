package storage

import (
    "context"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

type Post struct {
    ID            int       `json:"id"`
    RedditID      string    `json:"reddit_id"`
    Title         string    `json:"title"`
    Body          string    `json:"body"`
    MediaURLs     []string  `json:"media_urls"`
    TranslatedBody string   `json:"translated_body"`
    PublishedAt   *time.Time `json:"published_at"`
    CreatedAt     time.Time `json:"created_at"`
}

type Store interface {
    SavePost(ctx context.Context, p Post) error
    IsPostSeen(ctx context.Context, redditID string) (bool, error)
    ListUnpublishedPosts(ctx context.Context) ([]Post, error)
    MarkPublished(ctx context.Context, redditID string) error
    Close() error
}

type PostgresStore struct {
    pool *pgxpool.Pool
}

func NewPostgresStore(dsn string) (*PostgresStore, error) {
    pool, err := pgxpool.New(context.Background(), dsn)
    if err != nil {
        return nil, err
    }

    // Test connection
    if err := pool.Ping(context.Background()); err != nil {
        return nil, err
    }

    return &PostgresStore{pool: pool}, nil
}

func (s *PostgresStore) SavePost(ctx context.Context, p Post) error {
    query := `
        INSERT INTO posts (reddit_id, title, body, media_urls, translated_body)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (reddit_id) DO UPDATE SET
            title = EXCLUDED.title,
            body = EXCLUDED.body,
            media_urls = EXCLUDED.media_urls,
            translated_body = EXCLUDED.translated_body
    `
    
    _, err := s.pool.Exec(ctx, query, p.RedditID, p.Title, p.Body, p.MediaURLs, p.TranslatedBody)
    return err
}

func (s *PostgresStore) IsPostSeen(ctx context.Context, redditID string) (bool, error) {
    var exists bool
    query := `SELECT EXISTS(SELECT 1 FROM posts WHERE reddit_id = $1)`
    
    err := s.pool.QueryRow(ctx, query, redditID).Scan(&exists)
    return exists, err
}

func (s *PostgresStore) ListUnpublishedPosts(ctx context.Context) ([]Post, error) {
    query := `
        SELECT id, reddit_id, title, body, media_urls, translated_body, published_at, created_at
        FROM posts
        WHERE published_at IS NULL AND translated_body IS NOT NULL AND translated_body != ''
        ORDER BY created_at ASC
    `
    
    rows, err := s.pool.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var p Post
        
        err := rows.Scan(&p.ID, &p.RedditID, &p.Title, &p.Body, &p.MediaURLs, &p.TranslatedBody, &p.PublishedAt, &p.CreatedAt)
        if err != nil {
            return nil, err
        }
        
        posts = append(posts, p)
    }

    return posts, rows.Err()
}

func (s *PostgresStore) MarkPublished(ctx context.Context, redditID string) error {
    query := `UPDATE posts SET published_at = NOW() WHERE reddit_id = $1`
    _, err := s.pool.Exec(ctx, query, redditID)
    return err
}

func (s *PostgresStore) Close() error {
    s.pool.Close()
    return nil
}