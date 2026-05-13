package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mohamed8eo/url-shortener/internal/database"
	"github.com/mohamed8eo/url-shortener/model"
	"github.com/mohamed8eo/url-shortener/storage"
)

type Handler struct {
	store  *storage.Storage
	querie *database.Queries
	port   string
}

type CreateRequest struct {
	URL string `json:"url"`
}

type CreateResponse struct {
	ShortURL string `json:"short_url"`
	FullURL  string `json:"full_url"`
}

func NewHandler(store *storage.Storage, querie *database.Queries, port string) *Handler {
	return &Handler{
		store:  store,
		querie: querie,
		port:   port,
	}
}

func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}
	longURL := req.URL

	if !strings.HasPrefix(longURL, "http://") && !strings.HasPrefix(longURL, "https://") {
		longURL = "https://" + longURL
	}

	// First check if the longURL already exist or not
	existing, err := h.querie.GetURLByLongURL(r.Context(), longURL)
	if err == nil {
		fullURL := fmt.Sprintf("http://localhost:%s/%s", h.port, existing.ShortCode)
		res := CreateResponse{
			ShortURL: existing.ShortCode,
			FullURL:  fullURL,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
		return

	}

	shortURL, err := model.GenerateShortURL()
	if err != nil {
		http.Error(w, "faild to create the shortURL", http.StatusInternalServerError)
		return
	}

	_, err = h.querie.AddURL(r.Context(), database.AddURLParams{
		ShortCode:   shortURL,
		OriginalUrl: longURL,
	})
	if err != nil {
		http.Error(w, "failed to store URL", http.StatusInternalServerError)
		return
	}

	fullURL := fmt.Sprintf("http://localhost:%s/%s", h.port, shortURL)

	res := CreateResponse{
		ShortURL: shortURL,
		FullURL:  fullURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:]

	longURL, err := h.querie.GetLongURL(r.Context(), shortURL)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	h.querie.IncrementClicks(r.Context(), shortURL)

	http.Redirect(w, r, longURL, http.StatusFound)
}
