import React, { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { LogOut, LogIn } from "lucide-react"

interface AuthButtonProps {
  lang: string
}

export default function AuthButton({ lang }: AuthButtonProps) {
  const [isLoggedIn, setIsLoggedIn] = useState(false)
  const [username, setUsername] = useState("")
  const [avatar, setAvatar] = useState("")
  const [loading, setLoading] = useState(false)

  // 检查登录状态的函数
  const checkLoginStatus = () => {
    const token = localStorage.getItem("api_token")
    const user = localStorage.getItem("user")

    if (token && user) {
      setIsLoggedIn(true)
      try {
        const userData = JSON.parse(user)
        setUsername(userData.username || "User")
        setAvatar(userData.avatar || "")
      } catch {
        setUsername("User")
        setAvatar("")
      }
    } else {
      setIsLoggedIn(false)
      setUsername("")
      setAvatar("")
    }
  }

  useEffect(() => {
    // 初始化时检查登录状态
    checkLoginStatus()

    // 监听 storage 变化（其他标签页登录时）
    const handleStorageChange = () => {
      checkLoginStatus()
    }

    window.addEventListener("storage", handleStorageChange)
    return () => window.removeEventListener("storage", handleStorageChange)
  }, [])

  const handleLogin = async () => {
    try {
      setLoading(true)
      const state = Math.random().toString(36).substring(7)

      // Get GitHub auth URL from backend
      const response = await fetch(`/api/auth/url?state=${state}`)
      const { auth_url } = await response.json()

      // Redirect to GitHub
      window.location.href = auth_url
    } catch (err) {
      console.error("Failed to initiate login:", err)
      setLoading(false)
    }
  }

  const handleLogout = () => {
    localStorage.removeItem("api_token")
    localStorage.removeItem("user")
    setIsLoggedIn(false)
    setUsername("")
    window.location.reload()
  }

  if (isLoggedIn) {
    return (
      <div className="flex items-center gap-2">
        {avatar && (
          <img
            src={avatar}
            alt={username}
            className="w-8 h-8 rounded-full border border-gray-300 dark:border-gray-600"
          />
        )}
        <span className="text-sm text-gray-600 dark:text-gray-400">
          {username}
        </span>
        <Button
          onClick={handleLogout}
          variant="outline"
          size="sm"
          className="gap-1"
        >
          <LogOut className="h-4 w-4" />
          <span className="hidden sm:inline">Logout</span>
        </Button>
      </div>
    )
  }

  return (
    <Button
      onClick={handleLogin}
      disabled={loading}
      size="sm"
      className="gap-1"
    >
      <LogIn className="h-4 w-4" />
      <span className="hidden sm:inline">
        {loading ? "Loading..." : "登录"}
      </span>
    </Button>
  )
}
