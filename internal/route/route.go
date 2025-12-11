package route

import (
	"strings"
	"github.com/labstack/echo/v4"
	"tmail/ent"
	"tmail/ent/apitoken"
	"tmail/internal/api"
)

func Register(e *echo.Echo) {
	// Public auth routes
	authG := e.Group("/api/auth")
	authG.GET("/url", api.Wrap(api.GetAuthURL))
	authG.POST("/login", api.Wrap(api.Login))
	
	// Protected API routes (require API token)
	apiG := e.Group("/api")
	// Note: APITokenAuthMiddleware should be registered in app middleware
	apiG.GET("/profile", api.Wrap(api.GetProfile))
	apiG.POST("/mailbox", api.Wrap(api.CreateMailbox))
	apiG.GET("/mailboxes", api.Wrap(api.GetMailboxes))
	apiG.GET("/emails", api.Wrap(api.GetEmails))
	apiG.GET("/email/:id", api.Wrap(api.GetEmailDetail))
	
	// Original routes (keep backward compatibility)
	apiG.POST("/report", api.Wrap(api.Report))
	apiG.GET("/fetch", api.Wrap(api.Fetch))
	apiG.GET("/fetch/:id", api.Wrap(api.FetchDetail))
	apiG.GET("/fetch/latest", api.Wrap(api.FetchLatest))
	apiG.GET("/download/:id", api.Wrap(api.Download))
	apiG.GET("/domain", api.Wrap(api.DomainList))
}

// APITokenAuthMiddleware validates API token from Authorization header
func APITokenAuthMiddleware(client *ent.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			
			// Skip if no auth header
			if authHeader == "" {
				return next(c)
			}
			
			// Extract token from "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return next(c)
			}
			
			token := parts[1]
			
			// Query token from database
			apiToken, err := client.APIToken.Query().
				Where(apitoken.Token(token)).
				WithUser().
				First(c.Request().Context())
			
			if err != nil || apiToken == nil {
				// Invalid token, continue without user context
				return next(c)
			}
			
			if apiToken.Revoked {
				// Token is revoked
				return next(c)
			}
			
			// Store user in context
			c.Set("user", apiToken.Edges.User)
			c.Set("api_token", apiToken)
			
			return next(c)
		}
	}
}
