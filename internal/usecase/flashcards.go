package usecase

import (
	"context"
	"fmt"

	"anki-api/internal/entity"
)

func (f *FlashcardsUsecase) CreateFlashcard(ctx context.Context, flashcard entity.Flashcard) (int32, error) {
	collection, err := f.collectionsRepo.GetCollection(ctx, flashcard.CollectionID)
	if err != nil {
		return 0, fmt.Errorf("collectionsRepo.GetCollections: %w", err)
	}

	res, err := f.ChatGPTClient.Response(ctx, flashcard.Frontside, collection.LearningLanguage, collection.NativeLanguage)
	if err != nil {
		return 0, fmt.Errorf("f.ChatGPTClient.Response: %w", err)
	}

	res.CollectionID = flashcard.CollectionID
	res.Frontside = flashcard.Frontside
	// parse and save flashcard
	return f.flashcardsRepo.CreateFlashcard(ctx, res)
}

func (f *FlashcardsUsecase) GetFlashcards(ctx context.Context, collectionID int32, params entity.PageableQueryParams) (entity.GetFlashcardsResponse, error) {
	cards, err := f.flashcardsRepo.GetFlashcards(ctx, collectionID, params)
	if err != nil {
		return entity.GetFlashcardsResponse{}, fmt.Errorf("flashcardsRepo.GetFlashcards: %w", err)
	}

	total, err := f.flashcardsRepo.GetFlashcardsTotal(ctx, collectionID)
	if err != nil {
		return entity.GetFlashcardsResponse{}, fmt.Errorf("flashcardsRepo.GetFlashcardsTotal: %w", err)
	}

	return entity.GetFlashcardsResponse{
		Flashcards:     cards,
		PaginationMeta: entity.NewPaginationMeta(uint64(total), params.Size, params.Page),
	}, nil
}

func (f *FlashcardsUsecase) UpdateFlashcardStatus(ctx context.Context, flashcardID int, status entity.ReviewStatus) error {
	return f.flashcardsRepo.UpdateReview(ctx, flashcardID, status)
}

func (f *FlashcardsUsecase) DeleteFlashcard(ctx context.Context, flashcardID int) error {
	return f.flashcardsRepo.DeleteFlashcard(ctx, flashcardID)
}

func (f *FlashcardsUsecase) GetFlashcard(ctx context.Context, flashcardID int) (entity.Flashcard, error) {
	return f.flashcardsRepo.GetFlashcard(ctx, flashcardID)
}
