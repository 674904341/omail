# ✅ 本地前端调试设置完成

## 🎯 现在可以做什么

### 1. **查看前端页面**
- 访问 http://localhost:4321
- 你会看到登录页面或邮箱界面（取决于 localStorage）

### 2. **调试登录功能**
前端已集成完整的 GitHub OAuth 登录流程：
- ✅ Login 组件（显示登录按钮）
- ✅ Callback 页面（处理 GitHub 回调）
- ✅ 主页路由逻辑（根据登录状态显示不同内容）
- ✅ Token 存储在 localStorage

### 3. **测试 Mock API**
Mock 服务器提供所有必需的端点：
- ✅ `/api/auth/url` - 获取 GitHub 授权 URL
- ✅ `/api/auth/login` - 处理登录回调
- ✅ `/api/profile` - 获取用户信息
- ✅ `/api/mailboxes` - 获取邮箱列表
- ✅ `/api/mailbox` - 创建邮箱
- ✅ `/api/emails` - 获取邮件列表
- ✅ `/api/email/:id` - 获取邮件详情

## 🔧 调试技巧

### 快速测试登录状态
在浏览器控制台（F12）执行：
```javascript
// 模拟已登录
localStorage.setItem('api_token', 'test_token_xyz')
localStorage.setItem('user', JSON.stringify({
  id: 123,
  username: 'testuser',
  email: 'test@github.com',
  avatar_url: 'https://avatars.githubusercontent.com/u/123?v=4'
}))
location.reload()

// 清除登录状态
localStorage.clear()
location.reload()
```

### 监控 API 请求
1. 打开浏览器 F12
2. 点击 **Network** 标签
3. 执行登录或其他操作
4. 查看请求和响应详情

### 查看应用日志
1. 打开浏览器 F12
2. 点击 **Console** 标签
3. 查看所有 console.log 和错误信息

## 📝 文件结构回顾

```
前端项目:
  web/
  ├── src/
  │   ├── components/
  │   │   ├── Login.tsx                 ← 登录组件（GitHub 按钮）
  │   │   ├── Header.astro              ← 顶部导航
  │   │   ├── DisclosureModal.tsx       ← 免责声明弹窗
  │   │   └── ...
  │   ├── pages/
  │   │   └── [lang]/
  │   │       ├── index.astro           ← 主页（路由逻辑）
  │   │       └── callback.astro        ← GitHub 回调页面 ✨
  │   ├── layouts/
  │   │   └── Layout.astro
  │   └── lib/
  │       └── store/
  │           └── store.ts              ← 应用状态管理
  │
  ├── astro.config.mjs
  └── package.json

后端 Mock:
  ├── mock-server.js                    ← Mock API 服务器
  └── .env                              ← 环境变量

文档:
  ├── QUICK_START.md                    ← 快速开始（👈 你在这里）
  ├── LOCAL_DEBUG.md                    ← 详细调试指南
  ├── GITHUB_LOGIN_SETUP.md             ← GitHub OAuth 配置
  └── API_AUTH.md                       ← API 文档
```

## 🚀 下一步计划

### 第 1 阶段：本地调试（现在）✅
- [x] 前端开发服务器运行
- [x] Mock API 服务器运行
- [x] 登录流程设计完成
- [ ] **→ 下一步：编辑前端代码**

### 第 2 阶段：集成真实后端
- [ ] 使用 Docker Compose 启动完整后端
- [ ] 测试实际的 GitHub OAuth
- [ ] 测试数据库操作

### 第 3 阶段：部署到生产
- [ ] 配置 Docker 镜像构建
- [ ] 设置 GitHub Actions CI/CD
- [ ] 部署到服务器

## 💡 常见的开发任务

### 修改登录按钮样式
编辑：`web/src/components/Login.tsx`
```tsx
// 修改这部分来改变按钮外观
<Button
  onClick={handleLogin}
  disabled={!authUrl || loading}
  className="w-full bg-slate-900 hover:bg-slate-800"  // ← 这里
>
```

### 修改首页显示逻辑
编辑：`web/src/pages/[lang]/index.astro`
```astro
<script>
  // 修改这里来改变登录/邮箱界面的显示条件
  const token = localStorage.getItem("api_token")
  // ...
</script>
```

### 添加新的 API 端点
编辑：`mock-server.js`
```javascript
// 添加新的路由
app.get("/api/new-endpoint", (req, res) => {
  const token = extractToken(req)
  if (!token || !users[token]) {
    return res.status(401).json({ error: "unauthorized" })
  }
  // 你的逻辑...
})
```

## 🐛 常见问题

**Q: 点击登录按钮没有反应？**  
A: 检查浏览器 F12 → Console 是否有错误。确保 Mock API 服务器在运行。

**Q: 如何模拟不同的 GitHub 用户？**  
A: 修改 `mock-server.js` 中的用户信息生成逻辑。

**Q: 前端文件修改后没有自动刷新？**  
A: Astro dev server 应该自动检测变化。如果没有，手动刷新浏览器。

**Q: 如何查看 localStorage 中的数据？**  
A: 打开 F12 → Storage 标签 → LocalStorage。

## 📞 需要帮助？

- 查看 `LOCAL_DEBUG.md` 获取完整的故障排除指南
- 查看 `GITHUB_LOGIN_SETUP.md` 了解 GitHub OAuth 配置
- 查看 `API_AUTH.md` 了解 API 详细信息

## ✨ 下次快速启动

只需运行：
```powershell
.\dev-start.ps1
```

它会自动启动所有必要的服务！

---

**祝你调试愉快！🎉**
