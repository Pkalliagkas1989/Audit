package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"forum/models"
	"forum/repository"
	"forum/utils"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// OAuthHandler manages OAuth logins
type OAuthHandler struct {
	UserRepo     *repository.UserRepository
	SessionRepo  *repository.SessionRepository
	GoogleConfig *oauth2.Config
	GithubConfig *oauth2.Config
}

// NewOAuthHandler creates a new OAuthHandler
func NewOAuthHandler(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository) *OAuthHandler {
	googleConf := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:8080/forum/api/auth/google/callback",
		Scopes:       []string{"email", "profile"},
	}
	githubConf := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
		RedirectURL: "http://localhost:8080/forum/api/auth/github/callback",
		Scopes:      []string{"user:email"},
	}
	return &OAuthHandler{userRepo, sessionRepo, googleConf, githubConf}
}

func (h *OAuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := h.GoogleConfig.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *OAuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		utils.ErrorResponse(w, "missing code", http.StatusBadRequest)
		return
	}
	tok, err := h.GoogleConfig.Exchange(context.Background(), code)
	if err != nil {
		utils.ErrorResponse(w, "oauth exchange failed", http.StatusBadRequest)
		return
	}
	client := h.GoogleConfig.Client(context.Background(), tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		utils.ErrorResponse(w, "failed to fetch userinfo", http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()
	var info struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		utils.ErrorResponse(w, "invalid response", http.StatusBadRequest)
		return
	}
	if !info.EmailVerified {
		utils.ErrorResponse(w, "email not verified", http.StatusUnauthorized)
		return
	}
	user, err := h.UserRepo.FindOrCreateOAuthUser(info.Email, "google", info.Sub)
	if err != nil {
		utils.ErrorResponse(w, "failed to create user", http.StatusInternalServerError)
		return
	}
	h.finishOAuthLogin(w, r, user)
}

func (h *OAuthHandler) GithubLogin(w http.ResponseWriter, r *http.Request) {
	url := h.GithubConfig.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *OAuthHandler) GithubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		utils.ErrorResponse(w, "missing code", http.StatusBadRequest)
		return
	}
	tok, err := h.GithubConfig.Exchange(context.Background(), code)
	if err != nil {
		utils.ErrorResponse(w, "oauth exchange failed", http.StatusBadRequest)
		return
	}
	client := h.GithubConfig.Client(context.Background(), tok)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		utils.ErrorResponse(w, "failed to fetch user", http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()
	var ghUser struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ghUser); err != nil {
		utils.ErrorResponse(w, "invalid user data", http.StatusBadRequest)
		return
	}
	email := ghUser.Email
	if email == "" {
		resp2, err := client.Get("https://api.github.com/user/emails")
		if err != nil {
			utils.ErrorResponse(w, "failed to fetch emails", http.StatusBadRequest)
			return
		}
		defer resp2.Body.Close()
		var emails []struct {
			Email    string `json:"email"`
			Primary  bool   `json:"primary"`
			Verified bool   `json:"verified"`
		}
		if err := json.NewDecoder(resp2.Body).Decode(&emails); err != nil {
			utils.ErrorResponse(w, "invalid email data", http.StatusBadRequest)
			return
		}
		for _, e := range emails {
			if e.Primary && e.Verified {
				email = e.Email
				break
			}
		}
		if email == "" {
			for _, e := range emails {
				if e.Verified {
					email = e.Email
					break
				}
			}
		}
		if email == "" {
			utils.ErrorResponse(w, "no verified email", http.StatusBadRequest)
			return
		}
	}
	user, err := h.UserRepo.FindOrCreateOAuthUser(email, "github", fmt.Sprint(ghUser.ID))
	if err != nil {
		utils.ErrorResponse(w, "failed to create user", http.StatusInternalServerError)
		return
	}
	h.finishOAuthLogin(w, r, user)
}

func (h *OAuthHandler) finishOAuthLogin(w http.ResponseWriter, r *http.Request, user *models.User) {
	csrfToken := utils.GenerateCSRFToken()
	session, err := h.SessionRepo.Create(user.ID, r.RemoteAddr, csrfToken)
	if err != nil {
		utils.ErrorResponse(w, "failed to create session", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.SessionID,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	utils.JSONResponse(w, models.LoginResponse{User: *user, SessionID: session.SessionID, CSRFToken: csrfToken}, http.StatusOK)
}
