package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cthulhu-platform/auth/internal/pkg"
	"golang.org/x/oauth2"
)

// Helper functions
func generateCodeVerifier() string {
	b := make([]byte, 32)
	_, _ = io.ReadFull(rand.Reader, b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func generateCodeChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

func buildOAuthURL(provider, state, codeChallenge string) string {
	baseURL := "https://github.com/login/oauth/authorize"
	clientID := pkg.GITHUB_CLIENT_ID
	redirectURI := getRedirectURI(provider)

	return fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&state=%s&code_challenge=%s&code_challenge_method=S256&scope=read:user user:email",
		baseURL, clientID, redirectURI, state, codeChallenge)
}

func getRedirectURI(provider string) string {
	if provider == "github" {
		return pkg.GITHUB_REDIRECT_URI
	}
	return ""
}

func getOAuthConfig(provider string) *oauth2.Config {
	if provider == "github" {
		return &oauth2.Config{
			ClientID:     pkg.GITHUB_CLIENT_ID,
			ClientSecret: pkg.GITHUB_CLIENT_SECRET,
			RedirectURL:  pkg.GITHUB_REDIRECT_URI,
			Scopes:       []string{"read:user", "user:email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://github.com/login/oauth/authorize",
				TokenURL: "https://github.com/login/oauth/access_token",
			},
		}
	}
	return nil
}

type githubUserInfo struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
}

func fetchUserInfo(provider, accessToken string) (*struct {
	OAuthUserID string
	Email       string
	Username    *string
	AvatarURL   *string
}, error) {
	if provider == "github" {
		req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Accept", "application/vnd.github.v3+json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to fetch user info: status %d", resp.StatusCode)
		}

		var githubUser githubUserInfo
		if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
			return nil, err
		}

		// Get email if not in profile
		if githubUser.Email == "" {
			email, err := fetchGitHubEmail(accessToken)
			if err == nil && email != "" {
				githubUser.Email = email
			}
		}

		username := githubUser.Login
		avatarURL := githubUser.AvatarURL

		return &struct {
			OAuthUserID string
			Email       string
			Username    *string
			AvatarURL   *string
		}{
			OAuthUserID: fmt.Sprintf("%d", githubUser.ID),
			Email:       githubUser.Email,
			Username:    &username,
			AvatarURL:   &avatarURL,
		}, nil
	}

	return nil, fmt.Errorf("unsupported provider: %s", provider)
}

func fetchGitHubEmail(accessToken string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch emails: status %d", resp.StatusCode)
	}

	var emails []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	for _, email := range emails {
		if email.Primary {
			return email.Email, nil
		}
	}

	if len(emails) > 0 {
		return emails[0].Email, nil
	}

	return "", nil
}
