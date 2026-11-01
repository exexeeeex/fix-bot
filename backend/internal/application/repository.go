package application

import (
	"context"
	"fmt"

	"github.com/exexeeeex/fiks-api/internal/domain"
	"gorm.io/gorm"

	repo "github.com/exexeeeex/fiks-api/internal/repository"
)

type Repository interface {
	GetById(ctx context.Context, id string) (domain.Application, error)
	GetByGroupId(ctx context.Context, groupId string) ([]domain.Application, error)
	GetByMaxId(ctx context.Context, maxId string) (domain.Application, error)
	GetByUserId(ctx context.Context, userId string) (domain.Application, error)
	List(ctx context.Context, organizationId string) ([]domain.Application, error)
	Create(ctx context.Context, application *domain.Application) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// GetByUserId implements [Repository].
func (r *repository) GetByUserId(
	ctx context.Context,
	userId string,
) (domain.Application, error) {
	var result domain.Application

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userId).
		First(&result).
		Error

	if err != nil {
		return domain.Application{}, fmt.Errorf(
			"get by user id: %w",
			err,
		)
	}

	return result, nil
}

// GetByMaxId implements [Repository].
func (r *repository) GetByMaxId(ctx context.Context, maxId string) (domain.Application, error) {
	var result domain.Application

	err := r.db.WithContext(ctx).
		Model(&domain.Application{}).
		Joins("JOIN users ON users.id = application.user_id").
		Where("users.max_id = ?", maxId).
		Find(&result).
		Error
	if err != nil {
		return domain.Application{}, fmt.Errorf("get by max id: %w", err)
	}

	return result, nil
}

// GetByGroupId implements [Repository].
func (r *repository) GetByGroupId(ctx context.Context, groupId string) ([]domain.Application, error) {
	var applications []domain.Application

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Where("group_id = ?", groupId).
		Find(&applications).
		Error
	if err != nil {
		return nil, fmt.Errorf("get by group_id: %w", err)
	}
	return applications, nil
}

// GetById implements [Repository].
func (r *repository) GetById(ctx context.Context, id string) (domain.Application, error) {
	return repo.GetByID[domain.Application](ctx, r.db, id, "application")
}

// List implements [Repository].
func (r *repository) List(ctx context.Context, organizationId string) ([]domain.Application, error) {
	var applications []domain.Application

	err := r.db.WithContext(ctx).
		Where("applications.organization_id = ?", organizationId).
		Find(&applications).
		Error
	if err != nil {
		return nil, fmt.Errorf("list applications: %w", err)
	}
	return applications, nil
}

// Create implements [Repository].
func (r *repository) Create(ctx context.Context, application *domain.Application) error {
	return repo.CREATE(ctx, r.db, application, "application")
}
