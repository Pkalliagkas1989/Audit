package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"forum/models"
	"forum/repository"
	"forum/utils"
)

// GoogleOAuthHandler manages Google OAuth authentication.
type GoogleOAuthHandler struct {
	Config      *oauth2.Config
	UserRepo    *repository.UserRepository
	SessionRepo *repository.SessionRepository
}

// NewGoogleOAuthHandler creates a GoogleOAuthHandler with config from environment variables.
func NewGoogleOAuthHandler(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository) *GoogleOAuthHandler {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_CALLBACK_URL")
	if redirectURL == "" {
		log.Fatal("GOOGLE_CALLBACK_URL environment variable is not set")
	}
	cfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return &GoogleOAuthHandler{Config: cfg, UserRepo: userRepo, SessionRepo: sessionRepo}
}

// Login initiates the Google OAuth flow.
func (h *GoogleOAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	url := h.Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Callback handles the Google OAuth callback.
func (h *GoogleOAuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	code := r.URL.Query().Get("code")
	if code == "" {
		utils.ErrorResponse(w, "missing code", http.StatusBadRequest)
		return
	}

	token, err := h.Config.Exchange(ctx, code)
	if err != nil {
		utils.ErrorResponse(w, "token exchange failed", http.StatusInternalServerError)
		return
	}

	client := h.Config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		utils.ErrorResponse(w, "failed to fetch user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var data struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		utils.ErrorResponse(w, "invalid user info", http.StatusInternalServerError)
		return
	}

	email := strings.ToLower(data.Email)
	user, err := h.UserRepo.GetByEmail(email)
	if err != nil {
		if err == repository.ErrUserNotFound {
			username := strings.Split(email, "@")[0]
			if !utils.UsernameRegex.MatchString(username) {
				username = "user_" + utils.GenerateUUID()[:8]
			}
			reg := models.UserRegistration{
				Username: username,
				Email:    email,
				Password: utils.GenerateUUID(),
			}
			user, err = h.UserRepo.Create(reg)
			if err != nil {
				utils.ErrorResponse(w, "failed to create user", http.StatusInternalServerError)
				return
			}
		} else {
			utils.ErrorResponse(w, "failed to fetch user", http.StatusInternalServerError)
			return
		}
	}

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
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "http://localhost:8081/", http.StatusSeeOther)
}
