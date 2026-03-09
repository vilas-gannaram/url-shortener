package handlers

import (
	"encoding/json"
	"net/http"
)

// ListURLs handles GET /urls
func (h *URLHandler) ListURLs(w http.ResponseWriter, r *http.Request) {
	// var results []struct {
	// 	ID              uint   `json:"id"`
	// 	OriginalURL     string `json:"original_url"`
	// 	ShortKey        string `json:"short_key"`
	// 	RedirectedCount int    `json:"redirected_count"`
	// 	LastUpdated     int64  `json:"last_updated"`
	// }

	ctx := r.Context()
	results, err := h.Queries.ListURL(ctx)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
