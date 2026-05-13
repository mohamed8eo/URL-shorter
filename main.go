package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/mohamed8eo/url-shortener/handler"
	"github.com/mohamed8eo/url-shortener/internal/database"
	"github.com/mohamed8eo/url-shortener/middleware"
	"github.com/mohamed8eo/url-shortener/storage"

	"github.com/jackc/pgx/v5"
)

type apiConfig struct {
	dbQueries *database.Queries
	PORT      string
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	PORT := os.Getenv("PORT")

	dbQuerie := database.New(conn)
	cfg := &apiConfig{
		dbQueries: dbQuerie,
		PORT:      PORT,
	}

	mux := http.NewServeMux()

	storage := storage.NewStorage()
	handler := handler.NewHandler(storage, cfg.dbQueries, cfg.PORT)

	limiter := middleware.NewRateLimit(100, time.Minute)

	mux.HandleFunc("GET /", handler.Redirect)
	mux.HandleFunc("POST /create", handler.CreateShortURL)

	log.Print("DB is running")
	log.Printf("Server is running on PORT: %s\n", PORT)

	if err := http.ListenAndServe(":"+PORT, limiter.Middleware(mux)); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
