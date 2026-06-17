-- 001_create_users.down.sql
-- 回滚 users 和 refresh_tokens 表

DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;
