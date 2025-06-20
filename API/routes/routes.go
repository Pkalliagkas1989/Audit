package routes

import (
	"database/sql"
	"net/http"

	"forum/handlers"
	"forum/middleware"
	"forum/repository"
)

func SetupRoutes(db *sql.DB) http.Handler {
	// Create repositories
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	categoryRepo := repository.NewCategoryRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	reactionRepo := repository.NewReactionRepository(db)

	// Create handlers
	authHandler := handlers.NewAuthHandler(userRepo, sessionRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	postHandler := handlers.NewPostHandler(postRepo)
	myPostsHandler := handlers.NewMyPostsHandler(postRepo, commentRepo, reactionRepo)
	likedPostsHandler := handlers.NewLikedPostsHandler(postRepo, commentRepo, reactionRepo)
	commentHandler := handlers.NewCommentHandler(commentRepo)
	reactionHandler := handlers.NewReactionHandler(reactionRepo)
	guestHandler := handlers.NewGuestHandler(categoryRepo, postRepo, commentRepo, reactionRepo)

	// Create middleware
	registerLimiter := middleware.NewRateLimiter()
	authMiddleware := middleware.NewAuthMiddleware(sessionRepo, userRepo)
	csrfMiddleware := middleware.CSRFMiddleware(sessionRepo)
	corsMiddleware := middleware.NewCORSMiddleware("http://localhost:8081")

	// Create router
	mux := http.NewServeMux()

	// Public routes
	mux.Handle("/forum/api/categories", corsMiddleware.Handler(http.HandlerFunc(categoryHandler.GetCategories)))
	mux.Handle("/forum/api/allData", corsMiddleware.Handler(http.HandlerFunc(guestHandler.GetGuestData)))
	mux.Handle("/forum/api/register", corsMiddleware.Handler(http.HandlerFunc(registerLimiter.Limit(authHandler.Register))))
	mux.Handle("/forum/api/session/login", corsMiddleware.Handler(http.HandlerFunc(authHandler.Login)))
	mux.Handle("/forum/api/session/logout", corsMiddleware.Handler(http.HandlerFunc(authHandler.Logout)))
	mux.Handle("/forum/api/session/verify", corsMiddleware.Handler(http.HandlerFunc(authHandler.VerifySession)))

	// Protected routes with CSRF
	protected := func(h http.Handler) http.Handler {
		return corsMiddleware.Handler(authMiddleware.RequireAuth(csrfMiddleware(h)))
	}

	mux.Handle("/forum/api/posts/create", protected(http.HandlerFunc(postHandler.CreatePost)))
	mux.Handle("/forum/api/user/posts", protected(http.HandlerFunc(myPostsHandler.GetMyPosts)))
	mux.Handle("/forum/api/user/liked", protected(http.HandlerFunc(likedPostsHandler.GetLikedPosts)))
	mux.Handle("/forum/api/comments", protected(http.HandlerFunc(commentHandler.CreateComment)))
	mux.Handle("/forum/api/react", protected(http.HandlerFunc(reactionHandler.React)))

	// Global auth injection
	return authMiddleware.Authenticate(mux)
}
