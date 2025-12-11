#!/usr/bin/env pwsh
# 本地前端调试快速启动脚本
# 使用方法：./dev-start.ps1

Write-Host "╔════════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║        Tmail 本地开发环境启动脚本                                ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════════════════════════════╝" -ForegroundColor Cyan

# 检查依赖
Write-Host "`n[1/3] 检查依赖..." -ForegroundColor Yellow

if (-Not (Test-Path "web/package.json")) {
    Write-Host "✗ web/package.json 不存在" -ForegroundColor Red
    exit 1
}

if (-Not (Test-Path ".env")) {
    Write-Host "⚠ .env 文件不存在，正在创建..." -ForegroundColor Yellow
    Copy-Item ".env.example" ".env" -Force
    Write-Host "✓ .env 已创建，请编辑并填入 GitHub OAuth 凭证" -ForegroundColor Green
}

# 安装依赖
Write-Host "`n[2/3] 安装依赖..." -ForegroundColor Yellow

if (-Not (Test-Path "node_modules")) {
    Write-Host "安装根目录依赖..." -ForegroundColor Gray
    npm install --silent
}

if (-Not (Test-Path "web/node_modules")) {
    Write-Host "安装前端依赖..." -ForegroundColor Gray
    Push-Location web
    npm install --silent
    Pop-Location
}

Write-Host "✓ 依赖检查完毕" -ForegroundColor Green

# 启动服务
Write-Host "`n[3/3] 启动服务..." -ForegroundColor Yellow

Write-Host "`n启动 Mock API 服务器 (端口 3000)..." -ForegroundColor Cyan
Start-Process pwsh -ArgumentList "-NoExit", "-Command", "cd '$PSScriptRoot'; node mock-server.js" -WindowStyle Normal

Write-Host "启动前端开发服务器 (端口 4321)..." -ForegroundColor Cyan
Start-Process pwsh -ArgumentList "-NoExit", "-Command", "cd '$PSScriptRoot/web'; npm run dev" -WindowStyle Normal

Write-Host @"

╔════════════════════════════════════════════════════════════════╗
║                  🎉 开发环境已启动！                            ║
╠════════════════════════════════════════════════════════════════╣
║                                                                ║
║  前端:  http://localhost:4321                                 ║
║  API:   http://localhost:3000                                 ║
║                                                                ║
║  📱 在浏览器中打开: http://localhost:4321                      ║
║  🔍 打开 F12 查看控制台日志                                    ║
║                                                                ║
║  📝 编辑配置:                                                 ║
║     1. .env - GitHub OAuth 凭证                              ║
║     2. mock-server.js - 模拟 API 响应                         ║
║                                                                ║
║  📖 更多信息: 查看 LOCAL_DEBUG.md                             ║
║                                                                ║
╚════════════════════════════════════════════════════════════════╝

" -ForegroundColor Green

Write-Host "两个服务窗口已打开，请保持开启" -ForegroundColor Yellow
Write-Host "按 Ctrl+C 在各窗口中停止服务" -ForegroundColor Gray
