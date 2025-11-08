package main

import (
	"log"
	"net/http"
	"os"
	"redditclone/internal/redditclone/db"
	"redditclone/internal/redditclone/handlers"
	"redditclone/internal/redditclone/repository"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	var userRepo repository.UserRepository = db.CreateMemoryUserRepository()
	uh := handlers.NewUserHandler(userRepo, os.Getenv("JWT_SECRET"))

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/api/register", uh.RegisterApiHandler)
	apiMux.HandleFunc("/api/login", uh.LoginApiHandler)
	apiMux.HandleFunc("/api/health", handlers.HealthApiHandler)

	staticHTMLHandler := http.FileServer(http.Dir("web/html"))
	staticHandler := http.FileServer(http.Dir("web"))

	staticMux := http.NewServeMux()
	staticMux.Handle("/", staticHTMLHandler)
	staticMux.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	siteMux := http.NewServeMux()
	siteMux.Handle("/", staticMux)
	siteMux.Handle("/api/", apiMux)

	server := &http.Server{
		Addr: ":8080",
		//Handler:      stripTrailingSlash(mux),
		Handler:      siteMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Print("Starting server on :8080")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
