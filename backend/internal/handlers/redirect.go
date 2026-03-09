package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

// Redirect handles GET /{shortKey}
func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	// Check for the "Prediction" headers
	purpose := r.Header.Get("Sec-Purpose")
	isFakeRequest := strings.Contains(purpose, "prefetch") || strings.Contains(purpose, "prerender")

	shortKey := chi.URLParam(r, "shortKey")

	// Fetching the mapping from DB
	ctx := r.Context()
	urlMapping, err := h.Queries.GetURLByCode(ctx, shortKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Incrementing the count in background, making the redirect faster
	if !isFakeRequest {
		go func(urlID int64) {
			// Use background context for background task to ensure it completes
			// even if the original request context is cancelled.
			err := h.Queries.UpsertStats(context.Background(), urlID)
			if err != nil {
				log.Println("Error updating stats:", err)
			}
		}(urlMapping.ID)
	}

	// Redirecting to the original URL
	http.Redirect(w, r, urlMapping.OriginalUrl, http.StatusFound)
}

// Stats handles GET /stats/{shortKey}
func (h *URLHandler) Stats(w http.ResponseWriter, r *http.Request) {
	shortKey := chi.URLParam(r, "shortKey")

	stats, err := h.Queries.GetStatsByCode(r.Context(), shortKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
