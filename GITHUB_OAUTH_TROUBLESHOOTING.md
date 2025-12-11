# GitHub OAuth 环境变量配置故障排除

## ❌ 问题现象

点击 GitHub 登录时，跳转的 URL 缺少认证信息：
```
https://github.com/login/oauth/authorize?client_id=&redirect_uri=&scope=user&state=xxx
```

client_id 和 redirect_uri 都为空，说明环境变量没有被正确读取。

---

## ✅ 解决方案

### 问题 1：Docker Compose 环境变量格式错误

**❌ 错误的格式：**
```yaml
environment:
  - "GITHUB_OAUTH_ID: Ov23liXwSJQhdRTprTxK"  # ❌ 使用冒号
  - "GITHUB_OAUTH_SECRET: 727a8dda..."        # ❌ 使用冒号
```

**✅ 正确的格式：**
```yaml
environment:
  - GITHUB_OAUTH_ID=Ov23liXwSJQhdRTprTxK     # ✅ 使用等号
  - GITHUB_OAUTH_SECRET=727a8dda...          # ✅ 使用等号
  - GITHUB_OAUTH_REDIRECT=https://mail.4w.ink/en/callback
```

或者使用更清晰的格式：
```yaml
environment:
  GITHUB_OAUTH_ID: Ov23liXwSJQhdRTprTxK
  GITHUB_OAUTH_SECRET: 727a8dda...
  GITHUB_OAUTH_REDIRECT: https://mail.4w.ink/en/callback
```

### 问题 2：回调 URL 错误

**❌ 你的配置：**
```yaml
GITHUB_OAUTH_REDIRECT: https://mail.4w.ink/api/auth/login
```

**为什么错误？**
- `/api/auth/login` 是 API 端点，不是页面
- 用户会被重定向到 API 响应，而不是前端回调页面
- 前端无法处理回调

**✅ 正确的配置：**
```yaml
# 英文页面
GITHUB_OAUTH_REDIRECT: https://mail.4w.ink/en/callback

# 或中文页面（两个都可以）
GITHUB_OAUTH_REDIRECT: https://mail.4w.ink/zh/callback
```

---

## 🔍 诊断步骤

### 步骤 1：验证环境变量是否在容器内

```bash
# 进入正在运行的容器
docker exec -it omail sh

# 检查环境变量是否被设置
echo $GITHUB_OAUTH_ID
echo $GITHUB_OAUTH_SECRET
echo $GITHUB_OAUTH_REDIRECT

# 如果为空，说明环境变量没有被正确传递
```

### 步骤 2：检查后端日志

```bash
# 查看容器日志
docker logs omail | tail -50

# 或实时查看
docker logs -f omail
```

查找类似的日志信息：
```
GitHub OAuth Config:
  Client ID: Ov23liXwSJQhdRTprTxK
  Redirect URI: https://mail.4w.ink/en/callback
```

### 步骤 3：验证 GitHub App 设置

访问 https://github.com/settings/developers，检查你的 OAuth App：

```
Authorization callback URL: https://mail.4w.ink/en/callback
```

必须与 Docker Compose 中的 `GITHUB_OAUTH_REDIRECT` 完全一致！

---

## 📝 完整的修复步骤

### 1. 更新 docker-compose.yml

使用正确的格式：

```yaml
version: '3.8'

services:
  tmail:
    container_name: omail
    image: ohoimager/omail:develop
    network_mode: host
    restart: unless-stopped
    environment:
      # 数据库
      DB_HOST: dbprovider.ap-northeast-1.clawcloudrun.com
      DB_PORT: 46788
      DB_NAME: tmail
      DB_USER: postgres
      DB_PASS: cv6cklqz
      DB_DRIVER: postgres
      
      # 服务
      HOST: 0.0.0.0
      PORT: 3000
      DOMAIN_LIST: 4w.ink
      TZ: Asia/Shanghai
      ADMIN_ADDRESS: 674904341@4w.ink
      
      # ✅ GitHub OAuth（正确格式）
      GITHUB_OAUTH_ID: Ov23liXwSJQhdRTprTxK
      GITHUB_OAUTH_SECRET: 727a8dda892d74e063fdee8ec605ebdc1c3faa26
      GITHUB_OAUTH_REDIRECT: https://mail.4w.ink/en/callback
      
      # 分析
      UMAMI_ID: e673f3bb-48ce-4388-a7f7-3c5063cdcb84
      UMAMI_URL: https://cloud.umami.is/script.js
      UMAMI_DOMAINS: mail.4w.ink
      PUBLIC_GA_ID: G-5H7JB6P345
      
    volumes:
      - ./tmail:/app/fs
```

### 2. 验证 GitHub App 配置

1. 访问 https://github.com/settings/developers
2. 编辑你的 OAuth App
3. 确认以下设置：
   - **Client ID**: `Ov23liXwSJQhdRTprTxK` ✅
   - **Client Secret**: `727a8dda892d74e063fdee8ec605ebdc1c3faa26` ✅
   - **Authorization callback URL**: `https://mail.4w.ink/en/callback` ✅

### 3. 重启容器

```bash
# 停止旧容器
docker-compose down

# 启动新容器
docker-compose -f docker-compose.yml up -d

# 验证环境变量
docker exec omail sh -c 'echo "Client ID: $GITHUB_OAUTH_ID"'
```

### 4. 测试登录

1. 访问 https://mail.4w.ink
2. 点击 "API Token" 或 "Login with GitHub"
3. 检查跳转的 GitHub 授权 URL 是否包含：
   ```
   client_id=Ov23liXwSJQhdRTprTxK
   redirect_uri=https%3A%2F%2Fmail.4w.ink%2Fen%2Fcallback
   ```

---

## 🔧 调试技巧

### 查看构建的完整命令

```bash
# 查看容器环境变量
docker inspect omail | grep -A 20 Env

# 查看容器运行命令
docker inspect omail | grep -i cmd
```

### 临时测试环境变量

```bash
# 创建临时容器测试
docker run -it \
  -e GITHUB_OAUTH_ID=test123 \
  -e GITHUB_OAUTH_SECRET=secret456 \
  ohoimager/omail:develop \
  sh -c 'echo "ID: $GITHUB_OAUTH_ID, Secret: $GITHUB_OAUTH_SECRET"'
```

### 查看后端源代码的环境变量读取

在 `internal/api/auth.go` 中：
```go
func NewAuthService(client *ent.Client) *AuthService {
	return &AuthService{
		client:       client,
		githubID:     os.Getenv("GITHUB_OAUTH_ID"),      // ← 读取这个
		githubSecret: os.Getenv("GITHUB_OAUTH_SECRET"),  // ← 读取这个
		redirectURL:  os.Getenv("GITHUB_OAUTH_REDIRECT"), // ← 读取这个
	}
}
```

---

## ❓ 常见问题

### Q: 为什么格式很重要？

A: Docker Compose 将环境变量传递给容器的方式很严格。使用冒号 `:` 会被解析为 YAML 键值对，而不是环境变量。

### Q: 如果改用 .env 文件呢？

可以创建 `.env` 文件：
```
GITHUB_OAUTH_ID=Ov23liXwSJQhdRTprTxK
GITHUB_OAUTH_SECRET=727a8dda892d74e063fdee8ec605ebdc1c3faa26
GITHUB_OAUTH_REDIRECT=https://mail.4w.ink/en/callback
```

然后在 docker-compose.yml 中：
```yaml
env_file:
  - .env
```

### Q: 回调 URL 可以是其他值吗？

回调 URL 必须是一个**前端页面**，不能是 API 端点。可选值：
- `https://mail.4w.ink/en/callback` ✅
- `https://mail.4w.ink/zh/callback` ✅
- `https://mail.4w.ink/api/auth/login` ❌ （API 端点）
- `http://localhost:3000/en/callback` ✅ （本地开发）

---

## ✨ 验证成功

修复后，访问 https://mail.4w.ink/en/callback 时应该看到：

```
处理登录中...
正在跳转...
```

然后自动重定向到邮箱界面，localStorage 中会有 `api_token`。

---

## 🔐 安全提示

⚠️ **不要将 Client Secret 提交到 Git！**

使用环境变量管理敏感信息：
```bash
# 不要这样
git add docker-compose.yml  # 包含 secret

# 这样做
git add docker-compose.example.yml
# 在 .gitignore 中添加
echo "docker-compose.yml" >> .gitignore
```

---

**如果问题仍未解决，检查：**
1. Docker 容器是否正在运行：`docker ps | grep omail`
2. 容器日志是否有错误：`docker logs omail`
3. GitHub App 是否正确配置
4. 防火墙是否阻止了外部连接
