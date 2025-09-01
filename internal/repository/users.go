package repository

import (
	"context"
	"fmt"

	"anki-api/internal/entity"
	"anki-api/internal/repository/generated"
)

type UsersRepo struct {
	generated.Querier
}

func NewUsersRepo(querier generated.Querier) *UsersRepo {
	return &UsersRepo{querier}
}

func (u *UsersRepo) CreateUser(ctx context.Context, user entity.User) (int32, error) {
	id, err := u.InsertUser(ctx, generated.InsertUserParams{
		Email:    user.Email,
		Username: user.Username,
		Role:     user.Role,
	})
	if err != nil {
		return 0, fmt.Errorf("createUser: couldn't create user email: %s, Cause: %w", user.Email, err)
	}

	return id, nil
}
