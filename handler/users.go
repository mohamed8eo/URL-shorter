package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/mohamed8eo/url-shortener/internal/database"
)

type UsersHandler struct {
	DB *database.Queries
}

func (h *UsersHandler) HandlerSignUp(w http.ResponseWriter, r *http.Request) {
	type requeset struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req requeset

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	hashed, err := argon2id.CreateHash()
}
