package service

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func FindUserByName[T any](ctx context.Context, tx *gorm.DB, name, email string) (*T, error) {
	user, err := gorm.G[T](tx).Where("user_name = ? or email = ?", name, email).First(ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func FindUserById[T any](ctx context.Context, tx *gorm.DB, id uuid.UUID) (*T, error) {
	user, err := gorm.G[T](tx).Where("id = ? ", id).First(ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
