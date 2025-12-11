package api

import (
	"strings"
	"tmail/ent"
	"tmail/ent/apitoken"

	"github.com/labstack/echo/v4"
)

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

// BearerToken returns a middleware that checks for API token authorization
func BearerToken(client *ent.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			
			if authHeader == "" {
				return c.JSON(401, map[string]string{"error": "missing authorization header"})
			}
			
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(401, map[string]string{"error": "invalid authorization header"})
			}
			
			token := parts[1]
			
			// Validate token
			authService := NewAuthService(client)
			user, err := authService.ValidateToken(c.Request().Context(), token)
			
			if err != nil {
				return c.JSON(401, map[string]string{"error": err.Error()})
			}
			
			// Store user in context
			c.Set("user", user)
			
			return next(c)
		}
	}
}

// User returns the authenticated user from context
func GetUser(c echo.Context) *ent.User {
	if user, ok := c.Get("user").(*ent.User); ok {
		return user
	}
	return nil
}
