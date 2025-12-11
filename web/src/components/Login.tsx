import React, { useState } from "react"
import { Button } from "@/components/ui/button"

function Login() {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState("")

  const handleGitHubLogin = async () => {
    try {
      setLoading(true)
      const state = Math.random().toString(36).substring(7)

      // Get GitHub auth URL from backend
      const response = await fetch(`/api/auth/url?state=${state}`)
      const { auth_url } = await response.json()

      // Redirect to GitHub
      window.location.href = auth_url
    } catch (err) {
      setError("Failed to initiate login")
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-white dark:bg-slate-950 px-4">
      <div className="w-full max-w-md space-y-8">
        {/* Header */}
        <div className="text-center space-y-2">
          <img src="/favicon.svg" alt="Tmail" className="mx-auto h-12 w-12" />
          <h1 className="text-3xl font-bold">Tmail</h1>
          <p className="text-gray-600 dark:text-gray-400">
            Login to access your temporary mailboxes
          </p>
        </div>

        {/* Error Message */}
        {error && (
          <div className="rounded-lg bg-red-50 dark:bg-red-950 p-4 text-red-800 dark:text-red-200">
            {error}
          </div>
        )}

        {/* Login Buttons */}
        <div className="space-y-4">
          <Button
            onClick={handleGitHubLogin}
            disabled={loading}
            className="w-full bg-slate-900 hover:bg-slate-800 dark:bg-white dark:text-black dark:hover:bg-gray-100"
          >
            {loading ? "Logging in..." : "Login with GitHub"}
          </Button>
        </div>

        {/* Disclaimer */}
        <div className="rounded-lg bg-amber-50 dark:bg-amber-950 p-4 text-sm text-amber-900 dark:text-amber-100 space-y-2">
          <p>❗ Emails are only kept for 10 days</p>
          <p>❗ These mailboxes are temporary and public</p>
          <p>❗ Do not use for important accounts</p>
        </div>
      </div>
    </div>
  )
}

export default Login
