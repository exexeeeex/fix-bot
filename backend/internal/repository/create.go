package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

func CREATE[T any](
	ctx context.Context,
	db *gorm.DB,
	entity *T,
	entityName string,
) error {
	err := db.WithContext(ctx).
		Create(entity).
		Error
	if err != nil {
		return fmt.Errorf("create %s: %w", entityName, err)
	}

	return nil
}