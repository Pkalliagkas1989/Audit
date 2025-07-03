package handlers

import (
	"net/http"

	"forum/models"
	"forum/oauth"
	"forum/repository"
	"forum/utils"
)

// OAuthHandler manages OAuth login flows
type OAuthHandler struct {
	UserRepo    *repository.UserRepository
	SessionRepo *repository.SessionRepository
}

func NewOAuthHandler(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository) *OAuthHandler {
	return &OAuthHandler{UserRepo: userRepo, SessionRepo: sessionRepo}
}

func (h *OAuthHandler) login(w http.ResponseWriter, r *http.Request, p oauth.Provider) {
	state := utils.GenerateCSRFToken()
	http.SetCookie(w, &http.Cookie{Name: "oauth_state", Value: state, Path: "/", HttpOnly: true, MaxAge: 300})
	http.Redirect(w, r, p.AuthCodeURL(state), http.StatusFound)
}

func (h *OAuthHandler) GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	h.login(w, r, oauth.GoogleProvider{})
}

func (h *OAuthHandler) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	h.handleCallback(w, r, oauth.GoogleProvider{})
}

func (h *OAuthHandler) GitHubLoginHandler(w http.ResponseWriter, r *http.Request) {
	h.login(w, r, oauth.GitHubProvider{})
}

func (h *OAuthHandler) GitHubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	h.handleCallback(w, r, oauth.GitHubProvider{})
}

func (h *OAuthHandler) handleCallback(w http.ResponseWriter, r *http.Request, p oauth.Provider) {
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || r.URL.Query().Get("state") != stateCookie.Value {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "code not found", http.StatusBadRequest)
		return
	}
	token, err := p.Exchange(code)
	if err != nil {
		http.Error(w, "token exchange failed", http.StatusInternalServerError)
		return
	}
	email, pid, name, err := p.User(token)
	if err != nil {
		http.Error(w, "failed to fetch user", http.StatusInternalServerError)
		return
	}
	user, err := h.UserRepo.FindOrCreateOAuthUser(email, name, p.Name(), pid)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}
	csrfToken := utils.GenerateCSRFToken()
	session, err := h.SessionRepo.Create(user.ID, r.RemoteAddr, csrfToken)
	if err != nil {
		http.Error(w, "failed to create session", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "session_id", Value: session.SessionID, Path: "/", HttpOnly: true, Secure: false, SameSite: http.SameSiteLaxMode})
	utils.JSONResponse(w, models.LoginResponse{User: *user, SessionID: session.SessionID, CSRFToken: csrfToken}, http.StatusOK)
}
