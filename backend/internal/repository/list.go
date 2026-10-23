package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func List[T any](
	ctx context.Context,
	db *gorm.DB,
	entityName string,
) ([]T, error) {
	var entity []T

	err := db.WithContext(ctx).
		Order("created_at DESC").
		Find(&entity).
		Error
	
	if err != nil {
		var zero []T

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return zero, fmt.Errorf("%s not found", entityName)
		}

		return zero, fmt.Errorf("list %s: %w", entityName, err)
	}

	return entity, nil
}