package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"forum/models"
	"forum/utils"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailTaken         = errors.New("email is already taken")
	ErrUsernameTaken      = errors.New("username is already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// UserRepository handles user-related database operations
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Create adds a new user to the database
func (r *UserRepository) Create(reg models.UserRegistration) (*models.User, error) {
	// Check if email is already taken
	var count int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM user WHERE email = ?", reg.Email).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrEmailTaken
	}

	// Check if username is already taken
	err = r.DB.QueryRow("SELECT COUNT(*) FROM user WHERE username = ?", reg.Username).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrUsernameTaken
	}

	// Start a transaction
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Generate UUID for the user
	userID := utils.GenerateUUID()
	createdAt := time.Now()
	// Insert user record
	_, err = tx.Exec(
		"INSERT INTO user (user_id, username, email, created_at) VALUES (?, ?, ?, ?)",
		userID, reg.Username, reg.Email, createdAt,
	)
	if err != nil {
		return nil, err
	}

	// Hash the password
	passwordHash, err := utils.HashPassword(reg.Password)
	if err != nil {
		return nil, err
	}

	// Insert authentication record
	_, err = tx.Exec(
		"INSERT INTO user_auth (user_id, password_hash) VALUES (?, ?)",
		userID, passwordHash,
	)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Return the newly created user
	user := &models.User{
		ID:        userID,
		Username:  reg.Username,
		Email:     reg.Email,
		CreatedAt: createdAt,
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	//var timestamp string
	var createdAt time.Time

	err := r.DB.QueryRow(
		"SELECT user_id, username, email, created_at FROM user WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &createdAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Parse the timestamp
	user.CreatedAt = createdAt
	//user.CreatedAt, err = time.Parse("2006-01-02 15:04:05.999999999-07:00", timestamp)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id string) (*models.User, error) {
	var user models.User
	//var timestamp string
	var createdAt time.Time
	err := r.DB.QueryRow(
		"SELECT user_id, username, email, created_at FROM user WHERE user_id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &createdAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	user.CreatedAt = createdAt
	// Parse the timestamp
	//user.CreatedAt, err = time.Parse("2006-01-02 15:04:05.999999999-07:00", timestamp)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetAuthByUserID retrieves user authentication data by user ID
func (r *UserRepository) GetAuthByUserID(userID string) (*models.UserAuth, error) {
	var auth models.UserAuth

	err := r.DB.QueryRow(
		"SELECT user_id, password_hash FROM user_auth WHERE user_id = ?",
		userID,
	).Scan(&auth.UserID, &auth.PasswordHash)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &auth, nil
}

// Authenticate validates a user's login credentials
func (r *UserRepository) Authenticate(login models.UserLogin) (*models.User, error) {
	// Get the user by email
	fmt.Println(login.Email)
	user, err := r.GetByEmail(login.Email)
	if err != nil {
		fmt.Println("User email not found")
		return nil, ErrInvalidCredentials
	}

	// Get the user's authentication data
	auth, err := r.GetAuthByUserID(user.ID)
	if err != nil {
		fmt.Println("Auth record not found")
		return nil, err
	}
	fmt.Println("Stored hash:", auth.PasswordHash)
	fmt.Println("Password match?", utils.CheckPasswordHash(login.Password, auth.PasswordHash))

	// Check the password
	if !utils.CheckPasswordHash(login.Password, auth.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
