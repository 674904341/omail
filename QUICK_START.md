# 🚀 本地调试快速指南

## 现在运行中的服务

你的本地开发环境已准备好：

| 服务 | 地址 | 说明 |
|------|------|------|
| **前端** | http://localhost:4321 | Astro 开发服务器 |
| **Mock API** | http://localhost:3000 | 模拟后端（用于调试） |

## 访问应用

1. 打开浏览器：http://localhost:4321
2. 你应该看到 **"Sign in with GitHub"** 登录按钮
3. 点击按钮测试登录流程

## 编辑代码

### 前端代码（自动热加载）
```
web/src/
├── components/        # React 组件
│   ├── Login.tsx      # 登录页面
│   ├── Header.astro   # 顶部导航
│   └── ...
├── pages/             # 页面
│   └── [lang]/
│       ├── index.astro        # 主页
│       └── callback.astro      # GitHub 回调
└── layouts/           # 布局模板
    └── Layout.astro
```

**编辑任何文件后会自动刷新浏览器！**

### 模拟 API（需要重启）
```
mock-server.js        # 修改此文件后需要重启 Mock 服务器
```

## 常用操作

### 1. 清除 localStorage（重置登录状态）
```javascript
// 在浏览器控制台执行：
localStorage.clear()
location.reload()
```

### 2. 模拟用户登录
```javascript
// 在浏览器控制台执行：
localStorage.setItem('api_token', 'test_token_123')
localStorage.setItem('user', JSON.stringify({
  id: 1,
  username: 'test_user',
  email: 'test@example.com',
  avatar_url: 'https://...'
}))
location.reload()
```

### 3. 测试 API 调用
```bash
# 获取邮箱列表
curl -H "Authorization: Bearer test_token_123" \
  http://localhost:3000/api/mailboxes

# 创建邮箱
curl -X POST \
  -H "Authorization: Bearer test_token_123" \
  http://localhost:3000/api/mailbox
```

### 4. 查看浏览器控制台
按 `F12` 或 `Ctrl+Shift+I`：
- **Console** 标签：查看日志和错误
- **Network** 标签：查看 API 请求/响应
- **Storage** 标签：查看 localStorage

## 问题排除

### 前端没有显示登录按钮？
1. 打开浏览器 F12 → Console
2. 检查是否有 JavaScript 错误
3. 检查 http://localhost:4321 是否正常加载

### API 请求失败？
1. 确认 Mock 服务器在运行：http://localhost:3000/api/domain
2. 检查浏览器 F12 → Network 标签
3. 查看响应状态码和错误信息

### 如何停止服务？
在服务窗口中按 `Ctrl+C`

## 文件说明

| 文件 | 说明 |
|-----|------|
| `LOCAL_DEBUG.md` | 详细调试指南 |
| `GITHUB_LOGIN_SETUP.md` | GitHub OAuth 配置 |
| `mock-server.js` | 模拟 API 服务器 |
| `.env` | 环境变量配置 |
| `dev-start.ps1` | 快速启动脚本（下次使用） |

## 下次启动

使用快速启动脚本：
```bash
.\dev-start.ps1
```

它会自动：
- 检查依赖
- 启动 Mock API 服务器
- 启动前端开发服务器

## 核心代码位置

### 登录流程
- 前端：`web/src/components/Login.tsx`
- 回调：`web/src/pages/[lang]/callback.astro`
- 主页逻辑：`web/src/pages/[lang]/index.astro`

### Mock API
- `mock-server.js` - 所有 API 端点的实现

### 后端（生产环境）
- `internal/api/handlers.go` - API 处理器
- `internal/api/auth.go` - GitHub OAuth 逻辑
- `internal/route/route.go` - 路由配置

## 后续步骤

1. **配置 GitHub OAuth**（仅实际登录时需要）
   - 访问 https://github.com/settings/developers
   - 创建 OAuth App
   - 编辑 `.env` 填入凭证

2. **测试实际登录**
   - 使用真实的 GitHub OAuth（需要上面的步骤）
   - 会调用真实的 GitHub 服务

3. **连接真实后端**
   - 使用完整的 Docker Compose 或安装 Go 环境
   - 修改前端请求地址

## 需要帮助？

查看详细文档：
- `LOCAL_DEBUG.md` - 完整调试指南
- `GITHUB_LOGIN_SETUP.md` - GitHub 认证配置
- `API_AUTH.md` - API 文档
