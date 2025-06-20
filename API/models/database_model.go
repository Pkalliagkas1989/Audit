package models

import (
	"database/sql"
	"fmt"
	"forum/config"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// initDB initializes the database and returns a connection
func InitDB() (*sql.DB, error) {
	dbPath := filepath.Join("./database", "forum.db")

	// Check if database file exists
	firstTime := false
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		firstTime = true
		// Create directory if it doesn't exist
		if err := os.MkdirAll("./database", 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %v", err)
		}
	}

	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Verify connection works
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Set some basic connection pool settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	// Initialize database schema and data
	if firstTime {
		if err := createTables(db); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to create tables: %v", err)
		}

		if err := createIndexes(db); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to create indexes: %v", err)
		}

		if err := populateCategories(db, config.Categories); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to populate categories: %v", err)
		}
		// if err := InsertMockData(db); err != nil {
		// 	db.Close()
		// 	return nil, fmt.Errorf("failed to insert mock data: %v", err)
		// }
		fmt.Println("Database initialized successfully.")
	} else {
		fmt.Println("Database already exists. Skipping initialization.")
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	// Start a transaction for atomicity
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Define all table creation SQL statements
	tableStatements := []string{
		config.CreateUserTable,
		config.CreateUserAuthTable,
		config.CreateSessionsTable,
		config.CreateCategoriesTable,
		config.CreatePostsTable,
		config.CreateCommentsTable,
		config.CreateReactionsTable,
		config.CreatePostCategoriesTable,
	}

	// Execute each table creation statement
	for i, stmt := range tableStatements {
		_, err = tx.Exec(stmt)
		if err != nil {
			return fmt.Errorf("statement %d failed: %v\nSQL: %s", i+1, err, stmt)
		}
	}

	// Commit transaction
	return tx.Commit()
}

func createIndexes(db *sql.DB) error {
	// Start a transaction for atomicity
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Define all index creation SQL statements
	indexStatements := []string{
		config.IdxPostsUserID,
		config.IdxPostCategoriesPostID,
		config.IdxPostCategoriesCategoryID,
		config.IdxCommentsPostID,
		config.IdxCommentsUserID,
		config.IdxReactionsUserID,
		config.IdxReactionsPostID,
		config.IdxReactionsCommentID,
	}

	// Execute each index creation statement
	for _, stmt := range indexStatements {
		_, err = tx.Exec(stmt)
		if err != nil {
			return fmt.Errorf("failed to execute create index statement: %s: %v", stmt, err)
		}
	}

	// Commit transaction
	return tx.Commit()
}

func populateCategories(db *sql.DB, categories []string) error {
	if len(categories) == 0 {
		return nil
	}

	// Use transaction for atomicity
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT OR IGNORE INTO categories (name) VALUES (?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	for _, category := range categories {
		if _, err := stmt.Exec(category); err != nil {
			return fmt.Errorf("failed to insert category '%s': %v", category, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	fmt.Println("Categories populated (duplicates ignored if existed).")
	return nil
}

// // InsertMockData inserts mock users, categories, posts, comments, and reactions for testing/demo purposes.
// func InsertMockData(db *sql.DB) error {
// 	now := "2025-05-26T14:30:00"

// 	users := map[string]string{
// 		"Alice": "alice@example.com",
// 		"Bob":   "bob@example.com",
// 	}

// 	userIDs := map[string]string{
// 		"Alice": "e6f50c45-77f7-45d6-a206-111111111111",
// 		"Bob":   "f7e61d55-88e8-46e6-b307-222222222222",
// 	}

// 	categories := map[string]int{
// 		"General":              1,
// 		"Software Development": 2,
// 		"Hobbies":              3,
// 		"Random":               4,
// 		"Pets":                 5,
// 		"Travel":               6,
// 		"EMVALOTIS":            7,
// 	}

// 	posts := map[string]string{
// 		"AlicePost": "a1e2b3c4-d5f6-7890-1234-aaaaaaaaaaaa",
// 		"BobPost":   "b1e2c3d4-f5g6-7890-1234-bbbbbbbbbbbb",
// 	}

// 	comments := map[string]string{
// 		"AliceOnBob": "c1d2e3f4-g5h6-7890-1234-cccccccccccc",
// 		"BobOnAlice": "d1e2f3g4-h5i6-7890-1234-dddddddddddd",
// 	}

// 	tx, err := db.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()

// 	// Insert users
// 	for name, email := range users {
// 		_, err := tx.Exec(`INSERT OR IGNORE INTO user (user_id, username, email, created_at) VALUES (?, ?, ?, ?)`,
// 			userIDs[name], name, email, now)
// 		if err != nil {
// 			return fmt.Errorf("insert user %s: %v", name, err)
// 		}
// 	}

// 	// Insert user passwords
// 	for name := range users {
// 		_, err := tx.Exec(`INSERT OR IGNORE INTO user_auth (user_id, password_hash) VALUES (?, ?)`,
// 			userIDs[name], "hashed_password_123")
// 		if err != nil {
// 			return fmt.Errorf("insert password for user %s: %v", name, err)
// 		}
// 	}

// 	// Insert categories
// 	for name, id := range categories {
// 		_, err := tx.Exec(`INSERT OR IGNORE INTO categories (category_id, name) VALUES (?, ?)`, id, name)
// 		if err != nil {
// 			return fmt.Errorf("insert category %s: %v", name, err)
// 		}
// 	}

// 	// Insert posts
// 	_, err = tx.Exec(`INSERT INTO posts (post_id, user_id, category_id, title, content, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
// 		posts["AlicePost"], userIDs["Alice"], categories["Software Development"], "Alice's Title", "Alice on tech.", now)
// 	if err != nil {
// 		return fmt.Errorf("insert Alice post: %v", err)
// 	}

// 	_, err = tx.Exec(`INSERT INTO posts (post_id, user_id, category_id, title, content, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
// 		posts["BobPost"], userIDs["Bob"], categories["Pets"], "Bob's Title", "Bob on lifestyle.", now)
// 	if err != nil {
// 		return fmt.Errorf("insert Bob post: %v", err)
// 	}

// 	// Insert comments
// 	_, err = tx.Exec(`INSERT INTO comments (comment_id, post_id, user_id, content, created_at) VALUES (?, ?, ?, ?, ?)`,
// 		comments["AliceOnBob"], posts["BobPost"], userIDs["Alice"], "Nice post, Bob!", now)
// 	if err != nil {
// 		return fmt.Errorf("insert Alice comment: %v", err)
// 	}

// 	_, err = tx.Exec(`INSERT INTO comments (comment_id, post_id, user_id, content, created_at) VALUES (?, ?, ?, ?, ?)`,
// 		comments["BobOnAlice"], posts["AlicePost"], userIDs["Bob"], "Interesting point, Alice.", now)
// 	if err != nil {
// 		return fmt.Errorf("insert Bob comment: %v", err)
// 	}

// 	// Insert reactions to comments
// 	_, err = tx.Exec(`INSERT INTO reactions (user_id, comment_id, post_id, reaction_type, created_at) VALUES (?, ?, ?, ?, ?)`,
// 		userIDs["Alice"], comments["BobOnAlice"], nil, 1, now)
// 	if err != nil {
// 		return fmt.Errorf("insert Alice reaction: %v", err)
// 	}

// 	_, err = tx.Exec(`INSERT INTO reactions (user_id, comment_id, post_id, reaction_type, created_at) VALUES (?, ?, ?, ?, ?)`,
// 		userIDs["Bob"], comments["AliceOnBob"], nil, 2, now)
// 	if err != nil {
// 		return fmt.Errorf("insert Bob reaction: %v", err)
// 	}

// 	// Additional reactions (likes on posts)
// 	_, err = tx.Exec(`INSERT INTO reactions (user_id, comment_id, post_id, reaction_type, created_at) VALUES (?, ?, ?, ?, ?)`,
// 		userIDs["Alice"], nil, posts["BobPost"], 1, now)
// 	if err != nil {
// 		return fmt.Errorf("insert Alice reaction on Bob's post: %v", err)
// 	}

// 	_, err = tx.Exec(`INSERT INTO reactions (user_id, comment_id, post_id, reaction_type, created_at) VALUES (?, ?, ?, ?, ?)`,
// 		userIDs["Bob"], nil, posts["AlicePost"], 2, now)
// 	if err != nil {
// 		return fmt.Errorf("insert Bob reaction on Alice's post: %v", err)
// 	}

// 	return tx.Commit()
// }
