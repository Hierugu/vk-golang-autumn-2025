package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"redditclone/internal/redditclone/auth"
)

func (h *UserHandler) RegisterApiHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"message": "Failed to read request body"}`, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bodyBytes, &req)
	if err != nil {
		http.Error(w, `{"message": "Failed to unmarshal JSON"}`, http.StatusBadRequest)
		return
	}

	u, err := h.Repo.Register(req.Username, req.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf(`"{"message": "Failed to register: %s"}"`, err), http.StatusUnauthorized)
		return
	}

	token, err := h.jwtManager.Generate(u)
	if err != nil {
		http.Error(w, `{"message": "Could not create token"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(auth.JWTResponse{Token: token})
}
