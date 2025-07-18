package main

import (
	_ "database/sql"
	"flag"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"reviews/internal/config"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}
	cfg := config.LoadConfig()

	port := os.Getenv("PORT")
	if port == "" {
		port = ":4001"
	} else {
		port = ":" + port
	}

	addr := flag.String("addr", port, "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(cfg.Database.URL)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := initializeApp(db, errorLog, infoLog)

	fs := http.FileServer(http.Dir("./uploads"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"http://localhost:5173",
			"http://localhost:5174",
			"https://preview-admin-panel-creation-kzmqzeigysqrkgsy8lzt.vusercontent.net",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	})

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      addSecurityHeaders(c.Handler(app.routes())),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	if err := srv.ListenAndServe(); err != nil {
		errorLog.Fatal(err)
	}
}
