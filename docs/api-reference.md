# Show Love API 参考文档

> Base URL: `http://localhost:8080/api/v1`
> 认证方式: Bearer Token (JWT)
> 响应格式: `{ "code": 0, "message": "ok", "data": {...} }`

## 认证

### POST /auth/register — 邮箱注册

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"Secure123","nickname":"小温暖"}'
```

**响应** `201`:
```json
{
  "code": 0,
  "message": "注册成功",
  "data": {
    "user": {"id": "uuid", "email": "user@example.com", "nickname": "小温暖"},
    "access_token": "eyJ...",
    "refresh_token": "abc123..."
  }
}
```

### POST /auth/login — 邮箱登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"Secure123"}'
```

### POST /auth/refresh — 刷新 Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"abc123..."}'
```

## 用户

### GET /users/me — 获取当前用户

```bash
curl http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN"
```

### PUT /users/me — 更新个人资料

```bash
curl -X PUT http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"nickname":"新昵称","bio":"介绍一下自己"}'
```

## 帖子

### GET /posts — 帖子列表

```bash
curl "http://localhost:8080/api/v1/posts?sort=latest&page=1&size=20" \
  -H "Authorization: Bearer $TOKEN"
```

参数: `sort=latest|most_helped`, `page=1`, `size=20`

### POST /posts — 创建帖子

```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":"今天心情不太好","mood_tag":"sad","images":[]}'
```

### GET /posts/:id — 帖子详情

```bash
curl http://localhost:8080/api/v1/posts/$POST_ID \
  -H "Authorization: Bearer $TOKEN"
```

### DELETE /posts/:id — 删除帖子

```bash
curl -X DELETE http://localhost:8080/api/v1/posts/$POST_ID \
  -H "Authorization: Bearer $TOKEN"
```

### POST /posts/:id/stickers — 发送暖心贴纸

```bash
curl -X POST http://localhost:8080/api/v1/posts/$POST_ID/stickers \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"sticker_type":"hug"}'
```

贴纸类型: `hug`(抱抱), `cheer`(加油), `understand`(我懂你)

### POST /posts/:id/report — 举报帖子

```bash
curl -X POST http://localhost:8080/api/v1/posts/$POST_ID/report \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"reason":"内容不当"}'
```

## 评论

### GET /posts/:id/comments — 评论列表

```bash
curl "http://localhost:8080/api/v1/posts/$POST_ID/comments?page=1&size=20" \
  -H "Authorization: Bearer $TOKEN"
```

### POST /posts/:id/comments — 发表评论

```bash
curl -X POST http://localhost:8080/api/v1/posts/$POST_ID/comments \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"content":"加油！一切都会好起来的！"}'
```

## 情绪

### POST /moods — 记录今日情绪

```bash
curl -X POST http://localhost:8080/api/v1/moods \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"mood_level":7,"mood_label":"平静","note":"还不错"}'
```

`mood_level`: 1-10, 1最低落, 10最开心

### GET /moods — 获取情绪记录

```bash
curl "http://localhost:8080/api/v1/moods?from=2026-06-10&to=2026-06-17" \
  -H "Authorization: Bearer $TOKEN"
```

### GET /moods/weekly — 本周情绪曲线

```bash
curl http://localhost:8080/api/v1/moods/weekly \
  -H "Authorization: Bearer $TOKEN"
```

## 每日语录

### GET /quotes/today — 今日语录

```bash
curl http://localhost:8080/api/v1/quotes/today \
  -H "Authorization: Bearer $TOKEN"
```

## 设备推送

### POST /devices — 注册推送设备

```bash
curl -X POST http://localhost:8080/api/v1/devices \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"token":"fcm-device-token","platform":"android"}'
```

platform: `ios` | `android` | `web`

## 文件上传

### POST /upload/image — 上传图片

```bash
curl -X POST http://localhost:8080/api/v1/upload/image \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@photo.jpg"
```

限制: JPEG/PNG/GIF/WebP, 最大 10MB

## 健康检查

### GET /health — 服务健康状态

```bash
curl http://localhost:8080/api/v1/health
```

## 错误码

| Code | 说明 |
|------|------|
| 0 | 成功 |
| 400 | 参数错误 |
| 401 | 未认证 (Token缺失/过期/无效) |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 409 | 资源已存在 |
| 429 | 请求过于频繁 |
| 500 | 服务器内部错误 |
