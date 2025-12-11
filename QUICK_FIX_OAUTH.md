# 🔧 快速修复：GitHub OAuth 环境变量问题

## 问题摘要
Docker Compose 中的 GitHub OAuth 环境变量格式错误，导致认证参数为空。

---

## ⚡ 快速修复（3步）

### 步骤 1：修正 docker-compose.yml

**将你的配置从：**
```yaml
environment:
  - "GITHUB_OAUTH_ID: Ov23liXwSJQhdRTprTxK"
  - "GITHUB_OAUTH_SECRET: 727a8dda..."
  - "GITHUB_OAUTH_REDIRECT: https://mail.4w.ink/api/auth/login"
```

**改为：**
```yaml
environment:
  GITHUB_OAUTH_ID: Ov23liXwSJQhdRTprTxK
  GITHUB_OAUTH_SECRET: 727a8dda892d74e063fdee8ec605ebdc1c3faa26
  GITHUB_OAUTH_REDIRECT: https://mail.4w.ink/en/callback
```

**关键修改：**
- ✅ 移除引号（`"..."` → 无引号）
- ✅ 用冒号分隔（但不在 `-` 列表中）
- ✅ 回调 URL：`/api/auth/login` → `/en/callback`

### 步骤 2：验证 GitHub App 设置

访问 https://github.com/settings/developers，编辑你的 OAuth App：

确认 **Authorization callback URL** 设置为：
```
https://mail.4w.ink/en/callback
```

### 步骤 3：重启容器

```bash
# 停止旧容器
docker-compose down

# 启动新容器
docker-compose up -d

# 验证配置已读取
docker logs omail | grep -i github
```

应该看到类似的日志：
```
Application config loaded github_id_set=true github_secret_set=true github_redirect=https://mail.4w.ink/en/callback
```

---

## 📝 完整的修复后配置文件

```yaml
version: '3.8'

services:
  tmail:
    container_name: omail
    image: ohoimager/omail:develop
    network_mode: host
    restart: unless-stopped
    environment:
      # 数据库配置
      DB_HOST: dbprovider.ap-northeast-1.clawcloudrun.com
      DB_PORT: 46788
      DB_NAME: tmail
      DB_USER: postgres
      DB_PASS: cv6cklqz
      DB_DRIVER: postgres
      
      # 服务配置
      HOST: 0.0.0.0
      PORT: 3000
      DOMAIN_LIST: 4w.ink
      TZ: Asia/Shanghai
      ADMIN_ADDRESS: 674904341@4w.ink
      
      # ✅ GitHub OAuth（正确格式）
      GITHUB_OAUTH_ID: Ov23liXwSJQhdRTprTxK
      GITHUB_OAUTH_SECRET: 727a8dda892d74e063fdee8ec605ebdc1c3faa26
      GITHUB_OAUTH_REDIRECT: https://mail.4w.ink/en/callback
      
      # 分析配置
      UMAMI_ID: e673f3bb-48ce-4388-a7f7-3c5063cdcb84
      UMAMI_URL: https://cloud.umami.is/script.js
      UMAMI_DOMAINS: mail.4w.ink
      PUBLIC_GA_ID: G-5H7JB6P345
      
    volumes:
      - ./tmail:/app/fs
```

---

## 🧪 测试

修复后：

1. **访问应用**：https://mail.4w.ink
2. **点击 API Token 按钮**
3. **检查跳转 URL**：应该包含完整的 client_id 和 redirect_uri
   ```
   https://github.com/login/oauth/authorize?
   client_id=Ov23liXwSJQhdRTprTxK&
   redirect_uri=https%3A%2F%2Fmail.4w.ink%2Fen%2Fcallback&
   scope=user&state=...
   ```

4. **完成授权后**：应该自动回到 https://mail.4w.ink 并显示用户名

---

## ❌ 常见错误对照表

| 错误 | 原因 | 修复 |
|------|------|------|
| `client_id=&redirect_uri=` | 环境变量为空 | 检查 docker-compose 格式 |
| 回调失败白屏 | 回调 URL 错误 | 改为 `/en/callback` |
| `connection refused` | 容器未运行 | `docker-compose up -d` |
| 日志显示环境变量为空 | 格式错误 | 使用冒号分隔，不用引号 |

---

## 📖 详细参考

更完整的故障排除指南见：`GITHUB_OAUTH_TROUBLESHOOTING.md`

**包含内容：**
- ✅ 完整的诊断步骤
- ✅ 查看容器日志的方法
- ✅ 安全最佳实践
- ✅ 常见问题解答
