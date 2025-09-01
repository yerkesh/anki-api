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

func (h *Handler) CreateFlashcard(w http.ResponseWriter, r *http.Request) {
	var in entity.Flashcard
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "json.NewDecode: "+err.Error(), http.StatusBadRequest)
		return
	}

	collectionID, err := strconv.Atoi(chi.URLParam(r, "collectionID"))
	if err != nil {
		http.Error(w, "strconv.Atoi: "+err.Error(), http.StatusBadRequest)
	}

	in.CollectionID = int32(collectionID)

	if err = h.validator.StructCtx(r.Context(), in); err != nil {
		http.Error(w, "validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.flashcard.CreateFlashcard(r.Context(), in)
	if err != nil {
		http.Error(w, "cannot create flashcard: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	_ = json.NewEncoder(w).Encode(entity.IDResp{ID: int(id)})
}

func (h *Handler) GetFlashcards(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "collectionID")
	collectionID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid collection id", http.StatusBadRequest)
		return
	}

	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		log.Printf("invalid page params: %v, err: %v", pageStr, err)
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		log.Printf("invalid page_size params: %v, err: %v", pageSize, err)
	}

	pageParams := entity.PageableQueryParams{
		Page: uint64(page),
		Size: uint64(pageSize),
	}
	pageParams.Normalize()

	collections, err := h.flashcard.GetFlashcards(ctx, int32(collectionID), pageParams)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		http.Error(w, "no flashcards found", http.StatusNotFound)
		return
	case err != nil:
		log.Printf("get flashcards: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(collections); err != nil {
		log.Printf("encode flashcards: %v", err)
	}
}

func (h *Handler) UpdateFlashcardStatus(w http.ResponseWriter, r *http.Request) {
	type updateFlashcardStatusReq struct {
		Status string `json:"status" validate:"required,oneof=easy hard repeat"`
	}

	var in updateFlashcardStatusReq
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "json.NewDecode: "+err.Error(), http.StatusBadRequest)
		return
	}

	flashcardID, err := strconv.Atoi(chi.URLParam(r, "flashcardID"))
	if err != nil {
		http.Error(w, "strconv.Atoi: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.validator.StructCtx(r.Context(), in); err != nil {
		http.Error(w, "validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.flashcard.UpdateFlashcardStatus(r.Context(), flashcardID, entity.ReviewStatus(in.Status)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "cannot update status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204
}

func (h *Handler) DeleteFlashcard(w http.ResponseWriter, r *http.Request) {
	flashcardID, err := strconv.Atoi(chi.URLParam(r, "flashcardID"))
	if err != nil {
		http.Error(w, "strconv.Atoi: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.flashcard.DeleteFlashcard(r.Context(), flashcardID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "cannot delete flashcard: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
func (h *Handler) GetFlashcard(w http.ResponseWriter, r *http.Request) {
	flashcardID, err := strconv.Atoi(chi.URLParam(r, "flashcardID"))
	if err != nil {
		http.Error(w, "strconv.Atoi: "+err.Error(), http.StatusBadRequest)
		return
	}

	card, err := h.flashcard.GetFlashcard(r.Context(), flashcardID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		http.Error(w, "no flashcards found", http.StatusNotFound)
		return
	case err != nil:
		log.Printf("get flashcards: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(card); err != nil {
		log.Printf("encode flashcards: %v", err)
	}
}
