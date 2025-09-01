package usecase

import (
	"context"

	"anki-api/internal/entity"
)

func (c *CollectionsUsecase) CreateCollection(ctx context.Context, collection entity.Collection) (int32, error) {
	return c.collectionsRepo.CreateCollections(ctx, collection)
}

func (c *CollectionsUsecase) GetCollections(ctx context.Context, userID int32) ([]entity.Collection, error) {
	return c.collectionsRepo.GetCollections(ctx, userID)
}
