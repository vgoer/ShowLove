# Show Love 部署指南

## 架构概览

```
Flutter App → API Gateway (Gin, :8080) → gRPC → 7微服务
                                              ↓
                        PostgreSQL · Redis · MinIO · NATS
```

## 前置条件

- [Docker](https://docs.docker.com/get-docker/) ≥ 24
- [Docker Compose](https://docs.docker.com/compose/) ≥ 2
- [Go](https://go.dev/dl/) ≥ 1.22 (本地开发)
- [protoc](https://grpc.io/docs/protoc-installation/) (修改 proto 时需要)

## 快速启动 (Docker Compose)

### 1. 克隆项目

```bash
git clone https://github.com/your-org/show-love.git
cd show-love/backend
```

### 2. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env，至少设置:
#   JWT_SECRET=随机生成的安全密钥
#   AI_PROVIDER=openai (或 deepseek)
#   OPENAI_API_KEY=sk-xxx (使用AI功能时需要)
```

### 3. 启动基础设施

```bash
docker compose up -d postgres redis minio minio-init nats
```

等待所有服务 healthy:
```bash
docker compose ps
# 确认 postgres, redis, minio, nats 都是 (healthy)
```

### 4. 启动全部微服务

```bash
docker compose up -d --build
```

### 5. 验证

```bash
# 健康检查
curl http://localhost:8080/api/v1/health

# 注册用户
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test1234","nickname":"测试"}'
```

### 6. 停止

```bash
docker compose down        # 保留数据卷
docker compose down -v     # 清理全部数据
```

## 本地开发

### 启动依赖服务

```bash
docker compose up -d postgres redis minio minio-init nats
```

### 创建数据库

```bash
# 手动创建各服务的数据库 (如果 init-multi-db.sh 未自动执行)
docker exec -it showlove-postgres psql -U showlove -c "CREATE DATABASE users_db"
docker exec -it showlove-postgres psql -U showlove -c "CREATE DATABASE posts_db"
docker exec -it showlove-postgres psql -U showlove -c "CREATE DATABASE comments_db"
docker exec -it showlove-postgres psql -U showlove -c "CREATE DATABASE moods_db"
docker exec -it showlove-postgres psql -U showlove -c "CREATE DATABASE quotes_db"
docker exec -it showlove-postgres psql -U showlove -c "CREATE DATABASE notifications_db"
```

### 启动各服务 (每个终端窗口一个)

```bash
# 终端 1: 用户服务
cd services/user-service && go run ./cmd/

# 终端 2: 帖子服务
cd services/post-service && go run ./cmd/

# 终端 3: 评论服务
cd services/comment-service && go run ./cmd/

# 终端 4: API 网关
cd gateway && go run ./cmd/

# ... 其他服务按需启动
```

### 运行测试

```bash
# 全部测试
go test ./...

# 单个服务
go test ./services/user-service/... -v

# 带覆盖率
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 生产部署

### 使用 docker-compose.prod.yml

```bash
# 设置生产环境变量
export POSTGRES_USER=showlove
export POSTGRES_PASSWORD=<strong-password>
export JWT_SECRET=<random-64-char-string>
export AI_PROVIDER=openai
export OPENAI_API_KEY=sk-xxx
export STORAGE_PROVIDER=minio
export MINIO_ACCESS_KEY=<access-key>
export MINIO_SECRET_KEY=<secret-key>
export GATEWAY_PORT=8080

# 启动
docker compose -f docker-compose.prod.yml up -d --build
```

### 安全注意事项

1. **JWT_SECRET**: 使用 `openssl rand -base64 64` 生成
2. **数据库密码**: 禁止使用默认密码
3. **API Keys**: 通过 Docker secrets 或 Vault 管理
4. **HTTPS**: 生产环境请在前面加 Nginx/Caddy 反向代理，启用 TLS

### 监控与日志

- **MinIO Console**: http://localhost:9001 (minioadmin/minioadmin)
- **NATS Monitoring**: http://localhost:8222
- **日志**: 所有服务输出 JSON 格式日志到 stdout，可使用 Loki/Promtail 采集

## 服务端口

| 服务 | gRPC 端口 | 说明 |
|------|-----------|------|
| gateway | 8080 (HTTP) | API 网关 |
| user-service | 50051 | 用户认证 |
| post-service | 50052 | 帖子管理 |
| comment-service | 50053 | 评论管理 |
| mood-service | 50054 | 情绪追踪 |
| quote-service | 50055 | 每日语录 |
| ai-service | 50056 | AI 暖心回复 |
| notification-service | 50057 | 推送通知 |
| postgres | 5432 | 数据库 |
| redis | 6379 | 缓存 |
| minio | 9000 | S3 存储 |
| nats | 4222 | 消息队列 |

## 故障排查

| 问题 | 排查命令 |
|------|----------|
| 服务无法连接数据库 | `docker compose ps` — 确认 postgres healthy |
| Gateway 路由 502 | 检查对应微服务是否启动 |
| MinIO 上传失败 | 确认 bucket 已创建: `mc ls local/showlove` |
| 迁移未执行 | 手动执行 SQL: `docker exec -i showlove-postgres psql -U showlove -d <db> < migration/*.up.sql` |
