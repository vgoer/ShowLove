# Spec: 显出爱心 (Show Love) — 治愈系社区互助平台

> **类型**: 创业 MVP · **阶段**: 规格确认完毕 · **最后更新**: 2026-06-17
> **架构**: Golang 微服务 + Flutter 前端

---

## 1. Objective — 我们为什么做这个

### 1.1 产品愿景
一个以「解忧杂货铺」为精神原型的治愈系社区——每个人都可以在这里卸下防备，说出自己的不开心或困难，
而社区里的人们用温暖的文字、贴心的互动来帮忙想办法。让每一个困境都有人看见，每一声叹息都有人回应。

### 1.2 目标用户

| 维度 | 定义 |
|------|------|
| 年龄 | 18~35 岁年轻人 |
| 地域 | 全球（初期中文 + 英文双语） |
| 场景 | 情绪低落、遇到困难、想倾诉但不想让熟人知道 |
| 心态 | 愿意接受温暖，也愿意给予温暖 |

### 1.3 用户故事（核心 MVP）

```
US-01 注册/登录
  作为一个新用户，我可以用邮箱注册账号，设置昵称和头像，以便在社区中拥有自己的身份。

US-02 发布帖子
  作为一个遇到困难的用户，我可以发布一篇帖子，写下我的困境（文字/图片/语音），
  选择一个心情标签，让其他人看到并提供帮助。

US-03 浏览帖子
  作为一个想帮助别人的用户，我可以浏览首页推荐的帖子卡片，
  按"最新"或"最需要帮助"排序，看到每个人正在经历什么。

US-04 评论/互动
  作为一个想支持他人的用户，我可以在帖子下留下温暖的评论，
  或者发送"暖心贴纸"（抱抱、加油、我懂你…）快捷表达关心。

US-05 情绪记录
  作为一个关注自我情绪的用户，我可以用"情绪温度计"记录每天的心情变化，
  看到一周的情绪曲线。

US-06 AI 暖心回复
  作为发帖用户，当帖子暂时没有人回复时，系统会生成一条 AI 暖心鼓励，
  确保没有人感到被遗忘。

US-07 每日语录推送
  作为一个需要温暖提醒的用户，我每天可以收到一条治愈语录推送。

US-08 内容安全
  作为一个普通用户，我可以举报不当内容；
  系统会自动过滤敏感词，确保社区温暖安全。
```

### 1.4 成功标准（可验证）

| # | 标准 | 测量方式 |
|---|------|----------|
| S1 | 用户能在 3 步内完成注册并发布第一篇帖子 | 可用性测试 |
| S2 | 新帖子在无人回复时，60 秒内收到 AI 回复 | 自动化测试 |
| S3 | 敏感词过滤准确率 ≥ 95% | 测试集验证 |
| S4 | 首页帖子列表 API 响应时间 < 200ms (p95) | 性能测试 |
| S5 | 所有微服务健康检查通过率 ≥ 99.9% | 监控告警 |
| S6 | 治愈风格一致性：所有页面通过 UI 审查清单 | 人工审查 |
| S7 | 核心 API 错误率 < 0.5% | 日志/指标监控 |

---

## 2. Tech Stack — 技术选型

### 2.1 后端技术栈

| 层级 | 技术 | 版本 | 选择理由 |
|------|------|------|----------|
| 语言 | Go | ≥ 1.22 | 高性能，并发原生支持，微服务生态成熟 |
| HTTP 框架 | Gin | ≥ 1.10 | 高性能 HTTP 路由，中间件生态丰富 |
| RPC 框架 | gRPC + Protocol Buffers | ≥ 1.64 | 服务间高性能通信，强类型契约 |
| 数据库 | PostgreSQL | ≥ 15 | 成熟的关系型数据库，JSONB 支持灵活查询 |
| 缓存 | Redis | ≥ 7.2 | 会话缓存 + 热数据加速 + 分布式锁 |
| 对象存储 | MinIO (S3-compatible) | latest | 自建 S3 兼容存储，图片/语音文件 |
| 文件存储 | RustFS (SeaweedFS) | latest | 分布式小文件存储，高可用 |
| 消息队列 | NATS / Redis Streams | latest | 轻量级异步消息，事件驱动 |
| 认证 | JWT + Refresh Token | — | 无状态认证，适合微服务架构 |
| ORM | GORM | ≥ 1.25 | Go 生态最成熟的 ORM |
| 配置管理 | Viper | ≥ 1.18 | 多格式配置，环境变量注入 |
| 日志 | Zerolog | ≥ 1.32 | 高性能零分配 JSON 日志 |
| 链路追踪 | OpenTelemetry | ≥ 1.24 | 分布式追踪标准 |
| 容器化 | Docker + Docker Compose | latest | MVP 阶段部署方案 |

### 2.2 前端技术栈（不变）

| 层级 | 技术 | 版本 | 选择理由 |
|------|------|------|----------|
| 前端框架 | Flutter | ≥ 3.27 | 单代码库覆盖 iOS/Android/Web |
| 语言 | Dart | ≥ 3.6 | — |
| 状态管理 | Riverpod | ≥ 2.6 | 编译时安全，测试友好 |
| 路由 | go_router | ≥ 14 | Flutter 官方推荐 |
| HTTP 客户端 | Dio | ≥ 5.4 | 拦截器、重试、文件上传 |
| 崩溃监控 | Sentry | latest | 跨平台错误追踪 |

### 2.3 AI 模型支持（多模型可切换）

| 模型 | 用途 | 配置方式 |
|------|------|----------|
| OpenAI GPT-4o-mini | 暖心回复（默认） | 环境变量 `AI_PROVIDER=openai` |
| DeepSeek | 暖心回复（备选） | 环境变量 `AI_PROVIDER=deepseek` |
| 通义千问 | 暖心回复（国内） | 环境变量 `AI_PROVIDER=qwen` |
| 本地模型 (Ollama) | 离线部署 | 环境变量 `AI_PROVIDER=ollama` |

### 2.4 最低平台版本

| 平台 | 最低版本 | 覆盖率 |
|------|----------|--------|
| iOS | 15.0+ | ~96% |
| Android | API 26 (Android 8.0) | ~95% |

---

## 3. Project Structure — 项目结构

```
show_love/
├── backend/                            # Golang 后端单体仓库（monorepo）
│   ├── gateway/                        # API 网关服务（Gin）
│   │   ├── cmd/
│   │   │   └── main.go                 # 网关入口
│   │   ├── internal/
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go             # JWT 验证中间件
│   │   │   │   ├── ratelimit.go        # 限流中间件
│   │   │   │   ├── cors.go             # CORS 中间件
│   │   │   │   └── logging.go          # 请求日志中间件
│   │   │   ├── handler/                # HTTP 处理器（路由 → gRPC 调用）
│   │   │   │   ├── auth_handler.go
│   │   │   │   ├── post_handler.go
│   │   │   │   ├── comment_handler.go
│   │   │   │   ├── mood_handler.go
│   │   │   │   ├── quote_handler.go
│   │   │   │   └── upload_handler.go
│   │   │   ├── router/
│   │   │   │   └── router.go           # 路由注册
│   │   │   └── client/                 # 下游 gRPC 客户端
│   │   │       ├── user_client.go
│   │   │       ├── post_client.go
│   │   │       ├── comment_client.go
│   │   │       ├── mood_client.go
│   │   │       └── ai_client.go
│   │   ├── config/
│   │   │   └── config.yaml             # 网关配置
│   │   └── Dockerfile
│   │
│   ├── services/                       # 微服务
│   │   ├── user-service/               # 用户服务
│   │   │   ├── cmd/main.go
│   │   │   ├── internal/
│   │   │   │   ├── server/             # gRPC 服务端实现
│   │   │   │   ├── repository/         # 数据访问层（GORM + PostgreSQL）
│   │   │   │   ├── model/              # 数据模型
│   │   │   │   └── service/            # 业务逻辑层
│   │   │   ├── migration/              # 数据库迁移脚本
│   │   │   ├── config/
│   │   │   └── Dockerfile
│   │   │
│   │   ├── post-service/               # 帖子服务
│   │   │   ├── cmd/main.go
│   │   │   ├── internal/
│   │   │   │   ├── server/
│   │   │   │   ├── repository/
│   │   │   │   ├── model/
│   │   │   │   ├── service/
│   │   │   │   └── moderation/         # 敏感词过滤
│   │   │   ├── migration/
│   │   │   ├── config/
│   │   │   └── Dockerfile
│   │   │
│   │   ├── comment-service/            # 评论服务
│   │   │   ├── cmd/main.go
│   │   │   ├── internal/
│   │   │   │   ├── server/
│   │   │   │   ├── repository/
│   │   │   │   ├── model/
│   │   │   │   └── service/
│   │   │   ├── migration/
│   │   │   ├── config/
│   │   │   └── Dockerfile
│   │   │
│   │   ├── mood-service/               # 情绪记录服务
│   │   │   ├── cmd/main.go
│   │   │   ├── internal/
│   │   │   │   ├── server/
│   │   │   │   ├── repository/
│   │   │   │   ├── model/
│   │   │   │   └── service/
│   │   │   ├── migration/
│   │   │   ├── config/
│   │   │   └── Dockerfile
│   │   │
│   │   ├── quote-service/              # 每日语录服务
│   │   │   ├── cmd/main.go
│   │   │   ├── internal/
│   │   │   │   ├── server/
│   │   │   │   ├── repository/
│   │   │   │   ├── model/
│   │   │   │   ├── service/
│   │   │   │   └── scheduler/          # 定时推送调度
│   │   │   ├── migration/
│   │   │   ├── config/
│   │   │   └── Dockerfile
│   │   │
│   │   ├── ai-service/                 # AI 暖心回复服务
│   │   │   ├── cmd/main.go
│   │   │   ├── internal/
│   │   │   │   ├── server/
│   │   │   │   ├── provider/           # 多模型适配器
│   │   │   │   │   ├── provider.go    # 接口定义
│   │   │   │   │   ├── openai.go
│   │   │   │   │   ├── deepseek.go
│   │   │   │   │   ├── qwen.go
│   │   │   │   │   └── ollama.go
│   │   │   │   ├── prompt/             # Prompt 模板
│   │   │   │   └── service/
│   │   │   ├── config/
│   │   │   └── Dockerfile
│   │   │
│   │   └── notification-service/       # 通知推送服务
│   │       ├── cmd/main.go
│   │       ├── internal/
│   │       │   ├── server/
│   │       │   ├── push/               # FCM / APNs 推送适配
│   │       │   └── service/
│   │       ├── config/
│   │       └── Dockerfile
│   │
│   ├── pkg/                            # 共享库
│   │   ├── jwt/                         # JWT 生成/验证
│   │   ├── storage/                     # 存储抽象层
│   │   │   ├── storage.go              # 存储接口
│   │   │   ├── minio.go               # MinIO 实现
│   │   │   └── rustfs.go              # RustFS/SeaweedFS 实现
│   │   ├── validator/                   # 通用验证工具
│   │   ├── pagination/                  # 分页工具
│   │   ├── errcode/                     # 统一错误码
│   │   └── events/                      # 事件定义（跨服务消息）
│   │
│   ├── proto/                           # Protocol Buffers 定义
│   │   ├── user/
│   │   │   └── user.proto
│   │   ├── post/
│   │   │   └── post.proto
│   │   ├── comment/
│   │   │   └── comment.proto
│   │   ├── mood/
│   │   │   └── mood.proto
│   │   ├── quote/
│   │   │   └── quote.proto
│   │   ├── ai/
│   │   │   └── ai.proto
│   │   └── notification/
│   │       └── notification.proto
│   │
│   ├── docker-compose.yml               # 本地开发编排
│   ├── docker-compose.prod.yml          # 生产环境编排
│   ├── Makefile                         # 统一构建命令
│   └── .env.example                     # 环境变量模板
│
├── frontend/                            # Flutter 前端
│   ├── android/
│   ├── ios/
│   ├── web/
│   ├── assets/
│   │   ├── images/
│   │   ├── animations/
│   │   └── fonts/
│   ├── lib/
│   │   ├── main.dart
│   │   ├── app.dart
│   │   ├── core/                        # 全局基础设施
│   │   │   ├── constants/
│   │   │   │   ├── app_colors.dart
│   │   │   │   ├── app_theme.dart
│   │   │   │   ├── app_text_styles.dart
│   │   │   │   └── api_constants.dart   # 后端 API 端点（改为网关地址）
│   │   │   ├── extensions/
│   │   │   ├── utils/
│   │   │   ├── widgets/                 # 全局共享组件
│   │   │   ├── router/
│   │   │   │   └── app_router.dart
│   │   │   └── l10n/
│   │   │       ├── app_zh.arb
│   │   │       └── app_en.arb
│   │   ├── data/                        # 数据层（调用 REST API）
│   │   │   ├── repositories/
│   │   │   │   ├── auth_repository.dart
│   │   │   │   ├── post_repository.dart
│   │   │   │   ├── mood_repository.dart
│   │   │   │   └── quote_repository.dart
│   │   │   └── datasources/
│   │   │       ├── api_datasource.dart           # Dio HTTP 客户端封装
│   │   │       ├── auth_datasource.dart
│   │   │       ├── post_datasource.dart
│   │   │       └── local_storage_datasource.dart  # Token 持久化
│   │   ├── domain/                      # 领域层（纯 Dart）
│   │   │   ├── models/
│   │   │   └── repositories/            # 抽象接口
│   │   └── features/                    # 功能模块
│   │       ├── auth/
│   │       ├── feed/
│   │       ├── post_detail/
│   │       ├── create_post/
│   │       ├── mood_tracker/
│   │       ├── daily_quote/
│   │       ├── profile/
│   │       └── settings/
│   ├── test/
│   │   ├── unit/
│   │   ├── widget/
│   │   └── integration_test/
│   ├── pubspec.yaml
│   └── analysis_options.yaml
│
├── docs/                                # 项目文档
│   ├── architecture.md                  # 架构说明
│   ├── api-reference.md                 # API 文档
│   └── deployment.md                    # 部署指南
│
├── scripts/                             # 工具脚本
│   ├── init-db.sh                       # 初始化数据库
│   └── gen-proto.sh                     # 生成 protobuf 代码
│
├── .github/
│   └── workflows/                       # CI/CD
│       ├── backend-ci.yml
│       └── frontend-ci.yml
│
└── SPEC.md                              # ← 本文件
```

---

## 4. Commands — 开发命令

### 4.1 基础设施

```bash
# === Docker 环境启动 ===
docker compose up -d                        # 启动所有基础设施（PG、Redis、MinIO、RustFS）
docker compose ps                           # 查看服务状态
docker compose down -v                      # 停止并清理数据卷

# === protobuf 代码生成 ===
cd backend && make proto                    # 从 proto/ 生成所有 Go 和 Dart gRPC 代码

# === 数据库迁移 ===
cd backend && make migrate-up              # 运行所有服务的数据库迁移
cd backend && make migrate-down            # 回滚所有迁移
```

### 4.2 后端开发

```bash
# === 构建 ===
cd backend && make build                    # 构建所有微服务
cd backend && make build-gateway           # 仅构建网关
cd backend && make build-user             # 仅构建用户服务

# === 运行 ===
cd backend && make run-gateway             # 启动网关（开发模式，热重载用 air）
cd backend && make run-user               # 启动用户服务
cd backend && make run-all                 # 启动所有服务（docker compose up）

# === 静态检查 ===
cd backend && make lint                    # golangci-lint run（CI 强制通过）
cd backend && make fmt                     # gofmt + goimports 格式化
cd backend && make vet                     # go vet 检查

# === 测试 ===
cd backend && make test                    # 运行全部后端测试
cd backend && make test-unit              # 仅单元测试
cd backend && make test-integration       # 仅集成测试（需 Docker 环境）
cd backend && make test-coverage          # 带覆盖率报告

# === 单个服务测试 ===
cd backend/services/user-service && go test ./... -v
cd backend/services/post-service && go test ./... -v
```

### 4.3 前端开发（不变）

```bash
# === 环境初始化 ===
cd frontend && flutter clean && flutter pub get
cd frontend && dart run build_runner build --delete-conflicting-outputs

# === 开发运行 ===
cd frontend && flutter run
cd frontend && flutter run -d chrome       # Web 预览

# === 静态检查 ===
cd frontend && flutter analyze             # Dart 静态分析（CI 强制通过）
cd frontend && dart format --set-exit-if-changed lib/ test/

# === 测试 ===
cd frontend && flutter test
cd frontend && flutter test --coverage

# === 构建 ===
cd frontend && flutter build apk --flavor prod
cd frontend && flutter build ios --flavor prod
cd frontend && flutter build web --flavor prod
```

---

## 5. Code Style — 代码风格

### 5.1 Go 代码风格

```go
// ✅ 文件命名: snake_case
// user_service.go, post_handler.go, jwt_middleware.go

// ✅ 包命名: 小写单词，不用下划线
package repository  // ✅
package user_repo   // ❌

// ✅ 接口: 单方法接口 + er 后缀；多方法接口名词
type Reader interface {
    Read(ctx context.Context, id string) (*User, error)
}

type UserRepository interface {
    FindByID(ctx context.Context, id string) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    Create(ctx context.Context, user *User) error
    Update(ctx context.Context, user *User) error
}

// ✅ 错误处理: 始终检查，包装上下文
func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("user service: get user %s: %w", id, err)
    }
    if user == nil {
        return nil, ErrNotFound
    }
    return user, nil
}

// ✅ 构造函数: NewXxx 返回接口，接受依赖
func NewUserService(repo UserRepository, cache *redis.Client) *UserService {
    return &UserService{repo: repo, cache: cache}
}

// ✅ gRPC 服务端: 小而专注的方法
func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    user, err := s.svc.GetUser(ctx, req.GetUserId())
    if err != nil {
        return nil, status.Errorf(codeFromError(err), "get user: %v", err)
    }
    return &pb.GetUserResponse{User: userToProto(user)}, nil
}

// ✅ Gin handler: 简洁，参数校验前置
func (h *PostHandler) CreatePost(c *gin.Context) {
    var req CreatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, ErrResponse(errcode.InvalidParam, err.Error()))
        return
    }
    userID := c.GetString("user_id") // 从 JWT 中间件注入

    resp, err := h.client.CreatePost(c.Request.Context(), &pb.CreatePostRequest{
        AuthorId: userID,
        Content:  req.Content,
        MoodTag:  req.MoodTag,
    })
    if err != nil {
        c.JSON(httpStatusFromGRPC(err), ErrResponse(errcode.Internal, err.Error()))
        return
    }
    c.JSON(201, OKResponse(resp.GetPost()))
}
```

### 5.2 关键约定（Go）

| 规则 | 说明 |
|------|------|
| Context 传递 | 所有函数第一个参数是 `ctx context.Context` |
| 错误包装 | 使用 `fmt.Errorf("...: %w", err)` 保留调用链 |
| 零值依赖 | 不使用 `init()` 做初始化，依赖注入在 `main.go` 中 |
| interface 定义在使用方 | 接口定义在调用方包，不是实现方 |
| 一行 ≤ 120 字符 | `gofumpt` 自动格式化 |
| 禁止 `panic` | 除非 `main.go` 初始化致命错误 |
| 数据库事务 | 通过 `context` 传递事务对象 |
| 日志 | 使用 `zerolog`，结构化日志字段用 `Str()/Int()/Err()` |

### 5.3 Dart/Flutter 代码风格（保持不变）

```dart
// ✅ 好的写法: ConsumerWidget + const 构造 + 小而专注
class PostCard extends ConsumerWidget {
  const PostCard({required this.post, super.key});

  final Post post;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Padding(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _PostHeader(post: post),
          const Gap(12),
          _PostBody(post: post),
          const Gap(12),
          _PostActions(post: post),
        ],
      ),
    );
  }
}
```

**Dart 关键约定：**

| 规则 | 说明 |
|------|------|
| `const` 优先 | 能用 `const` 的地方一律用 |
| Widget 拆分 | 单个 `build` 方法不超过 50 行 |
| 模型不可变 | 使用 `freezed` / `@immutable` |
| 异步状态 | 使用 `AsyncValue.guard` 统一错误处理 |
| API 调用 | 必须走 Repository 层，禁止 Widget 直接调 HTTP |

---

## 6. Microservices Architecture — 微服务架构

### 6.1 服务全景图

```
                              ┌─────────────────────────────────────┐
                              │        Flutter App (Mobile/Web)      │
                              └─────────────────┬───────────────────┘
                                                │ REST (HTTPS + JWT)
                                                ▼
                              ┌─────────────────────────────────────┐
                              │          API Gateway (Gin)           │
                              │   Auth · RateLimit · CORS · Log     │
                              └──────┬──────┬──────┬──────┬────────┘
                                     │      │      │      │
                           gRPC ────┼──────┼──────┼──────┼──── gRPC
                                     ▼      ▼      ▼      ▼
                              ┌──────┴──────┴──────┴──────┴────────┐
                              │         Microservices Layer          │
                              │                                      │
                    ┌─────────┴─────────┐  ┌─────────┐  ┌─────────┐
                    │   user-service    │  │  post-   │  │comment- │
                    │   (认证/用户/资料) │  │ service  │  │ service │
                    └─────────┬─────────┘  └────┬─────┘  └────┬────┘
                              │               │             │
                    ┌─────────┴─────────┐  ┌──┴─────┐  ┌──┴─────────┐
                    │   mood-service    │  │ quote- │  │    ai-     │
                    │   (情绪追踪)       │  │service │  │  service   │
                    └─────────┬─────────┘  └────┬────┘  └─────┬──────┘
                              │               │             │
                    ┌─────────┴─────────┐       │             │
                    │notification-svc   │       │             │
                    │ (推送通知)         │       │             │
                    └───────────────────┘       │             │
                                                 │             │
                              ┌──────────────────┴─────────────┴────┐
                              │            Infrastructure            │
                              │  PostgreSQL · Redis · NATS/Redis     │
                              │  MinIO / RustFS · OpenTelemetry      │
                              └──────────────────────────────────────┘
```

### 6.2 服务职责边界

| 服务 | 职责 | 数据库 | 暴露 API |
|------|------|--------|----------|
| **gateway** | 路由转发、JWT 验证、限流、CORS | 无 | REST（对外） |
| **user-service** | 注册、登录、Token 管理、用户资料、善意积分 | `users_db` | gRPC（对内） |
| **post-service** | 帖子 CRUD、敏感词过滤、帖子列表、贴纸计数 | `posts_db` | gRPC |
| **comment-service** | 评论 CRUD、暖心贴纸发送 | `comments_db` | gRPC |
| **mood-service** | 情绪记录、温度计、周曲线统计 | `moods_db` | gRPC |
| **quote-service** | 每日语录管理、定时任务调度 | `quotes_db` | gRPC |
| **ai-service** | AI 回复生成、多模型适配、Prompt 管理 | 无（调用外部 API） | gRPC |
| **notification-service** | 推送通知（FCM/APNs）、设备 Token 管理 | `notifications_db` | gRPC |

### 6.3 关键业务流程

```
【用户发帖 → AI 回复流程】
User ──POST /posts──▶ Gateway ──gRPC──▶ post-service
                                            │
                                    保存帖子到 PostgreSQL
                                            │
                                    发送事件到 Redis Streams
                                            │
                                     ┌──────▼──────┐
                                     │  ai-service  │ 消费事件
                                     │              │ 生成 AI 回复
                                     └──────┬──────┘
                                            │ gRPC
                                     ┌──────▼────────┐
                                     │comment-service │ 保存 AI 评论
                                     └──────┬────────┘
                                            │ 发送推送事件
                                     ┌──────▼────────────┐
                                     │notification-service│
                                     └───────────────────┘

【跨服务数据查询策略】
- 原则: 每个服务拥有自己的数据，不跨库 JOIN
- 帖子卡片展示作者昵称/头像 → post-service 本地缓存一份用户基础信息（通过事件同步）
- 评论列表展示评论者信息 → comment-service 本地缓存
- 数据同步: 用户更新资料时，user-service 发布事件，其他服务消费更新本地缓存
```

### 6.4 API 路由设计（Gateway 对外）

```
# === 认证（无需 Token）===
POST   /api/v1/auth/register          # 邮箱注册
POST   /api/v1/auth/login             # 邮箱登录
POST   /api/v1/auth/refresh           # 刷新 Token

# === 用户 ===
GET    /api/v1/users/me               # 获取当前用户
PUT    /api/v1/users/me               # 更新个人资料
PUT    /api/v1/users/me/avatar        # 上传头像
GET    /api/v1/users/:id              # 获取用户公开信息

# === 帖子 ===
GET    /api/v1/posts                  # 帖子列表（?sort=latest|most_helped&page=1&size=20）
POST   /api/v1/posts                  # 创建帖子（支持 multipart 上传图片/语音）
GET    /api/v1/posts/:id              # 帖子详情
DELETE /api/v1/posts/:id              # 删除自己的帖子
POST   /api/v1/posts/:id/stickers     # 发送暖心贴纸
POST   /api/v1/posts/:id/report       # 举报帖子

# === 评论 ===
GET    /api/v1/posts/:id/comments     # 帖子评论列表
POST   /api/v1/posts/:id/comments     # 发表评论
DELETE /api/v1/comments/:id           # 删除自己的评论

# === 情绪 ===
GET    /api/v1/moods                  # 获取情绪记录（?from=2026-06-10&to=2026-06-17）
POST   /api/v1/moods                  # 记录今日情绪
GET    /api/v1/moods/weekly           # 获取本周情绪曲线

# === 每日语录 ===
GET    /api/v1/quotes/today           # 获取今日语录

# === 文件上传 ===
POST   /api/v1/upload/image           # 上传图片（multipart/form-data）
POST   /api/v1/upload/voice           # 上传语音
```

---

## 7. Data Model — 数据模型（PostgreSQL）

### 7.1 users_db（用户服务）

```sql
-- 用户表
CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email       VARCHAR(255) NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL,          -- bcrypt hash
    nickname    VARCHAR(100) NOT NULL,
    avatar_url  TEXT DEFAULT '',
    bio         TEXT DEFAULT '',
    kindness_score INTEGER DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);

-- 刷新令牌表
CREATE TABLE refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token       VARCHAR(512) NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked     BOOLEAN DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
```

### 7.2 posts_db（帖子服务）

```sql
CREATE TABLE posts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    author_id       UUID NOT NULL,
    author_nickname VARCHAR(100) NOT NULL,      -- 冗余字段，事件同步更新
    author_avatar   TEXT DEFAULT '',
    content         TEXT NOT NULL,
    mood_tag        VARCHAR(50) NOT NULL,        -- sad, anxious, lonely, stressed, ...
    images          JSONB DEFAULT '[]',           -- ["url1", "url2"]
    voice_url       TEXT,
    sticker_hug     INTEGER DEFAULT 0,
    sticker_cheer   INTEGER DEFAULT 0,
    sticker_understand INTEGER DEFAULT 0,
    comment_count   INTEGER DEFAULT 0,
    has_ai_reply    BOOLEAN DEFAULT FALSE,
    is_reported     BOOLEAN DEFAULT FALSE,
    is_hidden       BOOLEAN DEFAULT FALSE,       -- 被举报后隐藏
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_posts_author ON posts(author_id);
CREATE INDEX idx_posts_created ON posts(created_at DESC);
CREATE INDEX idx_posts_mood_tag ON posts(mood_tag);
CREATE INDEX idx_posts_comment_count ON posts(comment_count DESC);
-- 全文搜索索引
CREATE INDEX idx_posts_content_search ON posts USING gin(to_tsvector('simple', content));
```

### 7.3 comments_db（评论服务）

```sql
CREATE TABLE comments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id         UUID NOT NULL,
    author_id       UUID NOT NULL,
    author_nickname VARCHAR(100) NOT NULL,
    author_avatar   TEXT DEFAULT '',
    content         TEXT NOT NULL,
    is_ai_generated BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_comments_post ON comments(post_id, created_at ASC);
CREATE INDEX idx_comments_author ON comments(author_id);
```

### 7.4 moods_db（情绪服务）

```sql
CREATE TABLE mood_entries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    mood_level      INTEGER NOT NULL CHECK (mood_level BETWEEN 1 AND 10),
    mood_label      VARCHAR(50) NOT NULL,         -- 开心, 平静, 难过, 焦虑 ...
    note            TEXT,
    created_at      DATE NOT NULL DEFAULT CURRENT_DATE
);

CREATE UNIQUE INDEX idx_moods_user_date ON mood_entries(user_id, created_at);
CREATE INDEX idx_moods_user ON mood_entries(user_id, created_at DESC);
```

### 7.5 quotes_db（语录服务）

```sql
CREATE TABLE daily_quotes (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    text_zh         TEXT NOT NULL,
    text_en         TEXT NOT NULL,
    author          VARCHAR(100),
    background_url  TEXT,
    scheduled_date  DATE NOT NULL UNIQUE,
    pushed          BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_quotes_date ON daily_quotes(scheduled_date);
```

### 7.6 notifications_db（通知服务）

```sql
CREATE TABLE device_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL,
    token       VARCHAR(512) NOT NULL,
    platform    VARCHAR(10) NOT NULL CHECK (platform IN ('ios', 'android', 'web')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_device_tokens_token ON device_tokens(token);
CREATE INDEX idx_device_tokens_user ON device_tokens(user_id);
```

---

## 8. Storage Strategy — 存储策略

### 8.1 存储抽象层

```go
// backend/pkg/storage/storage.go

// Storage 统一存储接口，同时支持 MinIO 和 RustFS
type Storage interface {
    // Upload 上传文件，返回访问 URL
    Upload(ctx context.Context, key string, reader io.Reader, opts UploadOptions) (string, error)
    // Delete 删除文件
    Delete(ctx context.Context, key string) error
    // GetURL 获取文件访问 URL（支持签名）
    GetURL(ctx context.Context, key string, ttl time.Duration) (string, error)
}

type UploadOptions struct {
    ContentType string
    Size        int64
    Public      bool            // 是否公开访问
}
```

### 8.2 使用策略

| 场景 | 存储 | 理由 |
|------|------|------|
| 用户头像 | RustFS / MinIO | 小文件，高频读取 |
| 帖子图片 | MinIO（S3） | 中等文件，需要缩略图处理 |
| 语音文件 | MinIO（S3） | 大文件，需要流式访问 |
| 数据备份 | MinIO（S3） | 批量导出，归档 |
| 开发环境 | MinIO（Docker） | 单机部署足够 |
| 生产环境 | 两者可切换 | 通过 `STORAGE_PROVIDER` 环境变量控制 |

### 8.3 存储配置切换

```yaml
# 使用 MinIO
STORAGE_PROVIDER=minio
MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=showlove
MINIO_USE_SSL=false

# 使用 RustFS
STORAGE_PROVIDER=rustfs
RUSTFS_MASTER=rustfs-master:9333
RUSTFS_FILER=rustfs-filer:8888
```

---

## 9. Testing Strategy — 测试策略

### 9.1 后端测试金字塔

```
         ╱  E2E 测试  ╲           ← 3 条核心 API 流程
        ╱──────────────╲
       ╱   集成测试     ╲          ← 服务间 gRPC 调用 + DB
      ╱─────────────────╲
     ╱     单元测试       ╲        ← 业务逻辑全覆盖
    ╱───────────────────────╲
```

### 9.2 各层测试要求

| 层级 | 框架 | 位置 | 覆盖要求 |
|------|------|------|----------|
| **Go 单元测试** | `testing` + `testify` + `gomock` | 各服务 `internal/` 下的 `*_test.go` | Service 层 ≥ 80% |
| **Go 集成测试** | `testcontainers-go` | `backend/services/*/test/integration/` | Repository + gRPC Server ≥ 60% |
| **Go E2E 测试** | `httptest` + Docker Compose | `backend/test/e2e/` | 核心 API 流程 3 条 |
| **Flutter 单元测试** | `flutter_test` + `mockito` | `frontend/test/unit/` | Domain + Provider ≥ 80% |
| **Flutter Widget 测试** | `flutter_test` | `frontend/test/widget/` | 所有通用组件 + 关键页面 |

### 9.3 Go 测试示例

```go
// ✅ 单元测试: mock 依赖，专注业务逻辑
func TestUserService_CreateUser_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mock.NewMockUserRepository(ctrl)
    svc := service.NewUserService(mockRepo, nil)

    mockRepo.EXPECT().
        Create(gomock.Any(), gomock.Any()).
        Return(nil)

    user, err := svc.Register(context.Background(), &service.RegisterParams{
        Email:    "test@example.com",
        Password: "securePassword123",
        Nickname: "小温暖",
    })

    require.NoError(t, err)
    require.Equal(t, "test@example.com", user.Email)
    require.NotEmpty(t, user.ID)
}

// ✅ 集成测试: testcontainers 启动真实 PostgreSQL
func TestUserRepository_Create(t *testing.T) {
    ctx := context.Background()
    container, db := setupTestDB(t)   // testcontainers-go
    defer container.Terminate(ctx)

    repo := repository.NewUserRepository(db)
    user := &model.User{
        Email:    "test@example.com",
        Password: "hashed_password",
        Nickname: "测试用户",
    }

    err := repo.Create(ctx, user)
    assert.NoError(t, err)
    assert.NotEmpty(t, user.ID)
}
```

### 9.4 核心集成测试场景

```
后端 E2E:
1. POST /auth/register → POST /auth/login → GET /users/me（完整认证流程）
2. POST /posts → GET /posts → GET /posts/:id → POST /posts/:id/comments（发帖到评论流程）
3. POST /moods → GET /moods/weekly（情绪记录到曲线展示）

前端集成测试:
1. 注册 → 设置头像 → 发布第一篇帖子 → 看到帖子出现在首页
2. 浏览首页 → 点开帖子 → 发送评论 → 发送暖心贴纸
3. 打开情绪温度计 → 记录今日心情 → 查看一周曲线
```

---

## 10. Deployment — 部署方案（Docker Compose MVP）

### 10.1 服务编排

```yaml
# backend/docker-compose.yml（核心基础设施 + 全部微服务）

services:
  # ========== 基础设施 ==========
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_MULTIPLE_DATABASES: users_db,posts_db,comments_db,moods_db,quotes_db,notifications_db
    volumes:
      - pgdata:/var/lib/postgresql/data
      - ./scripts/init-multi-db.sh:/docker-entrypoint-initdb.d/init.sh

  redis:
    image: redis:7.2-alpine
    command: redis-server --appendonly yes
    volumes:
      - redisdata:/data

  minio:
    image: minio/minio:latest
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    volumes:
      - miniodata:/data
    ports:
      - "9000:9000"   # S3 API
      - "9001:9001"   # Console

  nats:
    image: nats:2.10-alpine
    command: -js

  # ========== API 网关 ==========
  gateway:
    build: ./gateway
    ports:
      - "8080:8080"
    depends_on:
      - user-service
      - post-service
      - comment-service
      - mood-service
      - quote-service
      - ai-service
    environment:
      GIN_MODE: release
      JWT_SECRET: ${JWT_SECRET}
      USER_SVC_ADDR: user-service:50051
      POST_SVC_ADDR: post-service:50052
      COMMENT_SVC_ADDR: comment-service:50053
      MOOD_SVC_ADDR: mood-service:50054
      QUOTE_SVC_ADDR: quote-service:50055
      AI_SVC_ADDR: ai-service:50056

  # ========== 微服务 ==========
  user-service:
    build: ./services/user-service
    depends_on: [postgres, redis]
    environment:
      DB_DSN: postgres://user:pass@postgres:5432/users_db?sslmode=disable
      REDIS_ADDR: redis:6379
      JWT_SECRET: ${JWT_SECRET}

  post-service:
    build: ./services/post-service
    depends_on: [postgres, redis, minio, nats]
    environment:
      DB_DSN: postgres://user:pass@postgres:5432/posts_db?sslmode=disable
      REDIS_ADDR: redis:6379
      STORAGE_PROVIDER: minio
      MINIO_ENDPOINT: minio:9000

  comment-service:
    build: ./services/comment-service
    depends_on: [postgres]
    environment:
      DB_DSN: postgres://user:pass@postgres:5432/comments_db?sslmode=disable

  mood-service:
    build: ./services/mood-service
    depends_on: [postgres]
    environment:
      DB_DSN: postgres://user:pass@postgres:5432/moods_db?sslmode=disable

  quote-service:
    build: ./services/quote-service
    depends_on: [postgres, nats]
    environment:
      DB_DSN: postgres://user:pass@postgres:5432/quotes_db?sslmode=disable

  ai-service:
    build: ./services/ai-service
    depends_on: [nats]
    environment:
      AI_PROVIDER: ${AI_PROVIDER:-openai}
      OPENAI_API_KEY: ${OPENAI_API_KEY}
      OPENAI_BASE_URL: ${OPENAI_BASE_URL}
      DEEPSEEK_API_KEY: ${DEEPSEEK_API_KEY}
      QWEN_API_KEY: ${QWEN_API_KEY}

  notification-service:
    build: ./services/notification-service
    depends_on: [postgres, nats]
    environment:
      DB_DSN: postgres://user:pass@postgres:5432/notifications_db?sslmode=disable
      FCM_SERVICE_ACCOUNT: /run/secrets/fcm_service_account

volumes:
  pgdata:
  redisdata:
  miniodata:
```

### 10.2 快速启动

```bash
# 1. 创建环境变量
cp backend/.env.example backend/.env
# 编辑 .env，填入 JWT_SECRET 和 AI_API_KEY

# 2. 一键启动
cd backend && docker compose up -d

# 3. 验证
curl http://localhost:8080/api/v1/health
# → {"status":"ok","services":{"gateway":"healthy","user":"healthy","post":"healthy",...}}
```

---

## 11. Boundaries — 边界规则

### 11.1 ✅ Always — 必须做

**通用：**
- 提交前运行 `make lint && make test`，全绿才能 push
- 提交信息使用约定式提交: `feat:`, `fix:`, `refactor:`, `test:`, `docs:`, `chore:`
- 所有 API 请求必须经过 JWT 认证（白名单路由除外）
- 所有用户输入在客户端 + 服务端双重校验
- API 响应统一格式: `{ "code": 0, "message": "ok", "data": {...} }`

**后端：**
- 每个服务独立数据库，禁止跨库 JOIN
- gRPC 调用必须设置 deadline（context.WithTimeout）
- 数据库操作必须通过 Repository 层
- 生产日志使用 zerolog JSON 格式
- 敏感配置通过环境变量注入，禁止硬编码

**前端：**
- 网络请求必须有 loading / error 状态
- 使用 `.arb` 文件管理所有用户可见字符串
- Widget 文件控制在 200 行以内

### 11.2 ⚠️ Ask First — 先问再改

- 新增微服务（创建新的 `backend/services/<name>/`）
- 修改 protobuf 定义（影响前后端契约）
- 修改数据库 schema（需要迁移脚本 + 审批）
- 新增第三方 Go/Dart 依赖
- 修改 docker-compose 编排结构
- 更换 AI 模型提供商
- 修改 JWT 过期策略
- 添加新的 API 版本（`/api/v2/`）
- 修改 `analysis_options.yaml` 或 `golangci-lint` 规则

### 11.3 🚫 Never — 绝不

- 提交密钥文件（`.env`, `*.pem`, `service-account.json`）
- 硬编码 API Key / Token / Secret / 数据库密码
- 在生产环境输出 debug 级别日志
- 跳过错误处理（"happy path only"）
- 跨服务直接访问其他服务的数据库
- 在 Gin handler 中直接写业务逻辑
- 删除失败测试来"修复" CI
- 在未确认 SPEC 更新的情况下实现 spec 外的功能
- 编辑自动生成的 protobuf 代码
- 在 Widget 里直接调 HTTP API（必须走 Repository 层）

---

## 12. 治愈风格设计系统（不变）

### 12.1 色彩体系

```
主色调（温暖珊瑚橙）:  #FF8C69
辅助色（柔和奶油黄）:  #FFE4B5
辅助色（安抚薄荷绿）:  #98D8C8
辅助色（梦幻薰衣草）:  #C3AED6

背景色（暖米白）:      #FFFAF5
卡片色（纯白）:        #FFFFFF
文字主色（深棕）:      #4A3728
文字辅色（灰棕）:      #8B7355
分割线（浅驼）:        #E8D5C4

错误/警示:            #FF6B6B (柔和不刺眼)
成功确认:            #6BCB77
```

### 12.2 治愈元素清单

- [ ] 卡通 IP 角色（「小暖」），出现在空状态、加载、AI 回复
- [ ] 暖心贴纸系统：「抱抱🤗」「加油💪」「我懂你💛」「会好的✨」「你不是一个人🫂」
- [ ] 每日语录卡片：手写体字体 + 柔和渐变背景
- [ ] 情绪温度计：从蓝色（低落）渐变到橙色（温暖）到黄色（开心）
- [ ] 背景可选柔和动画渐变色 / 静态暖色调

---

## 13. 发布计划（MVP → V1）

### Phase 1: 核心 MVP（6-8 周）
- [ ] 项目骨架搭建：Monorepo 结构 + Docker Compose + protobuf 定义
- [ ] Gateway + JWT 认证中间件
- [ ] user-service：邮箱注册/登录 + Token 管理 + 用户资料
- [ ] post-service：帖子 CRUD + 列表（分页/排序）+ 敏感词过滤
- [ ] comment-service：评论 CRUD + 暖心贴纸
- [ ] ai-service：帖子创建后自动生成 AI 回复
- [ ] 文件上传：MinIO 集成（图片 + 语音）
- [ ] Flutter 前端适配：Dio HTTP 客户端替换 Firebase SDK
- [ ] CI/CD：后端 lint + test

### Phase 2: 治愈增强（+3-4 周）
- [ ] mood-service：情绪温度计 + 周曲线统计
- [ ] quote-service：每日语录管理 + 定时推送
- [ ] notification-service：FCM 推送集成
- [ ] 卡通 IP 角色 + 空状态插画
- [ ] 治愈系动效全面覆盖
- [ ] 中文 + 英文双语

### Phase 3: 社区深化（+3-4 周）
- [ ] 善意积分与勋章系统
- [ ] 第三方登录（Google / Apple）
- [ ] 暗色主题
- [ ] 帖子全文搜索
- [ ] 用户间私信（可选）
- [ ] 心理咨询师认证标识（可选）
- [ ] 性能优化 + 压测 + 生产就绪

---

## 14. Open Questions — 待确认

| # | 问题 | 决策 | 状态 |
|---|------|------|------|
| Q1 | 前端从 Firebase SDK 迁移到 REST API 调用 | 确认，保持 Flutter | ✅ |
| Q2 | 微服务拆分粒度：细粒度 7 个服务 | 确认 | ✅ |
| Q3 | 存储方案：同时支持 MinIO 和 RustFS | 确认 | ✅ |
| Q4 | MVP 部署：Docker Compose | 确认 | ✅ |
| Q5 | 服务间通信：gRPC | 确认 | ✅ |
| Q6 | API 网关：自建 Gin 网关 | 确认 | ✅ |
| Q7 | 认证方案：JWT + Refresh Token | 确认 | ✅ |
| Q8 | AI 模型：多模型可切换 | 确认 | ✅ |
| Q9 | 推送时间窗口：每天几点推送？ | 暂定 08:00（Asia/Shanghai） | 待确认 |
| Q10 | AI 回复签名用「小暖」？ | 是 | 待确认 |
| Q11 | 帖子默认排序：最新 or 最需要帮助？ | 最新 | 待确认 |

---

> 📌 **下一步**: 确认此 SPEC 后，进入 Phase 2 — 制定技术实现计划（Plan），
> 然后拆分为可执行任务（Tasks），最后进入编码实现。
