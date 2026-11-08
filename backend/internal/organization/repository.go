package organization

import (
	"context"
	"fmt"

	"github.com/exexeeeex/fiks-api/internal/domain"
	repo "github.com/exexeeeex/fiks-api/internal/repository"
	"gorm.io/gorm"
)

type Repository interface {
	GetById(ctx context.Context, id string) (domain.Organization, error)
	GetByName(ctx context.Context, name string) ([]domain.Organization, error)
	List(ctx context.Context) ([]domain.Organization, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

// GetById implements [Repository].
func (r *repository) GetById(ctx context.Context, id string) (domain.Organization, error) {
	return repo.GetByID[domain.Organization](ctx, r.db, id, "organization")
}

// GetByName implements [Repository].
func (r *repository) GetByName(ctx context.Context, name string) ([]domain.Organization, error) {	
	var organizations []domain.Organization

	err := r.db.WithContext(ctx).
		Where("name ILIKE ?", "%"+name+"%").
		Order("name DESC").
		Find(&organizations).Error
	if err != nil {
		return []domain.Organization{}, fmt.Errorf("get organizations by name: %w", err)
	}

	return organizations, nil
}

// List implements [Repository].
func (r *repository) List(ctx context.Context) ([]domain.Organization, error) {
	return repo.List[domain.Organization](ctx, r.db, "organization")
}
