package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"tmail/ent"
	"tmail/ent/user"
	"time"
)

// GitHubUser represents GitHub user response
type GitHubUser struct {
	ID      int    `json:"id"`
	Login   string `json:"login"`
	Avatar  string `json:"avatar_url"`
	Email   string `json:"email"`
	Name    string `json:"name"`
}

// AuthService handles GitHub OAuth authentication
type AuthService struct {
	client       *ent.Client
	githubID     string
	githubSecret string
	redirectURL  string
}

func NewAuthService(client *ent.Client) *AuthService {
	return &AuthService{
		client:       client,
		githubID:     os.Getenv("GITHUB_OAUTH_ID"),
		githubSecret: os.Getenv("GITHUB_OAUTH_SECRET"),
		redirectURL:  os.Getenv("GITHUB_OAUTH_REDIRECT"),
	}
}

// GetGitHubAuthURL returns the GitHub OAuth authorization URL
func (as *AuthService) GetGitHubAuthURL(state string) string {
	return fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email&state=%s",
		as.githubID,
		as.redirectURL,
		state,
	)
}

// ExchangeCode exchanges GitHub code for access token
func (as *AuthService) ExchangeCode(code string) (string, error) {
	req, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", nil)
	q := req.URL.Query()
	q.Add("client_id", as.githubID)
	q.Add("client_secret", as.githubSecret)
	q.Add("code", code)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	
	if accessToken, ok := result["access_token"].(string); ok {
		return accessToken, nil
	}
	return "", fmt.Errorf("failed to get access token")
}

// GetGitHubUser fetches GitHub user info with access token
func (as *AuthService) GetGitHubUser(accessToken string) (*GitHubUser, error) {
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var ghUser GitHubUser
	json.Unmarshal(body, &ghUser)
	return &ghUser, nil
}

// CreateOrUpdateUser creates or updates user from GitHub data
func (as *AuthService) CreateOrUpdateUser(ctx context.Context, ghUser *GitHubUser) (*ent.User, error) {
	u, err := as.client.User.Query().Where(user.GithubID(fmt.Sprintf("%d", ghUser.ID))).First(ctx)
	
	if err == nil {
		// User exists, update
		return as.client.User.UpdateOne(u).
			SetUsername(ghUser.Login).
			SetAvatarURL(ghUser.Avatar).
			SetEmail(ghUser.Email).
			SetUpdatedAt(time.Now()).
			Save(ctx)
	}
	
	// Create new user
	return as.client.User.Create().
		SetGithubID(fmt.Sprintf("%d", ghUser.ID)).
		SetUsername(ghUser.Login).
		SetAvatarURL(ghUser.Avatar).
		SetEmail(ghUser.Email).
		Save(ctx)
}

// GenerateAPIToken creates a new API token for user
func (as *AuthService) GenerateAPIToken(ctx context.Context, userID int, tokenName string) (string, error) {
	token := generateRandomToken(32)
	
	_, err := as.client.APIToken.Create().
		SetToken(token).
		SetName(tokenName).
		SetUserID(userID).
		Save(ctx)
	
	if err != nil {
		return "", err
	}
	return token, nil
}

// ValidateToken validates API token and returns user
func (as *AuthService) ValidateToken(ctx context.Context, token string) (*ent.User, error) {
	t, err := as.client.APIToken.Query().
		Where().
		WithUser().
		First(ctx)
	
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}
	
	if t.Revoked {
		return nil, fmt.Errorf("token revoked")
	}
	
	// Update last used time
	as.client.APIToken.UpdateOne(t).
		SetLastUsedAt(time.Now()).
		Exec(ctx)
	
	return t.Edges.User, nil
}

// Helper function to generate random token
func generateRandomToken(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	token := make([]byte, length)
	for i := range token {
		token[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(token)
}
