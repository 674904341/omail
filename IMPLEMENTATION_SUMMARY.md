# 登录认证和 API 系统实现总结

## 已实现的功能

### 1. GitHub OAuth 登录系统
- ✅ GitHub OAuth 2.0 认证集成
- ✅ 用户自动创建/更新
- ✅ API Token 自动生成

### 2. API Token 认证系统
- ✅ Bearer Token 认证中间件
- ✅ Token 验证和管理
- ✅ Token 撤销功能
- ✅ 最后使用时间追踪

### 3. 邮箱管理 API
- ✅ 创建随机邮箱（`POST /api/mailbox`）
- ✅ 获取用户邮箱列表（`GET /api/mailboxes`）
- ✅ 获取邮件列表（`GET /api/emails?email=xxx`）
- ✅ 获取邮件详情（`GET /api/email/:id`）

### 4. 用户资料 API
- ✅ 获取用户信息（`GET /api/profile`）

### 5. 前端登录 UI
- ✅ GitHub 登录页面（Login.tsx）
- ✅ OAuth 回调处理（LoginCallback.tsx）

## 数据库架构

新增三个表：

### User 表
```
- id (主键)
- github_id (GitHub ID，唯一)
- username (用户名)
- avatar_url (头像)
- email (邮箱)
- created_at (创建时间)
- updated_at (更新时间)
```

### APIToken 表
```
- id (主键)
- token (Token 值，唯一)
- name (Token 名称)
- user_id (用户外键)
- created_at (创建时间)
- last_used_at (最后使用时间)
- revoked (是否撤销)
```

### Mailbox 表
```
- id (主键)
- email (邮箱地址)
- user_id (用户外键)
- created_at (创建时间)
```

### Envelope 表（已有，新增关联）
```
- 新增 mailbox_id 关系
```

## 环境变量配置

```bash
# GitHub OAuth
GITHUB_OAUTH_ID=your_github_app_id
GITHUB_OAUTH_SECRET=your_github_app_secret
GITHUB_OAUTH_REDIRECT=https://yourdomain.com/api/auth/login
```

## 部署步骤

### 1. 生成数据库 Schema

```bash
cd /app
go generate ./ent
```

### 2. 创建/迁移数据库

```bash
# Ent 会自动创建表，或者运行迁移
# 如果使用 PostgreSQL，确保数据库存在
createdb tmail
```

### 3. 设置 GitHub OAuth App

访问 https://github.com/settings/developers：
1. New OAuth App
2. Application name: Tmail
3. Homepage URL: https://yourdomain.com
4. Authorization callback URL: https://yourdomain.com/api/auth/login
5. 复制 Client ID 和 Client Secret

### 4. 配置环境变量

```bash
export GITHUB_OAUTH_ID=xxx
export GITHUB_OAUTH_SECRET=xxx
export GITHUB_OAUTH_REDIRECT=https://yourdomain.com/api/auth/login
```

### 5. 启动应用

```bash
./tmail
```

## API 使用示例

### 1. 获取 GitHub 授权 URL

```bash
curl "http://localhost:3000/api/auth/url?state=random123"

# 响应
{
  "auth_url": "https://github.com/login/oauth/authorize?..."
}
```

### 2. GitHub 回调（用户点击授权后）

```bash
# GitHub 重定向到
GET /api/auth/login?code=abc123&state=random123

# 响应
{
  "user": {
    "id": 1,
    "username": "octocat",
    "avatar": "https://avatars.githubusercontent.com/u/1?v=4",
    "email": "octocat@github.com"
  },
  "api_token": "abcd1234..."
}
```

### 3. 创建邮箱

```bash
curl -X POST http://localhost:3000/api/mailbox \
  -H "Authorization: Bearer abcd1234..." \
  -H "Content-Type: application/json"

# 响应
{
  "email": "abc12345@mail.4w.ink",
  "created_at": "2025-12-11T00:00:00Z"
}
```

### 4. 获取邮箱列表

```bash
curl http://localhost:3000/api/mailboxes \
  -H "Authorization: Bearer abcd1234..."

# 响应
{
  "mailboxes": [
    {
      "email": "abc12345@mail.4w.ink",
      "created_at": "2025-12-11T00:00:00Z"
    }
  ]
}
```

### 5. 获取邮件列表

```bash
curl "http://localhost:3000/api/emails?email=abc12345@mail.4w.ink" \
  -H "Authorization: Bearer abcd1234..."

# 响应
{
  "emails": [
    {
      "id": 1,
      "from": "sender@example.com",
      "subject": "Hello",
      "created_at": "2025-12-11T00:00:00Z"
    }
  ]
}
```

### 6. 获取邮件详情

```bash
curl http://localhost:3000/api/email/1 \
  -H "Authorization: Bearer abcd1234..."

# 响应
{
  "id": 1,
  "from": "sender@example.com",
  "to": "abc12345@mail.4w.ink",
  "subject": "Hello",
  "content": "Hello World",
  "created_at": "2025-12-11T00:00:00Z",
  "attachments": [
    {
      "id": 1,
      "filename": "test.txt",
      "size": 1024
    }
  ]
}
```

## 文件变更清单

### 新增文件
- `ent/schema/user.go` - User 数据模型
- `ent/schema/api_token.go` - APIToken 数据模型
- `ent/schema/mailbox.go` - Mailbox 数据模型
- `internal/api/auth.go` - GitHub OAuth 认证逻辑
- `internal/api/handlers.go` - 认证和邮箱管理 API 处理器
- `internal/api/token_auth.go` - Token 认证中间件
- `web/src/components/Login.tsx` - 登录 UI
- `web/src/components/LoginCallback.tsx` - OAuth 回调处理
- `API_AUTH.md` - API 文档

### 修改文件
- `internal/route/route.go` - 添加新的认证路由和中间件
- `internal/app.go` - 集成 Token 认证中间件
- `internal/api/context.go` - 添加 User() 和 Client() 方法

## 下一步优化建议

1. **前端页面**
   - 创建仪表板页面显示用户邮箱
   - 添加邮件查看界面
   - 实现邮箱删除功能

2. **API 增强**
   - 删除邮箱（`DELETE /api/mailbox/:email`）
   - 列出所有 Token（`GET /api/tokens`）
   - 撤销 Token（`DELETE /api/token/:id`）
   - 刷新 Token（`POST /api/token/refresh`）

3. **安全增强**
   - 实现 CSRF 保护
   - 添加速率限制
   - 实现刷新 Token 机制
   - Token 过期时间设置

4. **数据库**
   - 确保所有新表都正确迁移
   - 添加索引优化查询

## 常见问题

### Q: Token 如何存储？
A: Token 以明文形式存储在数据库中，生产环境建议加密存储。

### Q: 如何撤销 Token？
A: 通过将 APIToken 的 `revoked` 字段设置为 true。

### Q: Token 过期吗？
A: 当前实现没有过期机制，生产环境建议添加。

### Q: 如何备份 API Token？
A: Token 只在创建时返回一次，用户应妥善保管。建议在用户首次登录时提示保存 Token。
