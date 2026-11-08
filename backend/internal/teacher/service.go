package teacher

import (
	"context"
	"fmt"
	"strings"
	"uuid"

	"github.com/exexeeeex/fiks-api/internal/domain"
)

type TeacherCache interface {
	Get(
		ctx context.Context,
		key string,
	) (domain.Teacher, bool, error)

	Set(
		ctx context.Context,
		key string,
		teacher domain.Teacher,
		ttl int,
	) error
}

type TeacherListCache interface {
	Get(
		ctx context.Context,
		key string,
	) ([]domain.Teacher, bool, error)

	Set(
		ctx context.Context,
		key string,
		teachers []domain.Teacher,
		ttl int,
	) error
}

type OrganizationService interface {
	GetById(ctx context.Context, id string) (domain.Organization, error)
}

type Service struct {
	repo                Repository
	cache               TeacherCache
	listCache           TeacherListCache
	organizationService OrganizationService
}

func NewService(
	repo Repository,
	cache TeacherCache,
	listCache TeacherListCache,
	organizationService OrganizationService,
) *Service {
	return &Service{
		repo:                repo,
		cache:               cache,
		listCache:           listCache,
		organizationService: organizationService,
	}
}

func (s *Service) Create(ctx context.Context, request domain.Teacher) (domain.Teacher, error) {
	request.Name = strings.TrimSpace(request.Name)

	if request.OrganizationId == "" {
		return domain.Teacher{}, fmt.Errorf("organization id cannot be empty")
	}
	if request.Name == "" {
		return domain.Teacher{}, fmt.Errorf("teacher name cannot be empty")
	}

	organization, err := s.organizationService.GetById(ctx, request.OrganizationId)
	if err != nil {
		return domain.Teacher{}, fmt.Errorf("get organization: %w", err)
	}

	teacher := domain.Teacher{
		Id:             uuid.NewV4().String(),
		Name:           request.Name,
		OrganizationId: organization.Id,
		UserId:         request.UserId,
	}

	if s.cache != nil {
		_ = s.cache.Set(ctx, teacher.Id, teacher, 10)
	}

	err = s.repo.Create(ctx, &teacher)
	if err != nil {
		return domain.Teacher{}, fmt.Errorf("create teacher: %w", err)
	}

	return teacher, nil
}

func (s *Service) GetByUserId(ctx context.Context, userId string) (domain.Teacher, error) {
	if userId == "" {
		return domain.Teacher{}, fmt.Errorf("user_id cannot be empty")
	}

	teacher, err := s.repo.GetByUserId(ctx, userId)
	if err != nil {
		return domain.Teacher{}, fmt.Errorf("get teacher by user_id: %w", err)
	}

	if s.cache != nil {
		_ = s.cache.Set(ctx, teacher.Id, teacher, 10)
	}
	return teacher, nil
}
