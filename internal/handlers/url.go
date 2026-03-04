package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jaevor/go-nanoid"

	"github.com/vilas-gannaram/url-shortener/internal/storage"
	"gorm.io/gorm"
)

// URLHandler holds our database connection
type URLHandler struct {
	DB *gorm.DB
}

type ShortenRequest struct {
	URL string `json:"url"`
}

// We use a Base32-style alphabet (no O, 0, I, 1, L, u)
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

	// Collision-Resistant Insertion
	var shortKey string
	var finalErr error

	for i := 0; i < 3; i++ {
		shortKey = canonicNanoid()

		finalErr = h.DB.Transaction(func(tx *gorm.DB) error {
			newRecord := storage.URLMapping{OriginalURL: req.URL, ShortKey: shortKey}
			if err := tx.Create(&newRecord).Error; err != nil {
				return err
			}

			stats := storage.URLStats{URLMappingID: newRecord.ID, RedirectedCount: 0}
			return tx.Create(&stats).Error
		})

		if finalErr == nil {
			break
		}
	}

	// Check if we actually succeeded
	if finalErr != nil {
		http.Error(w, "Database error: could not generate unique key", http.StatusInternalServerError)
		return
	}

	// Build response
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}

	fullURL := fmt.Sprintf("%s://%s/%s", scheme, r.Host, shortKey)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"short_url": fullURL,
	})
}

// Redirect handles GET /{shortKey}
func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {

	// Check for the "Prediction" headers
	purpose := r.Header.Get("Sec-Purpose")
	isFakeRequest := strings.Contains(purpose, "prefetch") || strings.Contains(purpose, "prerender")

	shortKey := chi.URLParam(r, "shortKey")
	var mapping storage.URLMapping

	// Fetching the mapping from DB
	if err := h.DB.Where("short_key = ?", shortKey).First(&mapping).Error; err != nil {
		http.NotFound(w, r)
		return
	}

	// Incrementing the count in background, making the redirect faster
	if !isFakeRequest {
		go func(id uint) {
			err := h.DB.Model(&storage.URLStats{}).
				Where("url_mapping_id = ?", id).
				UpdateColumn("redirected_count", gorm.Expr("redirected_count + ?", 1)).Error

			if err != nil {
				log.Println("Error updating stats:", err)
			}
		}(mapping.ID)
	}

	// Redirecting to the original URL
	http.Redirect(w, r, mapping.OriginalURL, http.StatusFound)
}

// Stats handles GET /stats/{shortKey}
func (h *URLHandler) Stats(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Stats endpoint not implemented yet"))
}
