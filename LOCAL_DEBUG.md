# 本地前端调试指南

## 快速开始

### 方案 1：前端 + Docker 后端（推荐）

**优点**：后端与生产环境一致，不需要安装 Go

#### 1. 准备 GitHub OAuth 凭证

1. 访问 https://github.com/settings/developers
2. 创建新的 OAuth App：
   - Application name: `Tmail Dev`
   - Homepage URL: `http://localhost:3000`
   - Authorization callback URL: `http://localhost:3000/en/callback`
3. 获取 `Client ID` 和 `Client Secret`

#### 2. 启动后端服务

```bash
# 复制环境变量文件
cp .env.example .env

# 编辑 .env，填入 GitHub 凭证
# GITHUB_OAUTH_ID=xxx
# GITHUB_OAUTH_SECRET=xxx

# 启动 Docker Compose（包含数据库和后端）
docker-compose -f docker-compose.dev.yml up -d
```

等待数据库就绪（通常需要 10-20 秒）：
```bash
# 检查日志
docker-compose -f docker-compose.dev.yml logs -f tmail
```

#### 3. 启动前端开发服务器

```bash
cd web
npm install  # 如果需要
npm run dev
```

#### 4. 访问应用

打开浏览器访问：
- **前端**：http://localhost:4321/
- **后端 API**：http://localhost:3000/api/
- **数据库**：localhost:5432 (PostgreSQL)

### 方案 2：完整本地开发（需要 Go）

如果你安装了 Go 1.24+：

```bash
# 终端 1: 启动数据库
docker run --name tmail-postgres \
  -e POSTGRES_DB=tmail \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine

# 终端 2: 启动后端
go run cmd/main.go

# 终端 3: 启动前端
cd web && npm run dev
```

## 调试技巧

### 1. 查看浏览器控制台

按 `F12` 打开开发者工具：
- **Console** 标签：查看 JavaScript 错误
- **Network** 标签：查看 API 调用
- **Storage** 标签：查看 localStorage（api_token 等）

### 2. 查看后端日志

```bash
# 查看 Docker 日志
docker-compose -f docker-compose.dev.yml logs -f tmail

# 或直接查看最后 50 行
docker-compose -f docker-compose.dev.yml logs --tail 50 tmail
```

### 3. 测试 API 端点

使用 curl 或 Postman 测试：

```bash
# 获取 GitHub 授权 URL
curl "http://localhost:3000/api/auth/url?state=test123"

# 返回示例
# {"auth_url":"https://github.com/login/oauth/authorize?client_id=xxx&redirect_uri=http%3A%2F%2Flocalhost%3A3000&scope=user&state=test123"}
```

### 4. 检查数据库

连接到 PostgreSQL：

```bash
# 使用 psql（需要安装）
psql -h localhost -U postgres -d tmail

# 或使用 Docker
docker-compose -f docker-compose.dev.yml exec postgres psql -U postgres -d tmail
```

## 常见问题

### 问题 1：前端无法连接后端 API

**症状**：浏览器控制台显示 CORS 或 404 错误

**解决**：
1. 检查后端是否运行：`curl http://localhost:3000/api/domain`
2. 检查 API 路由：确认 `/api/auth/url` 端点可访问
3. 检查 CORS 配置：可能需要在后端添加 CORS 中间件

### 问题 2：GitHub 登录失败

**症状**：点击登录后白屏或回到首页

**解决**：
1. 检查 GitHub OAuth App 设置
2. 确认 Authorization callback URL 正确：`http://localhost:3000/en/callback`
3. 检查浏览器控制台错误信息
4. 查看后端日志：`docker-compose -f docker-compose.dev.yml logs tmail`

### 问题 3：数据库连接失败

**症状**：后端启动失败，日志显示数据库连接错误

**解决**：
```bash
# 检查 PostgreSQL 是否就绪
docker-compose -f docker-compose.dev.yml ps

# 查看完整日志
docker-compose -f docker-compose.dev.yml logs postgres

# 重启数据库
docker-compose -f docker-compose.dev.yml restart postgres
```

### 问题 4：端口已被占用

**症状**：启动服务时显示端口已占用

**解决**：
```bash
# 查看占用 3000 端口的进程
lsof -i :3000  # 或在 Windows 上：netstat -ano | findstr :3000

# 停止现有服务
docker-compose -f docker-compose.dev.yml down

# 或更改端口（修改 docker-compose.dev.yml）
```

## 停止和清理

```bash
# 停止所有服务
docker-compose -f docker-compose.dev.yml down

# 删除所有数据（完全清理）
docker-compose -f docker-compose.dev.yml down -v

# 删除 Docker 镜像
docker-compose -f docker-compose.dev.yml down --rmi all
```

## 环境变量配置

创建 `.env` 文件（参考 `.env.example`）：

```env
# GitHub OAuth（必需）
GITHUB_OAUTH_ID=xxx
GITHUB_OAUTH_SECRET=xxx
GITHUB_OAUTH_REDIRECT=http://localhost:3000

# 数据库（可选，使用 docker-compose 自动配置）
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=postgres
DB_NAME=tmail

# 分析工具（可选）
PUBLIC_GA_ID=G-XXXXXXXXXX
PUBLIC_UMAMI_DOMAIN=analytics.example.com
PUBLIC_UMAMI_SCRIPT_URL=https://analytics.example.com/script.js
```

## 提示

- ✅ 前端代码修改会自动热加载（Astro dev server）
- ✅ 后端代码修改需要重启容器：`docker-compose -f docker-compose.dev.yml restart tmail`
- ✅ 使用 `localhost` 而不是 `127.0.0.1`（某些情况下 CORS 表现不同）
- ✅ 定期清理 Docker 数据以获得最佳性能：`docker system prune`
