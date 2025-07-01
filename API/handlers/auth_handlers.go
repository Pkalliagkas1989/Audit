package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"forum/models"
	"forum/repository"
	"forum/utils"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	UserRepo    *repository.UserRepository
	SessionRepo *repository.SessionRepository
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository) *AuthHandler {
	return &AuthHandler{
		UserRepo:    userRepo,
		SessionRepo: sessionRepo,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reg models.UserRegistration
	err := json.NewDecoder(r.Body).Decode(&reg)
	if err != nil {
		utils.ErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	reg.Username = strings.TrimSpace(reg.Username)
	reg.Email = strings.TrimSpace(strings.ToLower(reg.Email))
	reg.Password = strings.TrimSpace(reg.Password)

	if reg.Username == "" || reg.Email == "" || reg.Password == "" {
		utils.ErrorResponse(w, "Username, email, and password are required", http.StatusBadRequest)
		return
	}

	if !utils.UsernameRegex.MatchString(reg.Username) {
		utils.ErrorResponse(w, "Username must be 3-50 characters, letters/numbers/underscores only", http.StatusBadRequest)
		return
	}

	cleanEmail, err := utils.ValidateEmail(reg.Email)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	reg.Email = cleanEmail

	if !utils.IsStrongPassword(reg.Password) {
		utils.ErrorResponse(w, "Password must contain letters and numbers", http.StatusBadRequest)
		return
	}

	user, err := h.UserRepo.Create(reg)
	if err != nil {
		switch err {
		case repository.ErrEmailTaken:
			utils.ErrorResponse(w, "Email is already taken", http.StatusConflict)
		case repository.ErrUsernameTaken:
			utils.ErrorResponse(w, "Username is already taken", http.StatusConflict)
		default:
			utils.ErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// 🟢 Create a session right after registration
	csrfToken := utils.GenerateCSRFToken()
	session, err := h.SessionRepo.Create(user.ID, r.RemoteAddr, csrfToken)
	if err != nil {
		utils.ErrorResponse(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// 🟢 Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.SessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Change to true in production (with HTTPS)
		SameSite: http.SameSiteLaxMode,
	})

	utils.JSONResponse(w, map[string]string{
		"message": "Registration successful",
	}, http.StatusCreated)
}


// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var login models.UserLogin
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if login.Email == "" || login.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Authenticate user
	user, err := h.UserRepo.Authenticate(login)

	if err != nil {
		if err == repository.ErrInvalidCredentials {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	csrfToken := utils.GenerateCSRFToken()	

	// Create a new session
	log.Println("Creating session for user:", user.ID)
	session, err := h.SessionRepo.Create(user.ID, r.RemoteAddr, csrfToken)
	log.Println("Session created:", session)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
    Name:     "session_id",
    Value:    session.SessionID,
    Path:     "/",
    Expires:  session.ExpiresAt,
    HttpOnly: true,
    Secure:   false,
    SameSite: http.SameSiteLaxMode, // Change from None to Lax for localhost
})


	// Return response
	w.Header().Set("Content-Type", "application/json")
	response := models.LoginResponse{
		User:      *user,
		SessionID: session.SessionID,
		CSRFToken: csrfToken,
	}
	json.NewEncoder(w).Encode(response)
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the session cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		// If no cookie, nothing to do
		w.WriteHeader(http.StatusOK)
		return
	}

	// Delete the session
	err = h.SessionRepo.Delete(cookie.Value)
	if err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	// Clear the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
}

// Inside AuthHandler
func (h *AuthHandler) VerifySession(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	session, err := h.SessionRepo.GetBySessionID(sessionCookie.Value)
	if err != nil {
		http.Error(w, "Session invalid or expired", http.StatusUnauthorized)
		return
	}

	// Optionally fetch user and return profile
	user, err := h.UserRepo.GetByID(session.UserID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, user, http.StatusOK)
}
