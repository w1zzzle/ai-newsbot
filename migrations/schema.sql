-- AI NewsBot Database Schema

CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    reddit_id TEXT UNIQUE NOT NULL,
    title TEXT NOT NULL,
    body TEXT,
    media_urls TEXT[],
    translated_body TEXT,
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_posts_reddit_id ON posts(reddit_id);
CREATE INDEX IF NOT EXISTS idx_posts_published_at ON posts(published_at);
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at);