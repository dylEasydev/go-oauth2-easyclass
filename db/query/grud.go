package query

import (
	"context"

	"gorm.io/gorm"
)

// fonction génerique de création
func QueryCreate[T any](tx *gorm.DB, data *T) error {
	ctx := context.Background()
	return gorm.G[T](tx).Create(ctx, data)
}
