# 新的工作流程：无需登录使用邮箱，登录获取 API Token

## 🎯 变更说明

项目已调整为两种使用方式：

### **方式 1：Web 界面（无需登录）** ✅ 推荐大多数用户
- 访问 http://localhost:4321 或 https://mail.4w.ink
- **直接使用**临时邮箱功能
- 生成临时邮箱地址
- 接收和查看邮件
- **不需要登录**

### **方式 2：API 访问（需要 API Token）** ✅ 用于程序化访问
- 点击顶部导航栏的 **"API Token"** 按钮
- 使用 GitHub 登录获取 API Token
- 使用 Token 通过 API 调用所有功能
- 示例：
  ```bash
  curl -H "Authorization: Bearer YOUR_API_TOKEN" \
    http://localhost:3000/api/mailboxes
  ```

---

## 🔄 用户流程

### 普通用户
```
访问首页 → 使用邮箱功能 → 完成（不需要登录）
```

### API 用户
```
访问首页 → 点击"API Token"按钮 → GitHub 登录 → 获取 Token → 使用 API
```

---

## 📁 文件变更

### 前端更改

#### 1. **主页 (`[lang]/index.astro`)**
- ✅ 移除了登录检查逻辑
- ✅ 所有用户都可以访问邮箱功能
- ✅ 简化为直接显示 Content 组件

#### 2. **Header (`Header.astro`)**
- ✅ 添加了新的 `AuthButton` 组件
- ✅ 顶部导航栏显示"API Token"登录按钮
- ✅ 登录后显示用户名和登出按钮

#### 3. **新增 AuthButton 组件 (`AuthButton.tsx`)**
```tsx
// 功能：
// - 显示登录/登出按钮
// - 管理 API Token 状态
// - 处理 GitHub OAuth 登录
// - 显示当前登录用户
```

### 后端保持不变
- `/api/auth/url` - 获取 GitHub 授权链接（公开）
- `/api/auth/login` - GitHub 回调处理（公开）
- `/api/profile` - 需要 API Token
- `/api/mailboxes` - 需要 API Token
- `/api/domain` - 公开（用于获取可用域名）

---

## 🔐 API Token 说明

### 获取 Token

1. **点击顶部导航栏的 "API Token" 按钮**
2. **用 GitHub 账户登录**
3. **授权应用**
4. **系统自动生成并存储 Token**

### 使用 Token

```bash
# 获取用户信息
curl -H "Authorization: Bearer TOKEN_HERE" \
  http://localhost:3000/api/profile

# 获取邮箱列表
curl -H "Authorization: Bearer TOKEN_HERE" \
  http://localhost:3000/api/mailboxes

# 创建新邮箱
curl -X POST \
  -H "Authorization: Bearer TOKEN_HERE" \
  http://localhost:3000/api/mailbox

# 获取邮件列表
curl -H "Authorization: Bearer TOKEN_HERE" \
  "http://localhost:3000/api/emails?email=xxx@mail.4w.ink"
```

### Token 存储位置
- **浏览器**: `localStorage['api_token']`
- **后端**: `APIToken` 表（关联到 User）

---

## 💾 localStorage 结构

**未登录时：**
```javascript
localStorage.api_token        // 不存在
localStorage.user             // 不存在
```

**已登录时：**
```javascript
localStorage.api_token = "abc123def456..."

localStorage.user = JSON.stringify({
  id: 123,
  username: "github_username",
  email: "user@github.com",
  avatar_url: "https://..."
})
```

---

## 🧪 测试流程

### 测试Web界面（无需登录）
1. 打开 http://localhost:4321
2. 应该直接看到邮箱界面
3. 可以创建临时邮箱
4. 可以查看邮件

### 测试 API Token 获取
1. 点击顶部 "API Token" 按钮
2. 重定向到 GitHub 登录
3. 授权应用
4. 重定向回应用
5. Token 自动存储在 localStorage
6. 按钮变成 "Logout"

### 测试 API 访问
```bash
# 在浏览器控制台获取 Token
const token = localStorage.getItem('api_token')

# 使用 curl 测试
curl -H "Authorization: Bearer $token" \
  http://localhost:3000/api/profile
```

---

## 🔄 GitHub Callback 流程

流程不变，但现在是可选的：

```
用户点击"API Token"按钮
  ↓
/api/auth/url?state=xxx
  ↓ 获取 GitHub 授权 URL
https://github.com/login/oauth/authorize?...
  ↓ 用户授权
GitHub 回调 → /en/callback?code=xxx&state=xxx
  ↓
/api/auth/login?code=xxx&state=xxx
  ↓ 获取 access token，创建/更新用户，生成 API Token
  ↓
存储到 localStorage，显示成功信息
```

---

## 📊 访问权限矩阵

| 功能 | 需登录 | API Token | 说明 |
|------|---------|-----------|------|
| 查看邮箱界面 | ❌ | ✅ | 任何人都可以访问 |
| 创建临时邮箱 | ❌ | ✅ | Web UI：免费；API：需 Token |
| 查看邮件 | ❌ | ✅ | Web UI：免费；API：需 Token |
| 获取邮箱列表 | - | ✅ | 仅 API（需 Token）|
| 下载附件 | ❌ | ✅ | 可选功能 |
| 用户账户管理 | ✅ | ✅ | 登录后可见 |

---

## 🚀 部署注意事项

### 生产环境配置
在 `docker-compose.yml` 中：
```yaml
environment:
  GITHUB_OAUTH_ID: your_prod_client_id
  GITHUB_OAUTH_SECRET: your_prod_client_secret
  GITHUB_OAUTH_REDIRECT: https://mail.4w.ink
```

### GitHub App 配置
- **Homepage URL**: https://mail.4w.ink
- **Authorization callback URL**: https://mail.4w.ink/en/callback

---

## 🔒 安全考虑

### API Token
- ✅ 由后端随机生成
- ✅ 存储在 User 记录中
- ✅ 支持多个 Token 创建
- ✅ 可以撤销（未来功能）
- ✅ 记录最后使用时间

### 匿名使用（Web UI）
- ✅ 无法访问 API
- ✅ 无法导出数据
- ✅ 邮箱数据临时存储
- ✅ 无登录状态

---

## 💡 使用案例

### 案例 1：个人测试邮箱（Web）
```
1. 打开 https://mail.4w.ink
2. 生成邮箱：test123@mail.4w.ink
3. 用于注册测试账户
4. 接收验证邮件
5. 查看邮件内容
```

### 案例 2：自动化脚本（API）
```bash
# 1. 获取 API Token（一次性）
# 通过 Web UI 登录，复制 Token

# 2. 在脚本中使用
API_TOKEN="your_token_here"

# 创建邮箱
curl -X POST \
  -H "Authorization: Bearer $API_TOKEN" \
  http://api.example.com/api/mailbox

# 获取邮件
curl -H "Authorization: Bearer $API_TOKEN" \
  "http://api.example.com/api/emails?email=test@mail.4w.ink"

# 自动化流程...
```

### 案例 3：第三方集成（API）
```javascript
// Node.js 客户端库
const TmailAPI = require('tmail-api');

const client = new TmailAPI({
  token: process.env.TMAIL_API_TOKEN
});

// 创建邮箱
const mailbox = await client.mailbox.create();
console.log(`新邮箱：${mailbox.email}`);

// 轮询邮件
const emails = await client.emails.list(mailbox.email);
```

---

## 常见问题

**Q: 我想在 Web 上进行数据分析，需要登录吗？**  
A: 不需要。Web 界面对所有人开放。

**Q: 我想用 API 实现自动化，如何开始？**  
A: 
1. 访问应用，点击 "API Token" 按钮
2. GitHub 登录获取 Token
3. 在 API 请求中使用 Bearer Token

**Q: Token 会过期吗？**  
A: 当前实现中不会过期。生产环境可以添加过期时间。

**Q: 如何撤销 Token？**  
A: 当前需要数据库操作。未来可以添加 Web UI 管理界面。

---

## 下一步改进

- [ ] Web UI 中的 API Token 管理界面
- [ ] Token 过期时间设置
- [ ] Token 使用统计
- [ ] 多个 Token 支持
- [ ] Token 撤销功能
- [ ] API 速率限制

---

## 快速参考

| 任务 | 方式 | 说明 |
|------|------|------|
| 使用邮箱 | Web | 打开应用，直接使用 |
| 获取 Token | Web | 点击 "API Token" 按钮 |
| 创建邮箱（API）| API | `POST /api/mailbox` + Token |
| 查看邮件（API）| API | `GET /api/emails?email=xxx` + Token |
| 用户信息 | API | `GET /api/profile` + Token |

---

**最后更新**：2025-12-11
