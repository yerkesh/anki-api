package usecase

import (
	"context"

	"anki-api/internal/entity"
)

func (u *UsersUsecase) CreateUser(ctx context.Context, user entity.User) (int32, error) {
	return u.userRepo.CreateUser(ctx, user)
}
