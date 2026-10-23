package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func GetByID[T any](
	ctx context.Context,
	db *gorm.DB,
	id string,
	entityName string,
) (T, error) {
	var entity T

	err := db.WithContext(ctx).
		Where("id = ?", id).
		First(&entity).
		Error
	
	if err != nil {
		var zero T
		
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return zero, fmt.Errorf("%s not found!", entityName)
		}

		return zero, fmt.Errorf("get %s by id: %w", entityName, err)
	}

	return entity, nil
}