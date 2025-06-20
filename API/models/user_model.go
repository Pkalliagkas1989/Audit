package models

import (
	"time"
)

// User represents a forum user
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// UserAuth contains user authentication information
type UserAuth struct {
	UserID       string `json:"-"`
	PasswordHash string `json:"-"`
}

// UserRegistration is used for registration requests
type UserRegistration struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// UserLogin is used for login requests
type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Session represents a user session
type Session struct {
	UserID     string    `json:"user_id"`
	SessionID  string    `json:"session_id"`
	CSRFToken  string    `json:"csrf_token"` // CSRF token for security
	IPAddress  string    `json:"ip_address"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}

// LoginResponse is the response after successful login
type LoginResponse struct {
	User      User   `json:"user"`
	SessionID string `json:"session_id"`
	CSRFToken string `json:"csrf_token"`
}







