package handlers

import (
	"encoding/json"
	"net/http"

	"anki-api/internal/entity"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var in entity.User
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "json.NewDecode: "+err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.validator.StructCtx(r.Context(), in); err != nil {
		http.Error(w, "validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.user.CreateUser(r.Context(), in)
	if err != nil {
		http.Error(w, "cannot create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	_ = json.NewEncoder(w).Encode(id)
}
