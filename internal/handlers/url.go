package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/vilas-gannaram/url-shortener/internal/storage"
	"github.com/vilas-gannaram/url-shortener/internal/utils"
	"gorm.io/gorm"
)

// URLHandler holds our database connection
type URLHandler struct {
	DB *gorm.DB
}

type ShortenRequest struct {
	URL string `json:"url"`
}

// Shorten handles POST /shorten
func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Parse the URL
	u, err := url.ParseRequestURI(req.URL)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// Validate Scheme (Protocol)
	if u.Scheme != "http" && u.Scheme != "https" {
		http.Error(w, "Only HTTP and HTTPS protocols are supported", http.StatusBadRequest)
		return
	}

	// Validate Host -- Ensure the host isn't empty (e.g., "https://")
	if u.Host == "" || !strings.Contains(u.Host, ".") {
		http.Error(w, "URL must have a valid domain", http.StatusBadRequest)
		return
	}

	var shortKey string

	// DB Transaction
	err = h.DB.Transaction(func(tx *gorm.DB) error {

		newRecord := storage.URLMapping{OriginalURL: req.URL}
		if err := tx.Create(&newRecord).Error; err != nil {
			return err
		}

		shortKey = utils.Encode(int64(newRecord.ID))
		if err := tx.Model(&newRecord).Update("ShortKey", shortKey).Error; err != nil {
			return err
		}

		stats := storage.URLStats{URLMappingID: newRecord.ID, RedirectedCount: 0}
		if err := tx.Create(&stats).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	domain := r.Host
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	fullURL := fmt.Sprintf("%s://%s/%s", scheme, domain, shortKey)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"short_url": fullURL,
	})
}

// Redirect handles GET /{shortKey}
func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortKey := chi.URLParam(r, "shortKey")
	var mapping storage.URLMapping

	// Fetching the mapping from DB
	if err := h.DB.Where("short_key = ?", shortKey).First(&mapping).Error; err != nil {
		http.NotFound(w, r)
		return
	}

	// Incrementing the count in background, making the redirect faster
	go func(id uint) {
		err := h.DB.Model(&storage.URLStats{}).
			Where("url_mapping_id = ?", id).
			UpdateColumn("redirected_count", gorm.Expr("redirected_count + ?", 1)).Error

		if err != nil {
			log.Println("Error updating stats:", err)
		}
	}(mapping.ID)

	// Redirecting to the original URL
	http.Redirect(w, r, mapping.OriginalURL, http.StatusFound)
}

// Stats handles GET /stats/{shortKey}
func (h *URLHandler) Stats(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Stats endpoint not implemented yet"))
}
