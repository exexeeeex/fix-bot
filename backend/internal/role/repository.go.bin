package role

import (
	"context"

	"github.com/exexeeeex/fiks-api/internal/domain"
	repo "github.com/exexeeeex/fiks-api/internal/repository"
	"gorm.io/gorm"
)

type Repository interface {
	GetByName(ctx context.Context, name string) (domain.Role, error)
	GetById(ctx context.Context, id string) (domain.Role, error)
	List(ctx context.Context) ([]domain.Role, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// GetById implements [Repository].
func (r *repository) GetById(ctx context.Context, id string) (domain.Role, error) {
	return repo.GetByID[domain.Role](ctx, r.db, id, "role")
}

// GetByName implements [Repository].
func (r *repository) GetByName(ctx context.Context, name string) (domain.Role, error) {
	var role domain.Role

	err := r.db.WithContext(ctx).
		Model(&domain.Role{}).
		Where("name ILIKE ?", name).
		Find(&role).
		Error
	if err != nil {
		return domain.Role{}, err
	}

	return role, nil
}

// List implements [Repository].
func (r *repository) List(ctx context.Context) ([]domain.Role, error) {
	var roles []domain.Role

	err := r.db.WithContext(ctx).
		Model(&domain.Role{}).
		Find(&roles).
		Error

	if err != nil {
		return nil, err
	}

	return roles, nil
}
