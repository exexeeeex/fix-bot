package visit

import (
	"context"
	"fmt"
	"time"
	"uuid"

	"github.com/exexeeeex/fiks-api/internal/domain"
)

type Service struct {
	repo Repository
}

func NewService(
	repo Repository,
) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, visit domain.Visit) (domain.Visit, error) {
	if visit.StudentId == "" {
		return domain.Visit{}, fmt.Errorf("student id cannot be empty")
	}

	visits, err := s.repo.ListByStudentId(ctx, visit.StudentId)
	if err != nil {
		return domain.Visit{}, fmt.Errorf("list by student id: %w", err)
	}

	if len(visits) == 0 {
		return s.create(ctx, visit.StudentId)
	}

	lastVisit := visits[0]
	now := time.Now().In(lastVisit.CreatedAt.Location())
	lastVisitDate := lastVisit.CreatedAt.In(now.Location())

	if now.Year() == lastVisitDate.Year() &&
		now.YearDay() == lastVisitDate.YearDay() {
		return lastVisit, nil
	}

	return s.create(ctx, visit.StudentId)
}

func (s *Service) create(ctx context.Context, studentID string) (domain.Visit, error) {
	visit := domain.Visit{
		Id:        uuid.NewV4().String(),
		StudentId: studentID,
	}

	if err := s.repo.Create(ctx, &visit); err != nil {
		return domain.Visit{}, fmt.Errorf("create: %w", err)
	}

	return visit, nil
}

func (s *Service) ListByGroupId(ctx context.Context, groupId string) ([]domain.Visit, error) {
	visits, err := s.repo.ListByGroupId(ctx, groupId)
	if err != nil {
		return nil, fmt.Errorf("list visits by groupd id: %w", err)
	}

	return visits, nil
}

func (s *Service) ListByOrganizationId(ctx context.Context, organizationId string) ([]domain.Visit, error) {
	visits, err := s.repo.ListByOrganizationId(ctx, organizationId)
	if err != nil {
		return nil, fmt.Errorf("list visits by organization id: %w", err)
	}

	return visits, nil
}
