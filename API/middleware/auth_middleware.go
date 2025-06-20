package middleware

import (
	"context"
	"net/http"

	"forum/models"
	"forum/repository"
)

// Authentication middleware checks if the user is authenticated
type AuthMiddleware struct {
	SessionRepo *repository.SessionRepository
	UserRepo    *repository.UserRepository
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(sessionRepo *repository.SessionRepository, userRepo *repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		SessionRepo: sessionRepo,
		UserRepo:    userRepo,
	}
}

// Authenticate middleware verifies authentication and sets user in context
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the session cookie
		cookie, err := r.Cookie("session_id")
		if err != nil {
			// No cookie, proceed as unauthenticated
			next.ServeHTTP(w, r)
			return
		}

		// Validate the session
		session, err := m.SessionRepo.GetBySessionID(cookie.Value)
		if err != nil {
			// Invalid or expired session, clear the cookie and continue as unauthenticated
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			})
			next.ServeHTTP(w, r)
			return
		}

		// Get the user
		user, err := m.UserRepo.GetByID(session.UserID)
		if err != nil {
			// User not found, clear the cookie and continue as unauthenticated
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			})
			next.ServeHTTP(w, r)
			return
		}

		// Set user in context
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth middleware ensures the user is authenticated
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user")
		if user == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// GetCurrentUser returns the authenticated user from the context
func GetCurrentUser(r *http.Request) *models.User {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		return nil
	}
	return user
}