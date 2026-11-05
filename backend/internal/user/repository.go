package user

import (
	"context"
	"fmt"

	repo "github.com/exexeeeex/fiks-api/internal/repository"

	"github.com/exexeeeex/fiks-api/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	GetById(ctx context.Context, id string) (domain.User, error)
	GetByMaxId(ctx context.Context, maxId string) (domain.User, error)
	List(ctx context.Context, organizationId string) ([]domain.User, error)
	Create(ctx context.Context, user *domain.User) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// GetByMaxId implements [Repository].
func (r *repository) GetByMaxId(ctx context.Context, maxId string) (domain.User, error) {
	var user domain.User

	err := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("max_id = ?", maxId).
		First(&user).
		Error

	if err != nil {
		return domain.User{}, fmt.Errorf("get user by max_id: %w", err)
	}

	return user, nil
}

// Create implements [Repository].
func (r *repository) Create(ctx context.Context, user *domain.User) error {
	return repo.CREATE(ctx, r.db, user, "user")
}

// GetById implements [Repository].
func (r *repository) GetById(ctx context.Context, id string) (domain.User, error) {
	return repo.GetByID[domain.User](ctx, r.db, id, "user")
}

// List implements [Repository].
func (r *repository) List(ctx context.Context, organizationId string) ([]domain.User, error) {
	var users []domain.User

	err := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Joins("JOIN organizations ON organizations.id = user.organization_id").
		Where("organizations.id = ?", organizationId).
		Find(&users).
		Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
