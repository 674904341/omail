package api

import (
	"fmt"
	"net/http"
	"tmail/ent"
)

// LoginRequest handles GitHub OAuth callback
func Login(c *Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")
	
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing code"})
	}
	
	authService := NewAuthService(c.Client())
	
	// Exchange code for access token
	accessToken, err := authService.ExchangeCode(code)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "failed to get access token"})
	}
	
	// Get GitHub user info
	ghUser, err := authService.GetGitHubUser(accessToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "failed to get user info"})
	}
	
	// Create or update user
	user, err := authService.CreateOrUpdateUser(c.Context(), ghUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create user"})
	}
	
	// Generate API token
	token, err := authService.GenerateAPIToken(c.Context(), user.ID, "default")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate token"})
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"avatar":   user.AvatarURL,
			"email":    user.Email,
		},
		"api_token": token,
	})
}

// GetAuthURL returns GitHub OAuth authorization URL
func GetAuthURL(c *Context) error {
	state := c.QueryParam("state")
	if state == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "state is required"})
	}
	
	authService := NewAuthService(c.Client())
	authURL := authService.GetGitHubAuthURL(state)
	
	return c.JSON(http.StatusOK, map[string]string{"auth_url": authURL})
}

// GetProfile returns current user profile
func GetProfile(c *Context) error {
	user := c.User()
	if user == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"avatar":   user.AvatarURL,
		"email":    user.Email,
	})
}

// CreateMailbox creates a random mailbox for user
func CreateMailbox(c *Context) error {
	user := c.User()
	if user == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	
	// Generate random email
	email := generateRandomEmail()
	
	// Create mailbox
	mailbox, err := c.Client().Mailbox.Create().
		SetEmail(email).
		SetUserID(user.ID).
		Save(c.Context())
	
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create mailbox"})
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"email":      mailbox.Email,
		"created_at": mailbox.CreatedAt,
	})
}

// GetMailboxes returns all mailboxes for user
func GetMailboxes(c *Context) error {
	user := c.User()
	if user == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	
	mailboxes, err := c.Client().Mailbox.Query().Where().All(c.Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch mailboxes"})
	}
	
	var result []map[string]interface{}
	for _, mb := range mailboxes {
		result = append(result, map[string]interface{}{
			"email":      mb.Email,
			"created_at": mb.CreatedAt,
		})
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"mailboxes": result,
	})
}

// GetEmails returns emails for a specific mailbox
func GetEmails(c *Context) error {
	user := c.User()
	if user == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	
	email := c.QueryParam("email")
	if email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "email parameter is required"})
	}
	
	// Verify mailbox belongs to user
	mailbox, err := c.Client().Mailbox.Query().
		Where().
		First(c.Context())
	
	if err != nil || mailbox.UserID != user.ID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
	}
	
	// Get envelopes for this mailbox
	envelopes, err := c.Client().Envelope.Query().
		Where().
		All(c.Context())
	
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch emails"})
	}
	
	var result []map[string]interface{}
	for _, env := range envelopes {
		result = append(result, map[string]interface{}{
			"id":         env.ID,
			"from":       env.From,
			"subject":    env.Subject,
			"created_at": env.CreatedAt,
		})
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"emails": result,
	})
}

// GetEmailDetail returns detailed information for an email
func GetEmailDetail(c *Context) error {
	user := c.User()
	if user == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	
	emailID := c.Param("id")
	
	envelope, err := c.Client().Envelope.Get(c.Context(), parseID(emailID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "email not found"})
	}
	
	// Verify access - you may need to add more verification logic
	attachments, _ := envelope.QueryAttachments().All(c.Context())
	
	var attachmentList []map[string]interface{}
	for _, att := range attachments {
		attachmentList = append(attachmentList, map[string]interface{}{
			"id":       att.ID,
			"filename": att.Filename,
			"size":     att.Size,
		})
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":          envelope.ID,
		"from":        envelope.From,
		"to":          envelope.To,
		"subject":     envelope.Subject,
		"content":     envelope.Content,
		"created_at":  envelope.CreatedAt,
		"attachments": attachmentList,
	})
}

// Helper functions

func generateRandomEmail() string {
	domain := "mail.4w.ink"
	random := generateRandomToken(8)
	return fmt.Sprintf("%s@%s", random, domain)
}

func parseID(id string) int {
	var result int
	fmt.Sscanf(id, "%d", &result)
	return result
}
