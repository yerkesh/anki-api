package repository

import (
	"context"
	"fmt"

	"anki-api/internal/entity"
	"anki-api/internal/repository/generated"
)

type CollectionsRepo struct {
	generated.Querier
}

func NewCollectionsRepo(queries generated.Querier) *CollectionsRepo {
	return &CollectionsRepo{queries}
}

func (c *CollectionsRepo) CreateCollections(ctx context.Context, collection entity.Collection) (int32, error) {
	id, err := c.InsertCollection(ctx, generated.InsertCollectionParams{
		UserID:           collection.UserID,
		Name:             collection.Name,
		NativeLanguage:   string(collection.NativeLanguage), // TODO: handle language
		LearningLanguage: string(collection.LearningLanguage),
	})

	if err != nil {
		return 0, fmt.Errorf("createCollection: couldn't create collection, userID %d, Cause: %w", collection.UserID, err)
	}

	return id, nil
}

func (c *CollectionsRepo) GetCollections(ctx context.Context, userID int32) ([]entity.Collection, error) {
	collections, err := c.SelectCollections(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getCollections: couldn't get collections, userID: %d, Cause: %w", userID, err)
	}

	return collectionsToEntity(collections), nil
}

func (c *CollectionsRepo) GetCollection(ctx context.Context, collectionID int32) (entity.Collection, error) {
	collection, err := c.SelectCollection(ctx, collectionID)
	if err != nil {
		return entity.Collection{}, fmt.Errorf("c.SelectCollection: couldn't get collection, collectionID: %d, err: %w", collectionID, err)
	}

	return entity.Collection{
		ID:               collection.ID,
		UserID:           collection.UserID,
		Name:             collection.Name,
		NativeLanguage:   entity.Language(collection.NativeLanguage),
		LearningLanguage: entity.Language(collection.LearningLanguage),
	}, nil
}

func collectionsToEntity(collections []generated.SelectCollectionsRow) []entity.Collection {
	res := make([]entity.Collection, 0, len(collections))
	for _, clc := range collections {
		res = append(res, entity.Collection{
			ID:               clc.ID,
			Name:             clc.Name,
			UserID:           clc.UserID,
			NativeLanguage:   entity.Language(clc.NativeLanguage), // TODO: handle unknown
			LearningLanguage: entity.Language(clc.LearningLanguage),
		})
	}

	return res
}
