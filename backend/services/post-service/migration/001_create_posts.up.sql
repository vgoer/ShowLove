-- 001_create_posts.up.sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS posts (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    author_id          UUID NOT NULL,
    author_nickname    VARCHAR(100) NOT NULL,
    author_avatar      TEXT DEFAULT '',
    content            TEXT NOT NULL,
    mood_tag           VARCHAR(50) NOT NULL,
    images             JSONB DEFAULT '[]',
    voice_url          TEXT DEFAULT '',
    sticker_hug        INTEGER DEFAULT 0,
    sticker_cheer      INTEGER DEFAULT 0,
    sticker_understand INTEGER DEFAULT 0,
    comment_count      INTEGER DEFAULT 0,
    has_ai_reply       BOOLEAN DEFAULT FALSE,
    is_reported        BOOLEAN DEFAULT FALSE,
    is_hidden          BOOLEAN DEFAULT FALSE,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_posts_author ON posts(author_id);
CREATE INDEX IF NOT EXISTS idx_posts_created ON posts(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_posts_mood_tag ON posts(mood_tag);
CREATE INDEX IF NOT EXISTS idx_posts_comment_count ON posts(comment_count DESC);
