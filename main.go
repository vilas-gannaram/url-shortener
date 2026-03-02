package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vilas-gannaram/url-shortener/handlers"
	"github.com/vilas-gannaram/url-shortener/storage"
)

func main() {
	db := storage.InitDB()

	// Initialize our handlers with the DB connection
	h := &handlers.URLHandler{DB: db}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.CleanPath)
	r.Use(middleware.StripSlashes)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("URL Shortener is Live"))
	})

	// Use the methods from our handler struct
	r.Post("/shorten", h.ShortenURL)
	r.Get("/{shortKey}", h.Redirect)

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", r)
}
