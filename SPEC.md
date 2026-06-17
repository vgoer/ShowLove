# Spec: 显出爱心 (Show Love) — 治愈系社区互助 App

> **类型**: 创业 MVP · **阶段**: 规格确认完毕 · **最后更新**: 2026-06-17

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
| S4 | 首页帖子列表加载时间 < 2 秒（4G 网络） | 性能测试 |
| S5 | 核心崩溃率 < 0.5% | Firebase Crashlytics |
| S6 | 治愈风格一致性：所有页面通过 UI 审查清单 | 人工审查 |

---

## 2. Tech Stack — 技术选型

| 层级 | 技术 | 版本 | 选择理由 |
|------|------|------|----------|
| 前端框架 | Flutter | ≥ 3.27 | 单代码库覆盖 iOS/Android/Web |
| 语言 | Dart | ≥ 3.6 | — |
| 状态管理 | Riverpod | ≥ 2.6 | 编译时安全，测试友好，社区主流 |
| 路由 | go_router | ≥ 14 | Flutter 官方推荐 |
| 后端服务 | Firebase | — | 与 Flutter 深度集成，免运维 |
| 认证 | Firebase Auth | — | 邮箱 + Google + Apple 登录 |
| 数据库 | Cloud Firestore | — | 实时同步，适合帖子/评论 |
| 文件存储 | Cloud Storage | — | 图片/语音文件 |
| 服务端逻辑 | Cloud Functions | Node.js 20 | 敏感词过滤 + AI 调用 |
| AI | OpenAI GPT-4o-mini | — | 成本低，回复质量好 |
| 推送 | Firebase Cloud Messaging | — | 每日语录 + 评论通知 |
| 崩溃监控 | Firebase Crashlytics | — | Flutter 原生集成 |
| 分析 | Firebase Analytics | — | 用户行为追踪 |

### 2.1 最低平台版本

| 平台 | 最低版本 | 覆盖率 |
|------|----------|--------|
| iOS | 15.0+ | ~96% |
| Android | API 26 (Android 8.0) | ~95% |

---

## 3. Commands — 开发命令

```bash
# === 环境初始化 ===
flutter clean && flutter pub get          # 全新安装依赖
dart run build_runner build --delete-conflicting-outputs  # 重新生成代码

# === 开发运行 ===
flutter run                               # 以 debug 模式运行（默认设备）
flutter run -d chrome                     # Web 预览
flutter run --flavor dev                  # 开发环境运行

# === 静态检查 ===
flutter analyze                           # Dart 静态分析（CI 强制通过）
dart format --set-exit-if-changed lib/ test/  # 格式化检查

# === 测试 ===
flutter test                              # 全部单元 + Widget 测试
flutter test --coverage                   # 带覆盖率报告
flutter test test/unit/                   # 仅单元测试
flutter test test/widget/                 # 仅 Widget 测试
flutter test integration_test/            # 集成测试（需模拟器）

# === 构建 ===
flutter build apk --flavor prod           # Android 生产包
flutter build ios --flavor prod           # iOS 生产包
flutter build web --flavor prod           # Web 生产包

# === Firebase ===
firebase deploy --only functions          # 部署云函数
firebase deploy --only firestore:rules    # 部署数据库规则
firebase deploy --only storage:rules      # 部署存储规则
```

---

## 4. Project Structure — 项目结构

```
show_love/
├── android/                          # Android 原生代码
├── ios/                              # iOS 原生代码
├── web/                              # Web 平台代码
├── assets/
│   ├── images/                       # 静态图片（插画、贴纸、图标）
│   ├── animations/                   # Lottie 动画文件
│   └── fonts/                        # 自定义字体（治愈手写体）
├── lib/
│   ├── main.dart                     # 入口 → ProviderScope + App
│   ├── app.dart                      # MaterialApp.router 配置
│   │
│   ├── core/                         # 全局基础设施（不依赖任何 feature）
│   │   ├── constants/
│   │   │   ├── app_colors.dart       # 治愈系色彩体系
│   │   │   ├── app_theme.dart        # ThemeData（亮色/暗色）
│   │   │   ├── app_text_styles.dart  # 文字样式系统
│   │   │   └── api_constants.dart    # API 端点 / Cloud Function 名
│   │   ├── extensions/
│   │   │   ├── context_ext.dart      # BuildContext 扩展
│   │   │   └── string_ext.dart       # String 扩展
│   │   ├── utils/
│   │   │   ├── validators.dart       # 表单校验工具
│   │   │   ├── date_formatter.dart   # 日期格式化（"3 分钟前"）
│   │   │   └── logger.dart           # 统一日志工具
│   │   ├── widgets/                  # 全局共享组件
│   │   │   ├── love_button.dart      # 主按钮（圆角、阴影、微动效）
│   │   │   ├── love_text_field.dart  # 输入框（治愈风格）
│   │   │   ├── avatar_widget.dart    # 头像组件
│   │   │   ├── loading_indicator.dart# 加载动画
│   │   │   ├── error_widget.dart     # 错误状态组件
│   │   │   ├── empty_state_widget.dart # 空状态组件
│   │   │   └── animated_gradient_bg.dart # 动画渐变背景
│   │   ├── router/
│   │   │   └── app_router.dart       # GoRouter 路由配置
│   │   ├── l10n/                     # 国际化
│   │   │   ├── app_zh.arb            # 中文
│   │   │   └── app_en.arb            # 英文
│   │   └── errors/
│   │       ├── app_exception.dart    # 自定义异常
│   │       └── error_handler.dart    # 全局错误处理
│   │
│   ├── data/                         # 数据层（Repository 实现 + DataSource）
│   │   ├── repositories/             # Repository 实现
│   │   │   ├── auth_repository.dart
│   │   │   ├── post_repository.dart
│   │   │   ├── mood_repository.dart
│   │   │   └── quote_repository.dart
│   │   └── datasources/
│   │       ├── firebase_auth_datasource.dart
│   │       ├── firestore_datasource.dart
│   │       ├── storage_datasource.dart
│   │       ├── cloud_functions_datasource.dart
│   │       └── local_storage_datasource.dart  # SharedPreferences 封装
│   │
│   ├── domain/                       # 领域层（纯 Dart，零依赖）
│   │   ├── models/
│   │   │   ├── user.dart             # 用户模型
│   │   │   ├── post.dart             # 帖子模型
│   │   │   ├── comment.dart          # 评论模型
│   │   │   ├── mood_entry.dart       # 情绪记录模型
│   │   │   └── quote.dart            # 每日语录模型
│   │   └── repositories/             # Repository 接口（抽象）
│   │       ├── i_auth_repository.dart
│   │       ├── i_post_repository.dart
│   │       ├── i_mood_repository.dart
│   │       └── i_quote_repository.dart
│   │
│   ├── features/                     # 功能模块（每个 feature 自包含）
│   │   ├── auth/
│   │   │   ├── screens/
│   │   │   │   ├── login_screen.dart
│   │   │   │   ├── register_screen.dart
│   │   │   │   └── profile_setup_screen.dart
│   │   │   └── providers/
│   │   │       └── auth_provider.dart
│   │   │
│   │   ├── feed/
│   │   │   ├── screens/
│   │   │   │   └── feed_screen.dart          # 首页帖子流
│   │   │   ├── widgets/
│   │   │   │   ├── post_card.dart            # 帖子卡片
│   │   │   │   ├── post_card_skeleton.dart   # 骨架屏
│   │   │   │   └── mood_filter_chips.dart    # 心情筛选 chips
│   │   │   └── providers/
│   │   │       └── feed_provider.dart
│   │   │
│   │   ├── post_detail/
│   │   │   ├── screens/
│   │   │   │   └── post_detail_screen.dart
│   │   │   ├── widgets/
│   │   │   │   ├── comment_tile.dart
│   │   │   │   ├── sticker_bar.dart          # 暖心贴纸栏
│   │   │   │   └── ai_reply_banner.dart      # AI 回复横幅
│   │   │   └── providers/
│   │   │       └── post_detail_provider.dart
│   │   │
│   │   ├── create_post/
│   │   │   ├── screens/
│   │   │   │   └── create_post_screen.dart
│   │   │   ├── widgets/
│   │   │   │   ├── mood_picker.dart          # 心情标签选择器
│   │   │   │   ├── image_picker_grid.dart    # 图片选择网格
│   │   │   │   └── voice_recorder.dart       # 语音录制组件
│   │   │   └── providers/
│   │   │       └── create_post_provider.dart
│   │   │
│   │   ├── mood_tracker/
│   │   │   ├── screens/
│   │   │   │   └── mood_tracker_screen.dart
│   │   │   ├── widgets/
│   │   │   │   ├── emotion_thermometer.dart  # 情绪温度计
│   │   │   │   └── mood_chart.dart           # 一周情绪曲线
│   │   │   └── providers/
│   │   │       └── mood_provider.dart
│   │   │
│   │   ├── daily_quote/
│   │   │   ├── widgets/
│   │   │   │   └── daily_quote_card.dart     # 每日语录卡片
│   │   │   └── providers/
│   │   │       └── quote_provider.dart
│   │   │
│   │   ├── profile/
│   │   │   ├── screens/
│   │   │   │   ├── profile_screen.dart
│   │   │   │   └── edit_profile_screen.dart
│   │   │   └── providers/
│   │   │       └── profile_provider.dart
│   │   │
│   │   └── settings/
│   │       ├── screens/
│   │       │   └── settings_screen.dart
│   │       └── providers/
│   │           └── settings_provider.dart
│   │
│   └── gen/                          # 自动生成（国际化、资源）
│       └── ...
│
├── test/
│   ├── unit/                         # 单元测试
│   │   ├── domain/
│   │   │   └── models/               # 模型序列化/反序列化测试
│   │   ├── core/
│   │   │   └── utils/                # 工具函数测试
│   │   └── features/                 # Provider / Notifier 逻辑测试
│   │       ├── auth/
│   │       ├── feed/
│   │       ├── create_post/
│   │       └── mood_tracker/
│   ├── widget/                       # Widget 测试
│   │   ├── core/
│   │   │   └── widgets/              # 通用组件测试
│   │   └── features/
│   │       ├── feed/
│   │       └── post_detail/
│   └── integration_test/             # 集成测试（E2E）
│       ├── auth_flow_test.dart
│       ├── post_flow_test.dart
│       └── mood_flow_test.dart
│
├── functions/                        # Firebase Cloud Functions
│   ├── src/
│   │   ├── index.ts                  # 入口
│   │   ├── ai/
│   │   │   └── generateReply.ts      # AI 暖心回复
│   │   ├── moderation/
│   │   │   └── filterContent.ts      # 敏感词过滤
│   │   └── notifications/
│   │       └── dailyQuote.ts         # 每日语录定时推送
│   ├── package.json
│   └── tsconfig.json
│
├── firebase.json                     # Firebase 配置
├── firestore.rules                   # Firestore 安全规则
├── storage.rules                     # Storage 安全规则
├── pubspec.yaml                      # Flutter 依赖
├── analysis_options.yaml             # Dart 静态分析规则
└── SPEC.md                           # ← 本文件
```

---

## 5. Code Style — 代码风格

### 5.1 命名规范

```dart
// ✅ 文件命名: snake_case
// post_card.dart, auth_provider.dart, mood_tracker_screen.dart

// ✅ 类命名: PascalCase (Widget, Provider, Model)
class PostCard extends ConsumerWidget { ... }
class AuthNotifier extends StateNotifier<AuthState> { ... }

// ✅ 变量/函数命名: camelCase
final userName = 'Alice';
String formatRelativeTime(DateTime date) { ... }

// ✅ 常量: camelCase (不用 SCREAMING_SNAKE)
const primaryWarm = Color(0xFFFF8C69);

// ✅ 私有成员: 前缀 _
final _firestore = FirebaseFirestore.instance;
void _handleRefresh() { ... }
```

### 5.2 组件风格

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

// ❌ 避免: 巨大的 build 方法，一个 Widget 干所有事
class PostCard extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    // 300 行代码揉在一起 ...
  }
}
```

### 5.3 状态管理模式

```dart
// ✅ 使用 Riverpod AsyncNotifier 处理异步
@riverpod
class FeedNotifier extends _$FeedNotifier {
  @override
  Future<List<Post>> build() => _fetchPosts();

  Future<void> refresh() async {
    state = const AsyncValue.loading();
    state = await AsyncValue.guard(() => _fetchPosts());
  }
}

// ✅ 统一错误处理
state = await AsyncValue.guard(() async {
  return await _repository.fetchPosts();
});
// AsyncValue 自动捕获异常，UI 层用 .when() 优雅处理
```

### 5.4 关键约定

| 规则 | 说明 |
|------|------|
| `const` 优先 | 能用 `const` 的地方一律用（性能优化） |
| 一行 ≤ 80 字符 | `dart format --line-length 80` |
| Widget 拆分 | 单个 `build` 方法不超过 50 行 |
| 禁止 `dynamic` | 除非调用 Firebase 原生 API |
| 导入别名 | `import 'package:flutter/material.dart' as m;` 避免命名冲突 |
| 模型不可变 | 使用 `freezed` / `@immutable` |

---

## 6. Testing Strategy — 测试策略

### 6.1 测试金字塔

```
         ╱ 集成测试 ╲          ← 3 条核心流程
        ╱─────────────╲
       ╱  Widget 测试  ╲        ← 关键 UI 组件
      ╱─────────────────╲
     ╱    单元测试         ╲     ← 业务逻辑全覆盖
    ╱─────────────────────────╲
```

### 6.2 各层测试要求

| 层级 | 框架 | 位置 | 覆盖要求 |
|------|------|------|----------|
| **单元测试** | `flutter_test` + `mockito` | `test/unit/` | Domain 模型 & Provider 逻辑 ≥ 80% |
| **Widget 测试** | `flutter_test` | `test/widget/` | 所有通用组件 + 关键页面 |
| **集成测试** | `integration_test` | `integration_test/` | 3 条核心用户旅程（见下） |

### 6.3 核心集成测试场景

```
1. 注册 → 设置头像 → 发布第一篇帖子 → 看到帖子出现在首页
2. 浏览首页 → 点开帖子 → 发送评论 → 发送暖心贴纸
3. 打开情绪温度计 → 记录今日心情 → 查看一周曲线
```

### 6.4 测试原则

- **Provider 逻辑**：100% 单元测试（mock repository）
- **Repository**：单元测试 mock datasource，验证缓存/重试逻辑
- **Widget**：验证三种状态（loading / data / error）+ 交互回调
- **Cloud Functions**：Jest 单元测试 mock Firebase Admin SDK
- **不追求 100% 覆盖率**，但核心路径必须全覆盖

---

## 7. Boundaries — 边界规则

### 7.1 ✅ Always — 必须做

- 提交前运行 `flutter analyze` + `flutter test`，全绿才能 push
- 遵循 `const` 优先原则
- 所有用户输入都做客户端 + 服务端双重校验
- 网络请求必须有 loading / error 状态
- 异常必须通过 `AsyncValue.guard` 或自定义 ErrorHandler 捕获
- 使用 `.arb` 文件管理所有用户可见字符串（国际化就绪）
- Widget 文件控制在 200 行以内
- 提交信息使用约定式提交: `feat:`, `fix:`, `refactor:`, `test:`, `docs:`

### 7.2 ⚠️ Ask First — 先问再改

- 新增第三方依赖（`flutter pub add xxx`）
- 修改 Firestore 数据结构 / 索引 / 安全规则
- 修改 Cloud Functions 部署配置
- 更换或升级核心依赖（Riverpod, go_router, Firebase SDK）
- 添加新的 feature 模块（创建新的 `lib/features/<name>/`）
- 修改 `analysis_options.yaml` 规则
- 任何涉及付费 API 配额变更的操作

### 7.3 🚫 Never — 绝不

- 提交密钥文件（`google-services.json`, `GoogleService-Info.plist`, `.env`）
- 硬编码 API Key / Token / Secret
- 跳过错误处理（"happy path only"）
- 编辑 `lib/gen/` 自动生成的文件
- 在 Widget 里直接调 Firebase API（必须走 Repository 层）
- 删除失败测试来"修复" CI
- 在未确认 SPEC 更新的情况下实现 spec 外的功能
- 在生产环境使用 debug 日志

---

## 8. 数据模型（核心）

### 8.1 Firestore 集合设计

```
users/{uid}
  ├─ nickname: string
  ├─ avatarUrl: string
  ├─ bio: string
  ├─ createdAt: timestamp
  └─ kindnessScore: number        // 善意积分（发帖/评论/被感谢）

posts/{postId}
  ├─ authorId: string (→ users/{uid})
  ├─ content: string
  ├─ moodTag: string               // e.g. "sad", "anxious", "lonely", "stressed"
  ├─ images: string[]              // Cloud Storage URLs
  ├─ voiceUrl: string?             // 语音文件 URL
  ├─ stickerCounts: map<string,int> // {hug: 5, cheer: 3, understand: 2}
  ├─ commentCount: number
  ├─ hasAiReply: bool
  ├─ isReported: bool
  ├─ createdAt: timestamp
  └─ updatedAt: timestamp

comments/{commentId}
  ├─ postId: string (→ posts/{postId})
  ├─ authorId: string (→ users/{uid})
  ├─ content: string
  ├─ isAiGenerated: bool
  └─ createdAt: timestamp

mood_entries/{entryId}
  ├─ userId: string (→ users/{uid})
  ├─ moodLevel: number             // 1-10，温度计刻度
  ├─ moodLabel: string             // "开心", "平静", "难过" ...
  ├─ note: string?
  └─ createdAt: timestamp

daily_quotes/{quoteId}
  ├─ textZh: string                // 中文语录
  ├─ textEn: string                // 英文语录
  ├─ author: string
  ├─ backgroundImage: string?      // 配图 URL
  └─ scheduledDate: timestamp      // 预定推送日期
```

---

## 9. 治愈风格设计系统

### 9.1 色彩体系

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

### 9.2 动效原则

- 页面转场：`SlideTransition` + `FadeTransition`，200ms 缓入缓出
- 暖心贴纸点击：缩放弹跳动画（`Curves.elasticOut`）
- 卡片出现：从下往上滑入 + 轻微淡入（`staggered animation`）
- 情绪温度计：缓慢呼吸式微动（`TweenAnimationBuilder`）
- AI 回复出现：打字机效果逐字显示

### 9.3 治愈元素清单

- [ ] 卡通 IP 角色（可命名，如「小暖」），出现在空状态、加载、AI 回复
- [ ] 暖心贴纸系统：「抱抱🤗」「加油💪」「我懂你💛」「会好的✨」「你不是一个人🫂」
- [ ] 每日语录卡片：手写体字体 + 柔和渐变背景
- [ ] 情绪温度计：从蓝色（低落）渐变到橙色（温暖）到黄色（开心）
- [ ] 背景可选柔和动画渐变色 / 静态暖色调

---

## 10. API & Cloud Functions 设计

### 10.1 Cloud Functions

```typescript
// 1. AI 暖心回复（帖子创建时触发）
export const onPostCreated = firestore
  .document('posts/{postId}')
  .onCreate(async (snap, context) => {
    // 调用 OpenAI 生成鼓励回复
    // 以匿名 "小暖" 身份写入 comments 子集合
  });

// 2. 敏感词过滤（帖子/评论创建时触发）
export const filterContent = firestore
  .document('posts/{postId}')
  .onCreate(async (snap, context) => {
    // 检测敏感词 → 标记 isReported / 屏蔽
  });

// 3. 每日语录推送（定时触发, pubsub schedule）
export const sendDailyQuote = pubsub
  .schedule('every day 08:00')
  .timeZone('Asia/Shanghai')
  .onRun(async (context) => {
    // 随机选取一条语录 → FCM 推送给所有用户
  });
```

### 10.2 客户端 API 抽象

```dart
// lib/domain/repositories/i_post_repository.dart
abstract interface class IPostRepository {
  Future<List<Post>> getFeed({required String sortBy, int limit = 20});
  Future<Post> getPostById(String postId);
  Future<Post> createPost(CreatePostParams params);
  Future<void> addComment(String postId, String content);
  Future<void> sendSticker(String postId, String stickerType);
  Future<void> reportPost(String postId);
}
```

---

## 11. 发布计划（MVP → V1）

### Phase 1: 核心 MVP（4-6 周）
- [ ] 邮箱注册/登录
- [ ] 发布帖子（文字 + 图片 + 心情标签）
- [ ] 首页帖子流（最新/最需要帮助 排序）
- [ ] 评论 + 暖心贴纸
- [ ] AI 暖心回复
- [ ] 基础敏感词过滤 + 举报
- [ ] 个人主页（简单版）

### Phase 2: 治愈增强（+2-3 周）
- [ ] 情绪温度计 + 周曲线
- [ ] 每日语录推送
- [ ] 卡通 IP 角色 + 空状态插画
- [ ] 语音发布
- [ ] 治愈系动效全面覆盖
- [ ] 中文 + 英文双语

### Phase 3: 社区深化（+3-4 周）
- [ ] 善意积分与勋章系统
- [ ] Google / Apple 第三方登录
- [ ] 暗色主题
- [ ] 帖子搜索
- [ ] 用户间私信（可选）
- [ ] 心理咨询师认证标识（可选）

---

## 12. Open Questions — 待确认

| # | 问题 | 决策人 | 状态 |
|---|------|--------|------|
| Q3 | 推送时间窗口：每天几点推送？默认 08:00？ | 创始人 | ✅ 每天 08:00（Asia/Shanghai） |
| Q4 | AI 回复签名用「小暖」还是其他名字？ | 创始人 | ✅ 小暖 |
| Q5 | 帖子排序：默认"最新"还是"最需要帮助"？ | 创始人 | ✅ 默认最新 |

---

> 📌 **下一步**: 确认此 SPEC 后，进入 Phase 2 — 制定技术实现计划（Plan），
> 然后拆分为可执行任务（Tasks），最后进入编码实现。
