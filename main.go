package main

import (
	"log"
	"math/rand"
	"net/http"
	"reviews/internal/handler"
	"reviews/internal/repository/memory"
	"reviews/internal/service"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	repo := memory.New()
	reviewService := service.NewReviewService(repo)
	authService := service.NewAuthService("admin", "password")

	authHandler := handler.NewAuthHandler(authService)
	reviewHandler := handler.NewReviewHandler(reviewService, authService)

	http.HandleFunc("/login", authHandler.Login)
	http.Handle("/reviews/", reviewHandler)
	http.Handle("/reviews", reviewHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
