package handlers

import (
	"context"

	"anki-api/internal/entity"
)

type structValidator interface {
	StructCtx(ctx context.Context, s interface{}) error
}

type flashcardUsecase interface {
	CreateFlashcard(ctx context.Context, flashcard entity.Flashcard) (int32, error)
	GetFlashcards(ctx context.Context, collectionID int32, params entity.PageableQueryParams) (entity.GetFlashcardsResponse, error)
	UpdateFlashcardStatus(ctx context.Context, flashcardID int, status entity.ReviewStatus) error
	DeleteFlashcard(ctx context.Context, flashcardID int) error
	GetFlashcard(ctx context.Context, flashcardID int) (entity.Flashcard, error)
}

type userUsecase interface {
	CreateUser(ctx context.Context, user entity.User) (int32, error)
}

type collectionUsecase interface {
	CreateCollection(ctx context.Context, collection entity.Collection) (int32, error)
	GetCollections(ctx context.Context, userID int32) ([]entity.Collection, error)
}

type Options struct {
	validator  structValidator
	flashcard  flashcardUsecase
	user       userUsecase
	collection collectionUsecase
}

type Handler struct {
	Options
}

func NewHandler(o Options) *Handler {
	return &Handler{
		Options: o,
	}
}
