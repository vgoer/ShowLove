CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE IF NOT EXISTS mood_entries (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL,
    mood_level INTEGER NOT NULL CHECK (mood_level >= 1 AND mood_level <= 10),
    mood_label VARCHAR(50) NOT NULL,
    note       TEXT DEFAULT '',
    created_at DATE NOT NULL DEFAULT CURRENT_DATE
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_moods_user_date ON mood_entries(user_id, created_at);
CREATE INDEX IF NOT EXISTS idx_moods_user ON mood_entries(user_id, created_at DESC);
