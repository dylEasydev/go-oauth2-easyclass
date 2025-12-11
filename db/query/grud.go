package query

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// fonction génerique de création
func QueryCreate[T any](tx *gorm.DB, data *T) error {
	ctx := context.Background()
	return gorm.G[T](tx).Create(ctx, data)
}

func QueryDeleteById[T any](tx *gorm.DB, id uuid.UUID) error {
	ctx := context.Background()
	_, err := gorm.G[T](tx).Where("id = ?", id).Delete(ctx)
	return err
}
