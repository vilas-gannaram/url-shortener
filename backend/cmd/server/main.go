package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/vilas-gannaram/url-shortener/internal/db"
	"github.com/vilas-gannaram/url-shortener/internal/handlers"
)

func main() {
	godotenv.Load()
	
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is empty! Check your environment variables.")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}
	defer pool.Close()

	ctx := context.Background()
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Could not connect to Supabase: %v", err)
	}

	queries := db.New(pool)
	h := &handlers.URLHandler{
		Queries: queries,
		Pool:    pool,
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.CleanPath)
	r.Use(middleware.StripSlashes)

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("URL Shortener is Live"))
	})

	r.Post("/shorten", h.Shorten)
	r.Get("/{shortKey}", h.Redirect)
	// r.Get("/stats/{shortKey}", h.Stats)
	r.Get("/urls", h.ListURLs)

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", r)
}
