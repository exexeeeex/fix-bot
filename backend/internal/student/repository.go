package student

import (
	"context"
	"fmt"

	"github.com/exexeeeex/fiks-api/internal/domain"
	repo "github.com/exexeeeex/fiks-api/internal/repository"
	"gorm.io/gorm"
)

type Repository interface {
	GetById(ctx context.Context, id string) (domain.Student, error)
	GetByGroupId(ctx context.Context, groupId string) ([]domain.Student, error)
	List(ctx context.Context, organizationId string) ([]domain.Student, error)
	Create(ctx context.Context, student *domain.Student) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetById(ctx context.Context, id string) (domain.Student, error) {
	return repo.GetByID[domain.Student](ctx, r.db, id, "student")
}

func (r *repository) GetByGroupId(ctx context.Context, groupId string) ([]domain.Student, error) {
	var students []domain.Student

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Where("group_id = ?", groupId).
		Find(&students).
		Error
	if err != nil {
		return nil, fmt.Errorf("get by group_id: %w", err)
	}
	return students, nil
}

func (r *repository) List(ctx context.Context, organizationId string) ([]domain.Student, error) {
	var students []domain.Student

	err := r.db.WithContext(ctx).
		Model(&domain.Student{}).
		Joins("JOIN groups ON groups.id = student.group.id").
		Where("groups.organization_id = ?", organizationId).
		Find(&students).
		Error
	if err != nil {
		return nil, fmt.Errorf("list students: %w", err)
	}

	return students, nil
}

func (r *repository) Create(ctx context.Context, student *domain.Student) error {
	return repo.CREATE(ctx, r.db, student, "student")
}
