# .gitignore 配置说明

## 已排除的开发环境文件

本项目的 `.gitignore` 已配置为排除以下本地开发文件，避免将它们提交到 GitHub：

### 环境变量文件
- `.env` - 主环境变量文件（包含敏感信息）
- `.env.local` - 本地环境变量
- `.env.*.local` - 特定环境的本地变量

### Node.js 依赖
- `node_modules/` - npm 包依赖目录
- `package-lock.json`（如果使用）

### 日志文件
- `*.log`
- `npm-debug.log*`
- `yarn-debug.log*`
- `yarn-error.log*`
- `pnpm-debug.log*`

### IDE 和编辑器
- `.vscode/` - VS Code 配置
- `.idea/` - IntelliJ IDEA 配置
- `*.swp`, `*.swo` - Vim 交换文件
- `*~` - 备份文件
- `.sublime-project`、`.sublime-workspace` - Sublime Text 配置

### 操作系统文件
- `.DS_Store` - macOS 文件
- `Thumbs.db` - Windows 缩略图缓存

### 构建输出
- `web/dist/` - 前端构建输出（Astro）
- `web/.astro/` - Astro 缓存
- `dist/` - Go 后端构建输出

### 临时文件
- `*.pid` - 进程 ID 文件
- `*.seed` - 数据库种子文件
- `.cache/` - 缓存目录
- `.parcel-cache/` - Parcel 缓存

### Mock 服务器相关
- `mock-server.js` 本身不被忽略（作为开发工具保留），但运行时生成的文件会被忽略

## 如何检查被忽略的文件

```bash
# 查看所有被 .gitignore 忽略的文件
git check-ignore -v <文件名>

# 查看哪些本应被忽略的文件被追踪了
git ls-files --others --ignored --exclude-standard
```

## 从 Git 中移除已追踪的文件

如果某些开发文件已经被 Git 追踪，需要移除它们：

```bash
# 从 Git 追踪中移除（保留本地文件）
git rm --cached <文件名>

# 从 Git 追踪中移除整个目录
git rm --cached -r <目录名>

# 更新提交
git commit -m "Remove development files from tracking"
```

## 建议：添加本地配置的提交

如果你需要本地特定配置（不想提交给其他开发者），有以下选项：

### 1. 使用 `.env.example` 模板
在仓库中保留 `.env.example`：
```bash
# 项目根目录
cp .env.example .env
# 编辑 .env 填入本地配置
```

### 2. 使用 Git 本地配置
```bash
# 设置本地 Git 配置（不会提交）
git config user.email "your-email@example.com"
```

### 3. 使用 `.git/info/exclude`
对于只在你本地使用的忽略规则：
```bash
echo "your-local-file" >> .git/info/exclude
```

## 安全建议

⚠️ **重要：**
- 从不提交 `.env` 文件，特别是包含 GitHub OAuth Secret
- 从不提交包含密码或 API Key 的文件
- 始终在 `.env.example` 中使用占位符值
- 定期检查 `git status` 确保没有敏感信息被追踪

## 团队协作

对于团队开发，建议：

1. 所有开发者都创建自己的 `.env` 文件（从 `.env.example` 复制）
2. `.env.example` 在仓库中，包含所有必需的变量名（带占位符值）
3. 每个开发者在本地 `.env` 中填入自己的凭证
4. 通过文档告知如何设置（在 `LOCAL_DEBUG.md` 中已包含）

当前项目已在 `LOCAL_DEBUG.md` 和 `GITHUB_LOGIN_SETUP.md` 中提供了详细的本地配置指南。
