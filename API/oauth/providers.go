package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Provider defines methods required to implement an OAuth provider.
type Provider interface {
	// Name returns the provider identifier (e.g. "google").
	Name() string
	// AuthCodeURL returns the authorization URL for the given state.
	AuthCodeURL(state string) string
	// Exchange exchanges an authorization code for an access token.
	Exchange(code string) (string, error)
	// User retrieves the user's email, provider user ID and display name.
	User(token string) (email, id, name string, err error)
}

// ---------------- Google provider ----------------

type GoogleProvider struct{}

type googleToken struct {
	AccessToken string `json:"access_token"`
}

type googleUser struct {
	ID            string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
}

func (GoogleProvider) Name() string { return "google" }

func (GoogleProvider) AuthCodeURL(state string) string {
	q := url.Values{}
	q.Set("client_id", os.Getenv("GOOGLE_CLIENT_ID"))
	q.Set("redirect_uri", os.Getenv("GOOGLE_REDIRECT_URL"))
	q.Set("response_type", "code")
	q.Set("scope", "openid email profile")
	q.Set("state", state)
	return "https://accounts.google.com/o/oauth2/v2/auth?" + q.Encode()
}

func (GoogleProvider) Exchange(code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", os.Getenv("GOOGLE_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("GOOGLE_CLIENT_SECRET"))
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", os.Getenv("GOOGLE_REDIRECT_URL"))
	resp, err := http.PostForm("https://oauth2.googleapis.com/token", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var t googleToken
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return "", err
	}
	if t.AccessToken == "" {
		return "", errors.New("no access token")
	}
	return t.AccessToken, nil
}

func (GoogleProvider) User(token string) (string, string, string, error) {
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()
	var u googleUser
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return "", "", "", err
	}
	if !u.EmailVerified {
		return "", "", "", errors.New("email not verified")
	}
	return strings.ToLower(u.Email), u.ID, u.Name, nil
}

// ---------------- GitHub provider ----------------

type GitHubProvider struct{}

type githubToken struct {
	AccessToken string `json:"access_token"`
}

type githubUser struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (GitHubProvider) Name() string { return "github" }

func (GitHubProvider) AuthCodeURL(state string) string {
	q := url.Values{}
	q.Set("client_id", os.Getenv("GITHUB_CLIENT_ID"))
	q.Set("redirect_uri", os.Getenv("GITHUB_REDIRECT_URL"))
	q.Set("scope", "user:email")
	q.Set("state", state)
	return "https://github.com/login/oauth/authorize?" + q.Encode()
}

func (GitHubProvider) Exchange(code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", os.Getenv("GITHUB_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("GITHUB_CLIENT_SECRET"))
	data.Set("code", code)
	data.Set("redirect_uri", os.Getenv("GITHUB_REDIRECT_URL"))

	req, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var t githubToken
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return "", err
	}
	if t.AccessToken == "" {
		return "", errors.New("no access token")
	}
	return t.AccessToken, nil
}

func (GitHubProvider) User(token string) (string, string, string, error) {
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()
	var u githubUser
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return "", "", "", err
	}
	email := strings.ToLower(u.Email)
	if email == "" {
		req2, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
		req2.Header.Set("Authorization", "Bearer "+token)
		req2.Header.Set("Accept", "application/vnd.github+json")
		resp2, err := http.DefaultClient.Do(req2)
		if err == nil {
			defer resp2.Body.Close()
			var emails []struct {
				Email    string `json:"email"`
				Primary  bool   `json:"primary"`
				Verified bool   `json:"verified"`
			}
			if err := json.NewDecoder(resp2.Body).Decode(&emails); err == nil {
				for _, e := range emails {
					if e.Primary && e.Verified {
						email = strings.ToLower(e.Email)
						break
					}
				}
			}
		}
	}
	if email == "" {
		return "", "", "", errors.New("no verified email")
	}
	return email, fmt.Sprintf("%d", u.ID), u.Login, nil
}
