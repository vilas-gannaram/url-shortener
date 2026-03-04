package handlers

import (
	"encoding/json"
	"net/http"

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

	newRecord := storage.URLMapping{OriginalURL: req.URL}
	if err := h.DB.Create(&newRecord).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	stats := storage.URLStats{URLMappingID: newRecord.ID, RedirectedCount: 0}
	h.DB.Create(&stats)

	shortKey := utils.Encode(int64(newRecord.ID))
	h.DB.Model(&newRecord).Update("ShortKey", shortKey)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"short_url": "http://localhost:8080/" + shortKey,
	})
}

// Redirect handles GET /{shortKey}
func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortKey := chi.URLParam(r, "shortKey")
	var mapping storage.URLMapping

	if err := h.DB.Where("short_key = ?", shortKey).First(&mapping).Error; err != nil {
		http.NotFound(w, r)
		return
	}

	h.DB.Model(&storage.URLStats{}).
		Where("url_mapping_id = ?", mapping.ID).
		Update("redirected_count", gorm.Expr("redirected_count + ?", 1))

	http.Redirect(w, r, mapping.OriginalURL, http.StatusFound)
}

// Stats handles GET /stats/{shortKey}
func (h *URLHandler) Stats(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Stats endpoint not implemented yet"))
}
