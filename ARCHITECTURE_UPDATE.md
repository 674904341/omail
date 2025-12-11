# 🔄 项目架构更新说明

## 变更摘要

**之前：** 所有用户都需要登录才能使用邮箱  
**现在：** 公开使用邮箱，仅登录获取 API Token

---

## 核心变更

### 1️⃣ 前端路由逻辑

**文件：** `web/src/pages/[lang]/index.astro`

```before
// 旧的逻辑：根据登录状态显示不同内容
if (loggedIn) {
  show mailbox
} else {
  show login page
}
```

```after
// 新的逻辑：直接显示邮箱（不检查登录）
show mailbox
// 登录按钮在 Header 中（仅用于获取 API Token）
```

### 2️⃣ 顶部导航栏

**文件：** `web/src/components/Header.astro`

```before
[Logo] [Title] ———— [ThemeIcon]
```

```after
[Logo] [Title] ———— [ThemeIcon] [API Token / Logout]
```

### 3️⃣ 新增认证按钮组件

**文件：** `web/src/components/AuthButton.tsx`

- 显示 "API Token" 登录按钮（未登录时）
- 显示用户名 + "Logout" 按钮（已登录时）
- 管理登录/登出状态
- 处理 GitHub OAuth 回调

---

## API 端点分类

### 🔓 公开端点（无需认证）
```
GET  /api/domain              # 获取可用域名
GET  /api/fetch               # 获取邮件列表
GET  /api/fetch/:id           # 获取邮件详情
GET  /api/fetch/latest        # 获取最新邮件
GET  /api/download/:id        # 下载附件
GET  /api/auth/url            # 获取 GitHub 授权 URL
GET  /api/auth/login          # GitHub 回调处理
GET  /api/report              # 报告垃圾邮件
```

### 🔒 受保护端点（需要 API Token）
```
GET  /api/profile             # 获取用户信息
GET  /api/mailboxes           # 获取用户邮箱列表
POST /api/mailbox             # 创建新邮箱
GET  /api/emails              # 获取邮箱邮件（用户自己的）
GET  /api/email/:id           # 获取邮件详情（用户自己的）
```

**认证方式：**
```bash
Authorization: Bearer YOUR_API_TOKEN
```

---

## 使用场景

### 场景 1：Web 用户（临时邮箱）
```
访问首页
 ↓
使用邮箱功能（无需登录）
 ↓
生成临时邮箱地址
 ↓
接收并查看邮件
```

### 场景 2：API 用户（程序化访问）
```
访问首页
 ↓
点击 "API Token" 按钮
 ↓
GitHub 登录
 ↓
获取 API Token
 ↓
在应用中使用 Token 调用 API
```

### 场景 3：自动化脚本
```
获取 API Token（通过 Web UI 一次性）
 ↓
在脚本环境变量中配置 Token
 ↓
使用 curl/SDK 调用 API
 ↓
自动化完成邮箱相关任务
```

---

## 代码示例

### Web UI 使用（不需要登录）
```javascript
// 获取可用域名
const domains = await fetch('/api/domain').then(r => r.json())

// 获取邮件列表（公开邮箱）
const emails = await fetch('/api/fetch?to=test@mail.4w.ink')
  .then(r => r.json())
```

### API 使用（需要 Token）
```bash
API_TOKEN="your_token_from_web_ui"

# 创建邮箱
curl -X POST \
  -H "Authorization: Bearer $API_TOKEN" \
  http://localhost:3000/api/mailbox

# 获取用户的邮箱列表
curl -H "Authorization: Bearer $API_TOKEN" \
  http://localhost:3000/api/mailboxes

# 获取特定邮箱的邮件
curl -H "Authorization: Bearer $API_TOKEN" \
  "http://localhost:3000/api/emails?email=YOUR_EMAIL@mail.4w.ink"
```

---

## 后端逻辑

### GitHub OAuth 流程（仍然需要）

1. **前端请求认证 URL**
   ```
   GET /api/auth/url?state=xxx
   ```
   响应：`{ "auth_url": "https://github.com/login/oauth/..." }`

2. **用户授权（在 GitHub 上）**

3. **GitHub 回调**
   ```
   GET /api/auth/login?code=xxx&state=xxx
   ```

4. **后端处理**
   - 使用 code 换取 access token
   - 从 GitHub 获取用户信息
   - 创建或更新数据库中的 User 记录
   - 生成新的 API Token
   - 返回 token 给前端

5. **前端存储 Token**
   ```javascript
   localStorage.setItem('api_token', token)
   localStorage.setItem('user', JSON.stringify(userInfo))
   ```

### 认证中间件

对于受保护端点：
```go
// 从请求头获取 token
Authorization: Bearer TOKEN

// 验证 token 是否存在且有效
// 查询数据库中的 APIToken 记录
// 获取关联的 User
// 将 User 信息添加到请求上下文
```

---

## 文件结构对比

### 前 vs 后

```
之前：
web/src/pages/[lang]/
├── index.astro           # 显示登录或邮箱（有条件判断）
└── callback.astro        # GitHub 回调处理

之后：
web/src/pages/[lang]/
├── index.astro           # 直接显示邮箱（无条件）
└── callback.astro        # GitHub 回调处理（不变）

web/src/components/
├── Header.astro          # 新增 AuthButton
├── AuthButton.tsx        # 新增！处理登录/登出
├── Login.tsx             # 仍存在（可用于独立登录页）
└── ...其他组件
```

---

## 数据库变化

**无需数据库迁移！** 现有的数据结构完全兼容：

```sql
-- 用户表（存在）
CREATE TABLE users (
  id INTEGER PRIMARY KEY,
  github_id INTEGER UNIQUE,
  username VARCHAR,
  email VARCHAR,
  avatar_url VARCHAR,
  ...
);

-- API Token 表（存在）
CREATE TABLE api_tokens (
  id INTEGER PRIMARY KEY,
  user_id INTEGER,
  token VARCHAR UNIQUE,
  name VARCHAR,
  created_at TIMESTAMP,
  last_used_at TIMESTAMP,
  revoked BOOLEAN,
  ...
);

-- 邮箱表（存在）
CREATE TABLE mailboxes (
  id INTEGER PRIMARY KEY,
  user_id INTEGER,  -- 可以为 NULL（未登录用户）
  email VARCHAR,
  created_at TIMESTAMP,
  ...
);
```

---

## 迁移指南（如果从旧版本升级）

### 前端
1. 更新 `web/src/pages/[lang]/index.astro`
2. 更新 `web/src/components/Header.astro`
3. 添加 `web/src/components/AuthButton.tsx`

### 后端
- ✅ 无需更改（兼容现有 API）

### 数据库
- ✅ 无需迁移（兼容现有数据）

---

## 测试清单

- [ ] 未登录用户可以访问首页
- [ ] 未登录用户可以使用邮箱功能
- [ ] 顶部显示 "API Token" 按钮（未登录时）
- [ ] 点击按钮触发 GitHub 登录
- [ ] 登录后 Token 存储在 localStorage
- [ ] 登录后按钮变成 "Logout"
- [ ] Logout 清除 localStorage
- [ ] API 调用需要 Bearer Token
- [ ] 无 Token 的 API 调用返回 401

---

## 安全考虑

### 公开端点
- ✅ `/api/domain` - 返回可用域名（安全公开）
- ✅ `/api/fetch` - 任何人都可以查看任何邮箱（因为本就是临时公开的）
- ⚠️ 需要防止滥用（速率限制、验证码等）

### 受保护端点
- ✅ 需要有效的 API Token
- ✅ Token 绑定到特定 User
- ✅ 用户只能访问自己的数据

---

## 配置清单

### 本地开发
```env
# .env
GITHUB_OAUTH_ID=你的_client_id
GITHUB_OAUTH_SECRET=你的_client_secret
GITHUB_OAUTH_REDIRECT=http://localhost:4321
```

### 生产环境
```env
GITHUB_OAUTH_ID=生产_client_id
GITHUB_OAUTH_SECRET=生产_client_secret
GITHUB_OAUTH_REDIRECT=https://mail.4w.ink
```

### GitHub App 配置
- **Homepage URL**: https://mail.4w.ink
- **Authorization callback URL**: https://mail.4w.ink/en/callback

---

## 性能影响

- ✅ **无负面影响** - 公开端点无认证开销
- ✅ **更好的缓存** - 未登录用户的请求可以更激进地缓存
- ✅ **更快的首次加载** - 无需先登录再加载邮箱

---

## 反馈和改进

如有问题或建议，欢迎提交 Issue 或 Pull Request！

关键改进方向：
- [ ] Web UI Token 管理界面
- [ ] Token 过期策略
- [ ] 速率限制防止滥用
- [ ] 更详细的使用统计

---

**更新时间**：2025-12-11
