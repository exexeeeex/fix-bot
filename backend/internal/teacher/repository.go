package teacher

import (
	"context"

	repo "github.com/exexeeeex/fiks-api/internal/repository"

	"github.com/exexeeeex/fiks-api/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	GetById(ctx context.Context, id string) (domain.Teacher, error)
	GetByName(ctx context.Context, name string) ([]domain.Teacher, error)
	List(ctx context.Context, organizationId string) ([]domain.Teacher, error)
	Create(ctx context.Context, teacher *domain.Teacher) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create implements [Repository].
func (r *repository) Create(ctx context.Context, teacher *domain.Teacher) error {
	return repo.CREATE(ctx, r.db, teacher, "teacher")
}

// GetById implements [Repository].
func (r *repository) GetById(ctx context.Context, id string) (domain.Teacher, error) {
	return  repo.GetByID[domain.Teacher](ctx, r.db, id, "teacher")
}

// GetByName implements [Repository].
func (r *repository) GetByName(ctx context.Context, name string) ([]domain.Teacher, error) {
	panic("unimplemented")
}

// List implements [Repository].
func (r *repository) List(ctx context.Context, organizationId string) ([]domain.Teacher, error) {
	panic("unimplemented")
}
