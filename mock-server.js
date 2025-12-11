/**
 * 本地开发模式下的 API 模拟服务器
 * 用于前端调试，无需启动完整的后端
 * 
 * 使用方法：
 * npm install express cors dotenv
 * node mock-server.js
 */

const express = require("express")
const cors = require("cors")
require("dotenv").config()

const app = express()
app.use(cors())
app.use(express.json())

// 存储模拟用户和 token
let users = {}
let tokens = {}

// GitHub OAuth 配置
const GITHUB_CLIENT_ID = process.env.GITHUB_OAUTH_ID || "your_client_id"
const GITHUB_CLIENT_SECRET = process.env.GITHUB_OAUTH_SECRET || "your_client_secret"
const GITHUB_REDIRECT_URI = process.env.GITHUB_OAUTH_REDIRECT || "http://localhost:3000"

// 生成随机 token
function generateToken() {
  return Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15)
}

// 生成随机邮箱
function generateEmail() {
  const random = Math.random().toString(36).substring(2, 10)
  return `${random}@mail.4w.ink`
}

// 1. 获取 GitHub 授权 URL
app.get("/api/auth/url", (req, res) => {
  const state = req.query.state
  if (!state) {
    return res.status(400).json({ error: "state parameter required" })
  }

  const authUrl = `https://github.com/login/oauth/authorize?client_id=${GITHUB_CLIENT_ID}&redirect_uri=${encodeURIComponent(GITHUB_REDIRECT_URI)}&scope=user&state=${state}`

  res.json({ auth_url: authUrl })
})

// 2. GitHub 回调处理（模拟）
app.get("/api/auth/login", (req, res) => {
  const code = req.query.code
  const state = req.query.state

  if (!code) {
    return res.status(400).json({ error: "code parameter required" })
  }

  // 模拟 GitHub 用户信息
  const userId = Math.floor(Math.random() * 1000000)
  const username = `user_${userId}`
  const apiToken = generateToken()

  // 存储用户和 token
  users[apiToken] = {
    id: userId,
    username: username,
    email: `${username}@github.com`,
    avatar_url: `https://avatars.githubusercontent.com/u/${userId}?v=4`,
  }

  tokens[apiToken] = {
    user_id: userId,
    created_at: new Date().toISOString(),
  }

  res.json({
    api_token: apiToken,
    user: users[apiToken],
  })
})

// 3. 获取用户信息
app.get("/api/profile", (req, res) => {
  const token = extractToken(req)
  if (!token || !users[token]) {
    return res.status(401).json({ error: "unauthorized" })
  }

  res.json(users[token])
})

// 4. 获取邮箱列表
app.get("/api/mailboxes", (req, res) => {
  const token = extractToken(req)
  if (!token || !users[token]) {
    return res.status(401).json({ error: "unauthorized" })
  }

  // 返回模拟的邮箱列表
  res.json({
    mailboxes: [
      {
        email: generateEmail(),
        created_at: new Date(Date.now() - 3600000).toISOString(),
      },
      {
        email: generateEmail(),
        created_at: new Date(Date.now() - 7200000).toISOString(),
      },
    ],
  })
})

// 5. 创建新邮箱
app.post("/api/mailbox", (req, res) => {
  const token = extractToken(req)
  if (!token || !users[token]) {
    return res.status(401).json({ error: "unauthorized" })
  }

  const email = generateEmail()
  res.json({
    email: email,
    created_at: new Date().toISOString(),
  })
})

// 6. 获取邮件列表
app.get("/api/emails", (req, res) => {
  const token = extractToken(req)
  if (!token || !users[token]) {
    return res.status(401).json({ error: "unauthorized" })
  }

  const email = req.query.email
  if (!email) {
    return res.status(400).json({ error: "email parameter required" })
  }

  // 返回模拟的邮件列表
  res.json({
    emails: [
      {
        id: 1,
        from: "noreply@example.com",
        subject: "[Verify] Email Verification Code",
        created_at: new Date().toISOString(),
      },
      {
        id: 2,
        from: "support@service.com",
        subject: "Welcome to our service!",
        created_at: new Date(Date.now() - 3600000).toISOString(),
      },
    ],
  })
})

// 7. 获取邮件详情
app.get("/api/email/:id", (req, res) => {
  const token = extractToken(req)
  if (!token || !users[token]) {
    return res.status(401).json({ error: "unauthorized" })
  }

  const emailId = req.params.id
  res.json({
    id: emailId,
    from: "noreply@example.com",
    to: "test@mail.4w.ink",
    subject: "[Verify] Email Verification Code",
    content: "<p>Your verification code is: <strong>123456</strong></p>",
    created_at: new Date().toISOString(),
    attachments: [
      {
        id: 1,
        filename: "document.pdf",
        content_type: "application/pdf",
      },
    ],
  })
})

// 原始的域名列表端点
app.get("/api/domain", (req, res) => {
  res.json({
    domains: ["mail.4w.ink", "localhost"],
  })
})

// Helper: 从请求头提取 token
function extractToken(req) {
  const authHeader = req.headers.authorization
  if (!authHeader || !authHeader.startsWith("Bearer ")) {
    return null
  }
  return authHeader.substring(7)
}

// 启动服务器
const PORT = process.env.API_PORT || 3000
app.listen(PORT, () => {
  console.log(`
  ╔════════════════════════════════════════════════════════════════╗
  ║                  Mock API Server Started                        ║
  ╠════════════════════════════════════════════════════════════════╣
  ║  API Server:  http://localhost:${PORT}                           ║
  ║  Frontend:    http://localhost:4321                            ║
  ║                                                                ║
  ║  Available endpoints:                                          ║
  ║    GET  /api/auth/url                                         ║
  ║    GET  /api/auth/login                                       ║
  ║    GET  /api/profile                                          ║
  ║    GET  /api/mailboxes                                        ║
  ║    POST /api/mailbox                                          ║
  ║    GET  /api/emails                                           ║
  ║    GET  /api/email/:id                                        ║
  ║    GET  /api/domain                                           ║
  ╚════════════════════════════════════════════════════════════════╝

  GitHub OAuth Config:
    Client ID:     ${GITHUB_CLIENT_ID}
    Redirect URI:  ${GITHUB_REDIRECT_URI}

  Note: This is a mock server for development only.
        For production, use the real Go backend.
  `)
})
