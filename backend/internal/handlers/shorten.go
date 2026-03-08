package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaevor/go-nanoid"
	"github.com/vilas-gannaram/url-shortener/internal/db"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type URLHandler struct {
	Queries *db.Queries
	Pool    *pgxpool.Pool
}

// Base32-style alpha-numeric characters (without O, 0, I, 1, L, u).
// Using go-nanoid for random string generation
var canonicNanoid, _ = nanoid.CustomASCII("abcdefghjkmnpqrstvwxyz23456789", 8)

// Shorten handles POST /shorten
func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	// Decode Request Body JSON
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate URL
	u, err := url.ParseRequestURI(req.URL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	// Prevent Self-Shortening (Loop Check)
	if u.Host == r.Host {
		http.Error(w, "Cannot shorten URLs from this domain", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	var shortKey string
	var success bool
	var lastErr error

	for i := 0; i < 3; i++ {
		shortKey = canonicNanoid()

		// 1. Start a NEW transaction for each attempt
		tx, err := h.Pool.Begin(ctx)
		if err != nil {
			lastErr = err
			continue
		}

		// 2. Try the insert
		qtx := h.Queries.WithTx(tx)
		_, err = qtx.CreateURL(ctx, db.CreateURLParams{
			ShortCode:   shortKey,
			OriginalUrl: req.URL,
		})

		if err == nil {
			// 3. Commit only on success
			if err = tx.Commit(ctx); err == nil {
				success = true
				break
			}
		}

		// 4. If we reach here, something failed. Rollback this attempt.
		lastErr = err
		tx.Rollback(ctx)
	}

	if !success {
		// Log the actual error to your terminal/Render logs
		fmt.Printf("Final failure after 3 attempts. Last error: %v\n", lastErr)
		http.Error(w, "Could not generate unique short code", http.StatusConflict)
		return
	}

	// Build response...
	fullURL := fmt.Sprintf("https://%s/%s", r.Host, shortKey)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"short_url": fullURL})
}
