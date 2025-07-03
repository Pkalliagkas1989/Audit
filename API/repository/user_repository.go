package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"forum/models"
	"forum/utils"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailTaken         = errors.New("email is already taken")
	ErrUsernameTaken      = errors.New("username is already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrPasswordNotSet     = errors.New("password not set; login with oauth")
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
		if err == ErrUserNotFound {
			return nil, ErrPasswordNotSet
		}
		return nil, err
	}

	if auth.PasswordHash == "" {
		return nil, ErrPasswordNotSet
	}

	if !utils.CheckPasswordHash(login.Password, auth.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

// FindOrCreateOAuthUser looks up a user by provider credentials or email and creates
// a new account if none exists. Username uniqueness is ensured.
func (r *UserRepository) FindOrCreateOAuthUser(email, name, provider, providerID string) (*models.User, error) {
	// Try by provider
	user, err := r.GetByProvider(provider, providerID)
	if err == nil {
		return user, nil
	}

	// Try by email
	user, err = r.GetByEmail(email)
	if err == nil {
		// link provider
		if err := r.LinkProvider(user.ID, provider, providerID); err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		return user, nil
	}

	if err != ErrUserNotFound {
		return nil, err
	}

	// Ensure unique username derived from name or email
	base := strings.TrimSpace(name)
	if base == "" {
		base = strings.Split(email, "@")[0]
	}
	username := base
	i := 1
	for {
		var count int
		if err := r.DB.QueryRow("SELECT COUNT(*) FROM user WHERE username = ?", username).Scan(&count); err != nil {
			return nil, err
		}
		if count == 0 {
			break
		}
		username = fmt.Sprintf("%s%d", base, i)
		i++
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	userID := utils.GenerateUUID()
	createdAt := time.Now()
	if _, err := tx.Exec(`INSERT INTO user (user_id, username, email, created_at) VALUES (?, ?, ?, ?)`, userID, username, email, createdAt); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(`INSERT INTO user_providers (user_id, provider, provider_id) VALUES (?, ?, ?)`, userID, provider, providerID); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &models.User{ID: userID, Username: username, Email: email, CreatedAt: createdAt}, nil
}

// GetByProvider returns user by provider and id
func (r *UserRepository) GetByProvider(provider, providerID string) (*models.User, error) {
	query := `SELECT u.user_id, u.username, u.email, u.created_at FROM user u JOIN user_providers up ON u.user_id = up.user_id WHERE up.provider = ? AND up.provider_id = ?`
	var u models.User
	var created time.Time
	err := r.DB.QueryRow(query, provider, providerID).Scan(&u.ID, &u.Username, &u.Email, &created)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	u.CreatedAt = created
	return &u, nil
}

// LinkProvider associates an OAuth provider with an existing user
func (r *UserRepository) LinkProvider(userID, provider, providerID string) error {
	_, err := r.DB.Exec(`INSERT OR IGNORE INTO user_providers (user_id, provider, provider_id) VALUES (?, ?, ?)`, userID, provider, providerID)
	return err
}

// HasPassword checks if the user has a password set
func (r *UserRepository) HasPassword(userID string) (bool, error) {
	var hash sql.NullString
	err := r.DB.QueryRow(`SELECT password_hash FROM user_auth WHERE user_id = ?`, userID).Scan(&hash)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return hash.Valid && hash.String != "", nil
}
