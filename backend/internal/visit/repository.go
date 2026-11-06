package visit

import (
	"context"

	"github.com/exexeeeex/fiks-api/internal/domain"
	repo "github.com/exexeeeex/fiks-api/internal/repository"

	"gorm.io/gorm"
)

type Repository interface {
	GetById(ctx context.Context, id string) (domain.Visit, error)
	ListByStudentId(ctx context.Context, studentId string) ([]domain.Visit, error)
	ListByGroupId(ctx context.Context, groupId string) ([]domain.Visit, error)
	ListByOrganizationId(ctx context.Context, organizationId string) ([]domain.Visit, error)
	Create(ctx context.Context, visit *domain.Visit) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// GetById implements [Repository].
func (r *repository) GetById(ctx context.Context, id string) (domain.Visit, error) {
	return repo.GetByID[domain.Visit](ctx, r.db, id, "visit")
}

// ListByStudentId implements [Repository].
func (r *repository) ListByStudentId(ctx context.Context, studentId string) ([]domain.Visit, error) {
	var visits []domain.Visit

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Model(&domain.Visit{}).
		Where("student_id = ?", studentId).
		Find(&visits).
		Error
	if err != nil {
		return nil, err
	}

	return visits, nil
}

// ListByGroupId implements [Repository].
func (r *repository) ListByGroupId(ctx context.Context, groupId string) ([]domain.Visit, error) {
	var visits []domain.Visit

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Model(&domain.Visit{}).
		Joins("JOIN students ON students.id = visits.student_id").
		Joins("JOIN groups ON groups.id = students.group_id").
		Where("groups.id = ?", groupId).
		Find(&visits).
		Error
	if err != nil {
		return nil, err
	}

	return visits, err
}

// ListByOrganizationId implements [Repository].
func (r *repository) ListByOrganizationId(ctx context.Context, organizationId string) ([]domain.Visit, error) {
	var visits []domain.Visit

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Model(&domain.Visit{}).
		Joins("JOIN students ON students.id = visits.student_id").
		Joins("JOIN groups ON groups.id = students.group_id").
		Where("groups.organization_id = ?", organizationId).
		Find(&visits).
		Error
	if err != nil {
		return nil, err
	}

	return visits, err
}

// Create implements [Repository].
func (r *repository) Create(ctx context.Context, visit *domain.Visit) error {
	return repo.CREATE(ctx, r.db, visit, "visit")
}
