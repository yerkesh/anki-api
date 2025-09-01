package usecase

import (
	"context"

	"anki-api/internal/clients/chatgpt"
	"anki-api/internal/entity"
)

type userRepoer interface {
	CreateUser(ctx context.Context, user entity.User) (int32, error)
}

type UsersUsecase struct {
	userRepo userRepoer
}

func NewUsersUsecase(userRepo userRepoer) *UsersUsecase {
	return &UsersUsecase{
		userRepo: userRepo,
	}
}

type flashcardsRepoer interface {
	CreateFlashcard(ctx context.Context, card entity.Flashcard) (int32, error)
	GetFlashcards(ctx context.Context, collectionID int32, params entity.PageableQueryParams) ([]entity.Flashcard, error)
	GetFlashcardsTotal(ctx context.Context, collectionID int32) (int64, error)
	UpdateReview(ctx context.Context, flashcardID int, status entity.ReviewStatus) error
	GetFlashcard(ctx context.Context, flashcardID int) (entity.Flashcard, error)
	DeleteFlashcard(ctx context.Context, flashcardID int) error
}

type FlashcardsUsecase struct {
	flashcardsRepo  flashcardsRepoer
	collectionsRepo collectionsRepoer
	ChatGPTClient   *chatgpt.Client
}

func NewFlashcardsUsecase(usecase flashcardsRepoer, collectionRepo collectionsRepoer, gptClient *chatgpt.Client) *FlashcardsUsecase {
	return &FlashcardsUsecase{
		flashcardsRepo:  usecase,
		collectionsRepo: collectionRepo,
		ChatGPTClient:   gptClient,
	}
}

type collectionsRepoer interface {
	CreateCollections(ctx context.Context, collection entity.Collection) (int32, error)
	GetCollections(ctx context.Context, userID int32) ([]entity.Collection, error)
	GetCollection(ctx context.Context, collectionID int32) (entity.Collection, error)
}

type CollectionsUsecase struct {
	collectionsRepo collectionsRepoer
}

func NewCollectionsUsecase(usecase collectionsRepoer) *CollectionsUsecase {
	return &CollectionsUsecase{
		collectionsRepo: usecase,
	}
}
