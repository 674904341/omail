# GitHub OAuth 配置

## 环境变量

```bash
# GitHub OAuth 配置
GITHUB_OAUTH_ID=your_github_oauth_app_id
GITHUB_OAUTH_SECRET=your_github_oauth_secret
GITHUB_OAUTH_REDIRECT=http://localhost:3000/api/auth/login
```

## 获取 GitHub OAuth 凭证

1. 访问 https://github.com/settings/developers
2. 点击 "New OAuth App"
3. 填写应用信息：
   - Application name: Tmail
   - Homepage URL: http://localhost:3000 (或你的生产域名)
   - Authorization callback URL: http://localhost:3000/api/auth/login

4. 获取 Client ID 和 Client Secret，配置到环境变量

## API 认证流程

### 1. 获取 GitHub 授权 URL

```bash
GET /api/auth/url?state=random_state_string
```

响应:
```json
{
  "auth_url": "https://github.com/login/oauth/authorize?..."
}
```

### 2. 用户授权后回调

GitHub 会重定向到:
```
/api/auth/login?code=authorization_code&state=state_string
```

### 3. 获取 API Token

```bash
POST /api/auth/login?code=authorization_code&state=state_string
```

响应:
```json
{
  "user": {
    "id": 1,
    "username": "octocat",
    "avatar": "https://avatars.githubusercontent.com/u/1?v=4",
    "email": "octocat@github.com"
  },
  "api_token": "abcd1234efgh5678ijkl9012mnop3456"
}
```

## 使用 API Token

在所有需要认证的请求中添加 Authorization header:

```bash
Authorization: Bearer abcd1234efgh5678ijkl9012mnop3456
```

## API 端点

### 认证相关

- `GET /api/auth/url?state=xxx` - 获取 GitHub 授权 URL
- `POST /api/auth/login?code=xxx&state=xxx` - 处理 GitHub 回调

### 用户相关

- `GET /api/profile` - 获取当前用户信息（需要 Token）
- `POST /api/mailbox` - 创建随机邮箱（需要 Token）
- `GET /api/mailboxes` - 获取用户的所有邮箱（需要 Token）

### 邮件相关

- `GET /api/emails?email=xxx@mail.4w.ink` - 获取邮箱的邮件列表（需要 Token）
- `GET /api/email/:id` - 获取邮件详情（需要 Token）

## 示例代码

### JavaScript/Fetch

```javascript
// 1. 获取授权 URL
const state = Math.random().toString(36).substring(7);
const response = await fetch(`/api/auth/url?state=${state}`);
const { auth_url } = await response.json();
window.location.href = auth_url;

// 2. 处理回调（后端处理）
// GitHub 重定向回来时，从 code 参数中获取认证码
// 这通常在前端登录页面中处理

// 3. 获取用户信息（使用 Token）
const token = localStorage.getItem('api_token');
const profileResponse = await fetch('/api/profile', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});
const profile = await profileResponse.json();
```

### cURL

```bash
# 创建邮箱
curl -X POST http://localhost:3000/api/mailbox \
  -H "Authorization: Bearer YOUR_API_TOKEN"

# 获取邮箱列表
curl http://localhost:3000/api/mailboxes \
  -H "Authorization: Bearer YOUR_API_TOKEN"

# 获取邮件列表
curl "http://localhost:3000/api/emails?email=random@mail.4w.ink" \
  -H "Authorization: Bearer YOUR_API_TOKEN"

# 获取邮件详情
curl http://localhost:3000/api/email/1 \
  -H "Authorization: Bearer YOUR_API_TOKEN"
```

## 数据库字段

### User 表
- `github_id` - GitHub 用户 ID（唯一）
- `username` - GitHub 用户名
- `avatar_url` - 用户头像 URL
- `email` - 邮箱
- `created_at` - 创建时间
- `updated_at` - 更新时间

### APIToken 表
- `token` - API Token（唯一）
- `name` - Token 名称
- `created_at` - 创建时间
- `last_used_at` - 最后使用时间
- `revoked` - 是否已撤销
- `user_id` - 关联的用户 ID

### Mailbox 表
- `email` - 邮箱地址
- `created_at` - 创建时间
- `user_id` - 关联的用户 ID
