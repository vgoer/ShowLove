# 实现计划：显出爱心 (Show Love) 后端 Golang 微服务

> 基于 SPEC.md v2026-06-17 · 前后端交替对接 · AI 直连真实 API

---

## 1. Context

**现状**：项目仅有 SPEC.md，无任何后端代码。
**目标**：从零搭建 7 个 Golang 微服务 + Gin 网关 + Flutter 前端适配。
**策略**：后端服务 → 网关路由 → Flutter 数据层 三者交替推进，每完成一个服务域立即前后端打通。
**约束**：Docker Compose MVP 部署，gRPC 服务间通信，JWT 认证，AI 直接接入 OpenAI + DeepSeek。

---

## 2. 依赖图

```
┌──────────────────────────────────────────────────────────────────┐
│              T00, T01, T02, T03 (基础设施并行)                     │
└────────────────────────────┬─────────────────────────────────────┘
                             │
              轮次1: 认证闭环
              ┌──────────────┴──────────────┐
              │ T04: user-service           │
              │ T05: Gateway骨架+Auth路由    │ ← 网关随服务增量生长
              │ T06: Flutter Auth对接       │
              └──────────────┬──────────────┘
                             │
              轮次2: 内容闭环
              ┌──────────────┴──────────────┐
              │ T07: storage 落地           │
              │ T08: post-service           │
              │ T09: Gateway Post路由       │
              │ T10: Flutter Feed+发帖      │
              │ T11: comment-service        │
              │ T12: Gateway Comment路由    │
              │ T13: Flutter 帖子详情+评论   │
              └──────────────┬──────────────┘
                             │
              轮次3: 智能闭环
              ┌──────────────┴──────────────┐
              │ T14: ai-service (真实API)   │
              │ T15: mood-service + 对接     │
              │ T16: quote-service + 对接    │
              │ T17: notification-service    │
              └──────────────┬──────────────┘
                             │
                             ▼
              ┌──────────────────────────────┐
              │ T18: E2E 测试 + CI/CD + 文档  │
              └──────────────────────────────┘
```

---

## 3. 任务列表（18 个任务，4 轮迭代）

### 📦 第 0 轮：基础设施（T00-T03，完全并行）

---

### T00 · 项目骨架与构建系统

**产出**：`backend/` 完整目录骨架、Makefile、Go workspace、golangci-lint 配置

**验收标准**：
- `make help` 列出全部命令
- `make lint` 通过（空项目）
- `go work` 正确引用所有子模块
- 目录结构符合 SPEC §3 定义

**涉及文件**：`Makefile`, `go.work`, `.env.example`, `.golangci.yml`, 7 个 `services/*/go.mod`

---

### T01 · Protocol Buffers 全量契约

**产出**：8 个 `.proto` 文件，一次性冻结前后端 + 服务间全部契约

| proto 文件 | RPC 方法 |
|-----------|----------|
| `user/user.proto` | Register, Login, RefreshToken, GetUser, UpdateProfile, UploadAvatar |
| `post/post.proto` | CreatePost, GetPost, ListPosts, DeletePost, SendSticker, ReportPost |
| `comment/comment.proto` | CreateComment, ListComments, DeleteComment |
| `mood/mood.proto` | RecordMood, GetMoods, GetWeeklyMood |
| `quote/quote.proto` | GetTodayQuote, CreateQuote, ListQuotes |
| `ai/ai.proto` | GenerateReply（内部用） |
| `notification/notif.proto` | RegisterDevice, SendPush |
| `common/common.proto` | 共享类型：Pagination, Timestamp, Error |

**验收标准**：
- `make proto` 生成 Go 代码零报错
- 所有分页统一 `page + page_size → total`
- 时间戳使用 `google.protobuf.Timestamp`
- 每个 message 字段与 SPEC §7 数据模型对应

**涉及文件**：`proto/**/*.proto`

---

### T02 · 共享包 pkg/

**产出**：6 个共享 Go 库

| 包 | 核心能力 | 关键测试 |
|----|---------|---------|
| `pkg/jwt` | RS256 签发/验证、Access(15min)+Refresh(7d) | 签发→验证→过期拒绝 |
| `pkg/storage` | Storage 接口 + MinIO 实现 + RustFS 占位 | MinIO testcontainer 上传/下载/删除 |
| `pkg/validator` | 邮箱/密码/昵称/内容长度校验 | 边界值 + 注入攻击 |
| `pkg/pagination` | 分页解析、Offset 计算 | 默认值/零值/超大值 |
| `pkg/errcode` | 统一错误码 + gRPC status 双向映射 | 每个 code 映射正确 |
| `pkg/events` | NATS JetStream 发布/订阅封装 | 发布→订阅→消费 |

**验收标准**：每个包 ≥80% 覆盖 + `go vet` + `golangci-lint` 全部通过

**涉及文件**：`pkg/*/`

---

### T03 · Docker Compose 基础设施

**产出**：`docker-compose.yml`，一键启动全部基础设施

```
postgres:15-alpine  → 6 个数据库自动初始化
redis:7.2-alpine    → AOF 持久化
minio:latest        → bucket 自动创建 + Console :9001
nats:2.10-alpine    → JetStream 启用
```

**验收标准**：
- `docker compose up -d` 30 秒内全部 healthy
- `docker compose ps` 全部 `(healthy)`
- MinIO Console `http://localhost:9001` 可访问
- `docker compose down -v` 后数据完全清理

**涉及文件**：`docker-compose.yml`, `scripts/init-multi-db.sh`

---

### 🔐 第 1 轮：认证闭环（T04 → T05 → T06，严格串行）

---

### T04 · user-service 完整实现

**产出**：用户微服务 —— 注册/登录/Token 刷新/资料管理/头像上传

```
services/user-service/
├── cmd/main.go
├── internal/
│   ├── server/grpc.go          # gRPC server 注册
│   ├── service/user.go         # 业务编排
│   ├── repository/user.go      # GORM 数据访问
│   ├── repository/token.go     # RefreshToken 管理
│   ├── model/user.go           # GORM model
│   └── service/user_test.go    # 单元测试
├── migration/
│   ├── 001_create_users.up.sql
│   └── 001_create_users.down.sql
├── config/config.yaml
└── Dockerfile
```

**实现要点**：bcrypt cost=12 · 邮箱唯一索引 · Refresh Token 存在数据库支持吊销 · 注册自动生成 UUID

**验收标准（gRPC 直连验证）**：
- Register → 成功 + 返回 User（无 password 字段）
- Register 重复邮箱 → AlreadyExists
- Login 正确密码 → access_token + refresh_token + user
- Login 错误密码 → Unauthenticated
- RefreshToken → 新 access_token
- RefreshToken 已吊销 → Unauthenticated
- GetUser/UpdateProfile → 正确返回/更新（需 token 中的 user_id）
- 单元测试 6 条路径全绿

**涉及文件**：约 10 个

---

### T05 · Gateway 骨架 + Auth 路由 + 中间件

**产出**：Gin 网关首批能力 —— 中间件栈 + 认证路由 + user-service gRPC 客户端

```
gateway/
├── cmd/main.go
├── internal/
│   ├── middleware/
│   │   ├── auth.go              # Bearer Token → user_id 注入 Context
│   │   ├── ratelimit.go         # 令牌桶（100 req/s）
│   │   ├── cors.go              # CORS *
│   │   └── logging.go           # trace_id + method + path + latency
│   ├── handler/auth_handler.go  # register/login/refresh
│   ├── handler/user_handler.go  # get/update me
│   ├── handler/health_handler.go
│   ├── handler/upload_handler.go
│   ├── client/user_client.go    # gRPC client (连接池, 5s timeout)
│   ├── client/storage_client.go
│   └── router/router.go         # 路由注册
├── config/config.yaml
└── Dockerfile
```

**路由表（首批）**：
```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh
GET    /api/v1/users/me            [需 Token]
PUT    /api/v1/users/me            [需 Token]
PUT    /api/v1/users/me/avatar     [需 Token]
POST   /api/v1/upload/image
GET    /api/v1/health
```

**验收标准**：
- 无 Token → 401 `{"code":401,"message":"missing authorization header"}`
- 过期 Token → 401 `{"code":401,"message":"token expired"}`
- 伪造 Token → 401 `{"code":401,"message":"invalid token"}`
- 正常请求 → user_id 注入 Context，日志含 trace_id
- `GET /health` → 200 + `{"gateway":"ok","user-service":"ok"}`

**涉及文件**：约 12 个

---

### T06 · Flutter 认证对接

**产出**：Flutter 前端数据层 Firebase → REST 替换（认证部分）

**变更范围**：
```
lib/core/constants/api_constants.dart    # baseUrl → 网关地址
lib/data/datasources/
├── api_datasource.dart                  # Dio 封装 (Base URL, 拦截器, 超时)
├── auth_datasource.dart                 # Firebase Auth → REST API
└── local_storage_datasource.dart        # Token 持久化
lib/data/repositories/auth_repository.dart  # 替换实现
lib/features/auth/                       # 页面无需大改
```

**Dio 拦截器核心逻辑**：
```
请求前 → 附加 Bearer Token
收到 401 → 自动调 /auth/refresh → 重试原请求
Refresh 也失败 → 清除 Token → emit 未认证状态
```

**验收标准（Flutter App 真机/模拟器验证）**：
- 注册新用户 → 自动登录 → 看到首页
- 关闭 App 重开 → Token 有效无需重新登录
- Token 过期 → 自动刷新 → 用户无感知
- 退出登录 → Token 清除 → 回到登录页
- 登录页/注册页 UI 与之前 Firebase 版一致

**涉及文件**：约 8 个

> ✅ **检查点 1**：用户可在 Flutter App 中完成注册→登录→查看/编辑个人资料。后端 Gateway + user-service + Flutter Auth 全链路打通。

---

### 📝 第 2 轮：内容闭环（T07-T13，后端服务与前端页面交替推进）

---

### T07 · 存储层完整落地

**产出**：MinIO 生产可用 + RustFS 适配器 + 图片上传全链路

**验收标准**：
- `POST /api/v1/upload/image` → 201 + `{"url":"http://minio:9000/showlove/images/xxx.jpg"}`
- 上传后通过 URL 可访问图片
- 删除后 URL 返回 404
- 环境变量 `STORAGE_PROVIDER=minio|rustfs` 切换正常
- 图片自动生成缩略图（200×200）

**涉及文件**：`pkg/storage/minio.go`（补全）, `pkg/storage/rustfs.go`（新）, `pkg/storage/thumbnail.go`（新）

---

### T08 · post-service 完整实现

**产出**：帖子微服务 —— CRUD + 列表排序分页 + 贴纸计数 + 敏感词过滤 + 事件发布

```
services/post-service/
├── cmd/main.go
├── internal/
│   ├── server/grpc.go
│   ├── service/post.go         # 创建/查询/列表/删除/举报
│   ├── service/sticker.go      # 贴纸增量计数
│   ├── moderation/filter.go    # AC 自动机 + 中文敏感词库
│   ├── moderation/filter_test.go
│   ├── repository/post.go
│   ├── model/post.go
│   └── service/post_test.go
├── migration/001_create_posts.sql
├── config/config.yaml
└── Dockerfile
```

**关键设计**：
- 创建帖子 → 发布 `post.created` 事件到 NATS（ai-service 消费）
- 列表排序：`sort=latest`（created_at DESC）/ `most_helped`（comment_count DESC）
- 敏感词命中 → 拒绝创建，返回敏感词列表

**验收标准（gRPC 直连验证）**：
- CreatePost → 201 + Post（含完整字段）
- ListPosts(sort=latest) → 分页正确，按时间倒序
- ListPosts(sort=most_helped) → 按评论数倒序
- GetPost → 200 + 含作者昵称/头像（冗余缓存字段）
- DeletePost → 仅作者可删，他人 403
- SendSticker(hug) → 200 + sticker_hug 递增
- 含敏感词内容 → InvalidArgument + 敏感词提示
- 分页边界：page=0 → 默认第1页，超大page → 空列表

**涉及文件**：约 12 个

---

### T09 · Gateway Post 路由

**产出**：Gateway 新增 post handler + post-service gRPC client

**新增路由**：
```
GET    /api/v1/posts                   [需 Token]
POST   /api/v1/posts                   [需 Token，multipart]
GET    /api/v1/posts/:id               [需 Token]
DELETE /api/v1/posts/:id               [需 Token]
POST   /api/v1/posts/:id/stickers      [需 Token]
POST   /api/v1/posts/:id/report        [需 Token]
```

**验收标准**：curl 可完成 创建→列表→详情→贴纸→举报 完整流程

**涉及文件**：`gateway/internal/handler/post_handler.go`, `gateway/internal/client/post_client.go`, `gateway/internal/router/router.go`（追加）

---

### T10 · Flutter 帖子流 + 发帖对接

**产出**：Flutter feed 页面 + create_post 页面数据层对接

**变更范围**：
- `lib/data/datasources/post_datasource.dart`：Firestore → REST
- `lib/data/repositories/post_repository.dart`：替换实现
- `lib/features/feed/`：适配新的 Post 模型
- `lib/features/create_post/`：图片上传 + 发帖流程

**验收标准（Flutter App 验证）**：
- 首页展示帖子卡片列表（按最新排序）
- 下拉刷新加载新帖子
- 上拉加载更多（分页）
- 发布新帖子（文字+图片+心情标签）→ 出现在列表顶部
- 帖子卡片骨架屏 + 加载/错误/空状态三态展示

**涉及文件**：约 8 个

---

### T11 · comment-service 完整实现

**产出**：评论微服务 —— CRUD + 贴纸联动

**验收标准（gRPC 直连验证）**：
- CreateComment → 201 + 帖子 comment_count 自动 +1
- ListComments → 按时间正序
- DeleteComment → 仅作者可删
- is_ai_generated 标记正确持久化

**涉及文件**：约 8 个

---

### T12 · Gateway Comment 路由

**产出**：Gateway 新增 comment handler + client

**新增路由**：
```
GET    /api/v1/posts/:id/comments      [需 Token]
POST   /api/v1/posts/:id/comments      [需 Token]
DELETE /api/v1/comments/:id            [需 Token]
```

**涉及文件**：`gateway/internal/handler/comment_handler.go`, `gateway/internal/client/comment_client.go`, router 追加

---

### T13 · Flutter 帖子详情 + 评论对接

**产出**：Flutter post_detail 页面数据层对接

**验收标准（Flutter App 验证）**：
- 点击帖子卡片 → 进入详情页 → 显示完整帖子内容
- 评论列表正常加载
- 发送评论 → 立即出现在列表底部
- 发送暖心贴纸 → 计数动画更新
- AI 回复横幅展示（如果已生成）

**涉及文件**：约 6 个

> ✅ **检查点 2**：用户可在 App 中浏览帖子列表 → 查看详情 → 发表评论 → 发送贴纸。帖子+评论全链路打通。

---

### 🤖 第 3 轮：智能闭环（T14-T17，大部分可并行）

---

### T14 · ai-service（直连 OpenAI + DeepSeek）

**产出**：AI 服务 —— 监听帖子创建事件 → 调用真实 AI → 写入 AI 评论

```
services/ai-service/
├── cmd/main.go
├── internal/
│   ├── server/grpc.go
│   ├── provider/
│   │   ├── provider.go          # interface{ Generate(ctx, prompt) (string, error) }
│   │   ├── openai.go            # GPT-4o-mini (默认)
│   │   ├── deepseek.go          # DeepSeek
│   │   ├── qwen.go              # 通义千问（占位）
│   │   └── ollama.go            # 本地 Ollama（占位）
│   ├── prompt/templates.go      # 中英双语治愈 Prompt
│   ├── service/ai.go            # 消费 NATS → 调 AI → gRPC 评论
│   └── service/ai_test.go       # 用 mock provider 测试流程
├── config/config.yaml
└── Dockerfile
```

**实现要点**：
- 启动时验证 AI_PROVIDER 配置的 API Key 是否有效
- OpenAI 和 DeepSeek 都完整实现，其他两个占位
- 降级策略：AI 调用失败 → 记录日志 + 不阻塞帖子创建
- Prompt 包含「以治愈温暖的语气回复，50-150 字，署名小暖」

**验收标准**：
- 创建帖子 → 60 秒内评论区出现「小暖」的 AI 回复
- `AI_PROVIDER=openai` → 调 GPT-4o-mini，回复内容温暖治愈
- `AI_PROVIDER=deepseek` → 调 DeepSeek API，正常返回
- AI API 不可用 → 帖子创建正常，日志记录错误
- 并发创建 5 个帖子 → 每个都能收到 AI 回复

**涉及文件**：约 10 个

---

### T15 · mood-service + Gateway + Flutter 情绪追踪

**后端**：mood-service（gRPC）+ Gateway mood 路由
```
POST   /api/v1/moods                       [需 Token]
GET    /api/v1/moods?from=&to=             [需 Token]
GET    /api/v1/moods/weekly                [需 Token]
```

**前端**：Flutter mood_tracker 页面 Firebase → REST 替换

**验收标准**：
- 记录今日情绪（选温度计 1-10 + 标签）→ 成功
- 同一天再次记录 → 更新覆盖
- 周曲线页面 → 显示最近 7 天情绪折线图
- 没有记录的日期 → 曲线断点（不画0）

**涉及文件**：`mood-service/`（约 8 个）+ Gateway 增加 3 个文件 + Flutter 约 4 个文件

---

### T16 · quote-service + Gateway + Flutter 每日语录

**后端**：quote-service（gRPC + 定时任务）+ Gateway quote 路由
```
GET    /api/v1/quotes/today                [需 Token]
```

**种子数据**：30 条中英双语治愈语录，首次启动自动写入

**前端**：Flutter daily_quote 卡片 Firebase → REST 替换

**验收标准**：
- `GET /quotes/today` → 返回当天语录（含中英文+作者+配图）
- cron 在 08:00 CST 触发 → 发布事件到 NATS
- Flutter 首页每日语录卡片正常展示

**涉及文件**：`quote-service/`（约 10 个）+ Gateway 增加 3 个文件 + Flutter 约 2 个文件

---

### T17 · notification-service + Gateway 设备注册

**后端**：notification-service（gRPC + NATS 事件消费 + FCM 推送）
```
POST   /api/v1/devices                    [需 Token]
```

**监听事件**：`post.commented` → 推送给帖子作者 · `daily_quote.push` → 推送给全部设备

**验收标准**：
- 注册设备 Token → 201
- 重复 Token → 去重（upsert）
- 某人评论我的帖子 → 我收到推送「XX 评论了你的帖子」
- 推送失败 → 重试 3 次 → 仍失败记录日志

**涉及文件**：`notification-service/`（约 8 个）+ Gateway 增加 3 个文件

> ✅ **检查点 3**：全部 7 个微服务 + Gateway + Flutter 全链路可用。8 个用户故事全部可走通。

---

### ✅ 第 4 轮：质量保证

---

### T18 · E2E 测试 + CI/CD + 文档

**产出**：
- `backend/test/e2e/` — 3 条核心 API 全链路测试
- `.github/workflows/backend-ci.yml` + `frontend-ci.yml`
- `docs/api-reference.md` + `docs/deployment.md`

**E2E 测试三条核心流程**：
```
1. 注册→登录→获取个人信息→更新昵称（认证全链路）
2. 上传图片→创建帖子→列表查询→查看详情→发评论→发贴纸→验证 AI 回复出现
3. 记录情绪→同天更新→查看周曲线
```

**CI 流程**：
- Backend: checkout → docker compose up → migrate → test → lint → docker compose down
- Frontend: checkout → flutter pub get → flutter analyze → flutter test

**验收标准**：
- `make test-e2e`：全自动启动环境 → 运行 3 条 → 全部通过 → 清理
- CI push 触发自动运行，lint + test 全绿才能 merge
- API 文档：每个端点含 curl 示例 + 响应示例
- 部署指南：从 git clone 到完整运行，步骤 ≤ 10 步

**涉及文件**：约 10 个

---

## 4. 执行顺序

```
                     T00 ─┬─ T01 ─┬─ T02 ─┬─ T03
                          │       │       │
                          ▼       ▼       ▼
                     (4 个基础任务可完全并行，预计 1-2 天)

轮次1 ─── T04 ──→ T05 ──→ T06        (串行, 3-4 天)
轮次2 ─── T07 ──→ T08 ──→ T09 ──→ T10 ──→ T11 ──→ T12 ──→ T13
              │                    │
              └── 2-3 天 ──────────┘               (串行, 4-5 天)
轮次3 ─── T14 ──┬── T15 ──┬── T16 ──┬── T17
                │         │         │
                └─ 并行 ──┘         │            (并行, 3-4 天)
轮次4 ─── T18                                      (2 天)
```

### 严格串行依赖

| 依赖链 | 原因 |
|--------|------|
| T04 → T05 → T06 | user-service 就绪后 Gateway 才能加 Auth 路由，Gateway 好了 Flutter 才能对接 |
| T08 → T09 → T10 | post-service 就绪后 Gateway 才能加 Post 路由，Flutter 才能对接 feed |
| T11 → T12 → T13 | comment-service 就绪后 Gateway + Flutter 才能对接评论 |
| T08 → T14 | ai-service 消费 post.created 事件，需要 post-service 先就绪 |

### 可并行任务

| 并行组 | 任务 | 理由 |
|--------|------|------|
| 组1 | T00, T01, T02, T03 | 互不依赖 |
| 组2 | T07 可与 T04-T06 并行 | 存储不依赖认证 |
| 组3 | T14, T15, T16, T17 | 各自独立：AI监听事件、情绪无依赖、语录无依赖、通知无依赖 |

---

## 5. 验证检查点

| CP | 位置 | 验证方式 |
|----|------|----------|
| CP1 🔐 | T06 后 | Flutter App 完成 注册→登录→浏览资料。Gateway Auth 全部端点可用 curl 调通 |
| CP2 📝 | T13 后 | Flutter App 完成 浏览帖子列表→发帖→评论→贴纸。内容全链路可用 |
| CP3 🤖 | T17 后 | 全部 8 个用户故事可在 Flutter App 走通，7 个微服务健康检查全绿 |
| CP4 ✅ | T18 后 | E2E 全绿 + CI 流水线正常 + 文档完整 |

---

## 6. 风险与缓解

| 风险 | 缓解 |
|------|------|
| protobuf 契约频繁变更 | T01 冻结全部接口定义，后续变更走 PR review |
| AI API 不稳定/欠费 | ai-service 内部降级：失败不阻塞帖子创建；DeepSeek 作为 OpenAI 的 fallback |
| 前后端联调耗时 | 每个服务完成后立即对接 Flutter，问题早发现 |
| Docker 资源不足 | 每个服务可独立 `go run` 开发，不必须全量启动 Docker |
