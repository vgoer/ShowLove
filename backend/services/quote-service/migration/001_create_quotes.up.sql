CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE IF NOT EXISTS daily_quotes (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    text_zh        TEXT NOT NULL,
    text_en        TEXT NOT NULL,
    author         VARCHAR(100),
    background_url TEXT DEFAULT '',
    scheduled_date DATE NOT NULL UNIQUE,
    pushed         BOOLEAN DEFAULT FALSE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_quotes_date ON daily_quotes(scheduled_date);
