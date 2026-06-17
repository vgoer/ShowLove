#!/bin/bash
# =============================================================================
# 初始化 PostgreSQL 多数据库
# 在 postgres 容器首次启动时自动执行
# =============================================================================

set -e

# 使用 psql 创建多个数据库
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE DATABASE users_db;
    CREATE DATABASE posts_db;
    CREATE DATABASE comments_db;
    CREATE DATABASE moods_db;
    CREATE DATABASE quotes_db;
    CREATE DATABASE notifications_db;
EOSQL

echo "✅ 全部数据库初始化完成: users_db, posts_db, comments_db, moods_db, quotes_db, notifications_db"
