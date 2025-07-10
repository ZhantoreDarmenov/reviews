package main

import (
	_ "context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/google/uuid"
	_ "github.com/joho/godotenv"
	_ "google.golang.org/api/option"
	"log"
	"net/http"

	"reviews/internal/handlers"
	_ "reviews/internal/handlers"
	_ "reviews/internal/models"
	"reviews/internal/repositories"
	_ "reviews/internal/repositories"
	"reviews/internal/services"
	_ "reviews/internal/services"
	_ "reviews/utils"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	userHandler   *handlers.UserHandler
	userRepo      *repositories.UserRepository
	reviewHandler *handlers.ReviewHandler
	reviewRepo    *repositories.ReviewRepository
}

func initializeApp(db *sql.DB, errorLog, infoLog *log.Logger) *application {
	// Repositories
	userRepo := repositories.UserRepository{DB: db}
	reviewRepo := repositories.ReviewRepository{DB: db}

	// Services
	userService := &services.UserService{UserRepo: &userRepo}
	reviewService := &services.ReviewService{Repo: &reviewRepo}

	// Handlers
	userHandler := &handlers.UserHandler{Service: userService}
	reviewHandler := &handlers.ReviewHandler{Service: reviewService}

	return &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		userHandler:   userHandler,
		userRepo:      &userRepo,
		reviewHandler: reviewHandler,
		reviewRepo:    &reviewRepo,
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Failed to open DB: %v", err)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		log.Printf("Failed to ping DB: %v", err)
		return nil, err
	}
	db.SetMaxIdleConns(35)
	log.Println("Successfully connected to database")
	return db, nil
}

func addSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
		next.ServeHTTP(w, r)
	})
}
