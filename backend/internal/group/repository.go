package group

import (
	"context"
	"errors"
	"fmt"

	"github.com/exexeeeex/fiks-api/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	GetById(ctx context.Context, id string) (domain.Group, error)
	GetByName(ctx context.Context, name string, organizationId string) ([]domain.Group, error)
	List(ctx context.Context, organizationId string) ([]domain.Group, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// GetById implements [Repository].
func (r *repository) GetById(ctx context.Context, id string) (domain.Group, error) {
	var group domain.Group

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&group).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Group{}, errors.New("group not found")
		}

		return domain.Group{}, fmt.Errorf("get group by id: %w", err)
	}

	return group, nil
}

// GetByName implements [Repository].
func (r *repository) GetByName(
	ctx context.Context, 
	name string, 
	organizationId string,
) ([]domain.Group, error) {
	var groups []domain.Group

	err := r.db.WithContext(ctx).
		Where("name ILIKE ? AND organization_id = ?", "%"+name+"%", organizationId).
		Order("name DESC").
		Find(&groups).Error

	if err != nil {
		return nil, err
	}

	return groups, err
}

// List implements [Repository].
func (r *repository) List(ctx context.Context, organizationId string) ([]domain.Group, error) {
	var groups []domain.Group

	err := r.db.WithContext(ctx).
		Where("organization_id = ?", organizationId).
		Order("created_at DESC").
		Find(&groups).
		Error
	if err != nil {
		return nil, fmt.Errorf("list groups: %w", err)
	}

	return groups, nil
}
