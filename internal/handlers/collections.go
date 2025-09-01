package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"anki-api/internal/entity"
)

func (h *Handler) CreateCollection(w http.ResponseWriter, r *http.Request) {
	var in entity.Collection
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "json.NewDecode: "+err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "strconv.Atoi: "+err.Error(), http.StatusBadRequest)
	}

	in.UserID = int32(userID)

	if err = h.validator.StructCtx(r.Context(), in); err != nil {
		http.Error(w, "validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.collection.CreateCollection(r.Context(), in)
	if err != nil {
		http.Error(w, "cannot create collection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	_ = json.NewEncoder(w).Encode(id)
}

func (h *Handler) GetCollections(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	collections, err := h.collection.GetCollections(ctx, int32(userID))
	switch {
	case errors.Is(err, sql.ErrNoRows):
		http.Error(w, "no collections found", http.StatusNotFound)
		return
	case err != nil:
		log.Printf("get collections: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(collections); err != nil {
		log.Printf("encode collections: %v", err)
	}
}
