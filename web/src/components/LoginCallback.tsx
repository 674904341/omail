import React, { useEffect, useState } from "react"

function LoginCallback() {
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")

  useEffect(() => {
    const handleCallback = async () => {
      try {
        const params = new URLSearchParams(window.location.search)
        const code = params.get("code")
        const state = params.get("state")

        if (!code) {
          setError("No authorization code provided")
          setLoading(false)
          return
        }

        // Exchange code for API token
        const response = await fetch(`/api/auth/login?code=${code}&state=${state}`, {
          method: "POST",
        })

        if (!response.ok) {
          throw new Error("Login failed")
        }

        const { user, api_token } = await response.json()

        // Save token to localStorage
        localStorage.setItem("api_token", api_token)
        localStorage.setItem("user", JSON.stringify(user))

        // Redirect to dashboard
        window.location.href = "/dashboard"
      } catch (err) {
        setError("Login failed: " + (err instanceof Error ? err.message : "Unknown error"))
        setLoading(false)
      }
    }

    handleCallback()
  }, [])

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-white dark:bg-slate-950">
      {loading && (
        <div className="text-center space-y-4">
          <div className="inline-block h-8 w-8 animate-spin rounded-full border-4 border-slate-300 border-r-slate-900 dark:border-slate-600 dark:border-r-white"></div>
          <p className="text-gray-600 dark:text-gray-400">Logging you in...</p>
        </div>
      )}

      {error && (
        <div className="text-center space-y-4">
          <p className="text-red-600 dark:text-red-400">{error}</p>
          <a href="/" className="text-blue-600 dark:text-blue-400 hover:underline">
            Back to login
          </a>
        </div>
      )}
    </div>
  )
}

export default LoginCallback
