package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vilas-gannaram/url-shortener/internal/handlers"
	"github.com/vilas-gannaram/url-shortener/internal/storage"
)

func main() {
	db := storage.InitDB()
	h := &handlers.URLHandler{DB: db}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.CleanPath)
	r.Use(middleware.StripSlashes)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("URL Shortener is Live"))
	})

	r.Post("/shorten", h.Shorten)
	r.Get("/{shortKey}", h.Redirect)
	r.Get("/stats/{shortKey}", h.Stats)

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", r)
}
