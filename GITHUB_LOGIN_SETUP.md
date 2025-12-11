# GitHub 登录配置指南

## 1. 创建 GitHub OAuth App

1. 访问 https://github.com/settings/developers
2. 点击 "New OAuth App" 或 "OAuth Apps" → "New OAuth App"
3. 填写以下信息：

### 应用信息
- **Application name**: Tmail（或你想要的名称）
- **Homepage URL**: 
  - 本地开发: `http://localhost:3000`
  - 生产环境: `https://mail.4w.ink` (替换为你的域名)

- **Authorization callback URL**: 
  - 本地开发: `http://localhost:3000/en/callback`（英文）或 `http://localhost:3000/zh/callback`（中文）
  - 生产环境: `https://mail.4w.ink/en/callback`

4. 点击 "Register application"

## 2. 获取凭证

注册后，你会看到：
- **Client ID**: 保存此值
- **Client Secret**: 点击 "Generate a new client secret" 并保存

## 3. 配置环境变量

### 本地开发 (.env)
```env
GITHUB_OAUTH_ID=你的_Client_ID
GITHUB_OAUTH_SECRET=你的_Client_Secret
GITHUB_OAUTH_REDIRECT=http://localhost:3000
```

### Docker Compose (.env)
```env
GITHUB_OAUTH_ID=你的_Client_ID
GITHUB_OAUTH_SECRET=你的_Client_Secret
GITHUB_OAUTH_REDIRECT=http://localhost:8080
```

### 生产环境 (.env)
```env
GITHUB_OAUTH_ID=你的_Client_ID
GITHUB_OAUTH_SECRET=你的_Client_Secret
GITHUB_OAUTH_REDIRECT=https://mail.4w.ink
```

## 4. 启动应用

```bash
# 本地开发
npm run dev

# Docker 开发
docker-compose up

# 生产环境
docker run -e GITHUB_OAUTH_ID=xxx -e GITHUB_OAUTH_SECRET=xxx -e GITHUB_OAUTH_REDIRECT=xxx tmail
```

## 5. 测试登录流程

1. 访问 `http://localhost:3000/`（或你的域名）
2. 点击 "Sign in with GitHub" 按钮
3. 授权 Tmail 应用访问你的 GitHub 账户
4. 应该会跳转回 callback 页面并显示邮箱界面

## 故障排除

### 问题：回调 URL 不匹配
**解决**: 确保 GitHub App 中的 "Authorization callback URL" 与环境变量中的 `GITHUB_OAUTH_REDIRECT` 一致

### 问题：授权后白屏
**解决**: 
- 检查浏览器控制台是否有错误
- 确保后端 `/api/auth/login` 接口可用
- 检查数据库连接是否正常

### 问题：Token 存储失败
**解决**: 检查浏览器是否允许 localStorage 访问

## 关键端点

- `GET /api/auth/url?state=xxx` - 获取 GitHub 授权 URL
- `GET /api/auth/login?code=xxx&state=xxx` - GitHub 回调处理
- `GET /api/profile` - 获取当前用户信息（需要 Bearer token）
- `POST /api/mailbox` - 创建新邮箱（需要 Bearer token）
- `GET /api/mailboxes` - 获取用户所有邮箱（需要 Bearer token）

## 安全提示

1. ✅ 生产环境务必使用 HTTPS
2. ✅ Client Secret 不要硬编码在代码中，使用环境变量
3. ✅ 定期轮换 GitHub App 凭证
4. ✅ 使用强密码保护你的 GitHub 账户
