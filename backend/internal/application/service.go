package application

import (
	"context"
	"errors"
	"fmt"
	"uuid"

	"github.com/exexeeeex/fiks-api/internal/domain"
	"gorm.io/gorm"
)

type ApplicationsCache interface {
	Get(
		ctx context.Context,
		key string,
	) (domain.Application, bool, error)

	Set(
		ctx context.Context,
		key string,
		application domain.Application,
		ttl int,
	) error
}

type ApplicationsListCache interface {
	Get(
		ctx context.Context,
		key string,
	) ([]domain.Application, bool, error)

	Set(
		ctx context.Context,
		key string,
		applications []domain.Application,
		ttl int,
	) error
}

type GroupService interface {
	GetById(ctx context.Context, id string) (domain.Group, error)
}

type Service struct {
	repo         Repository
	cache        ApplicationsCache
	listCache    ApplicationsListCache
	groupService GroupService
}

func NewService(
	repo Repository,
	cache ApplicationsCache,
	listCahce ApplicationsListCache,
	groupService GroupService,
) *Service {
	return &Service{
		repo:         repo,
		cache:        cache,
		listCache:    listCahce,
		groupService: groupService,
	}
}

func (s *Service) GetByUserId(ctx context.Context, userId string) (domain.Application, error) {
	if userId == "" {
		return domain.Application{}, fmt.Errorf("user id cannot be empty")
	}

	if s.cache != nil {
		cached, found, err := s.cache.Get(ctx, "user-id:"+userId)
		if err == nil && found {
			return cached, nil
		}
	}

	result, err := s.repo.GetByUserId(ctx, userId)
	if err != nil {
		return domain.Application{}, fmt.Errorf("get by user id: %w", err)
	}

	if s.cache != nil {
		_ = s.cache.Set(ctx, "user-id:"+userId, result, 10)
	}

	return result, nil
}

func (s *Service) GetById(ctx context.Context, id string) (domain.Application, error) {
	if id == "" {
		return domain.Application{}, fmt.Errorf("id cannot be empty")
	}

	if s.cache != nil {
		cached, found, err := s.cache.Get(ctx, id)
		if err == nil && found {
			return cached, nil
		}
	}

	result, err := s.repo.GetById(ctx, id)
	if err != nil {
		return domain.Application{}, fmt.Errorf("get application by id: %w", err)
	}

	if s.cache != nil {
		_ = s.cache.Set(ctx, id, result, 10)
	}

	return result, nil
}

func (s *Service) GetByGroupId(ctx context.Context, groupId string) ([]domain.Application, error) {
	if groupId == "" {
		return []domain.Application{}, fmt.Errorf("group_id cannot be empty")
	}

	if s.cache != nil {
		cached, found, err := s.listCache.Get(ctx, "group:"+groupId)
		if err == nil && found {
			return cached, nil
		}
	}

	result, err := s.repo.GetByGroupId(ctx, groupId)
	if err != nil {
		return []domain.Application{}, fmt.Errorf("get by group id: %w", err)
	}

	if s.listCache != nil {
		_ = s.listCache.Set(ctx, "group:"+groupId, result, 10)
	}

	return result, nil
}

func (s *Service) List(ctx context.Context, organizationId string) ([]domain.Application, error) {
	if organizationId == "" {
		return []domain.Application{}, fmt.Errorf("organization id cannot be empty")
	}

	if s.listCache != nil {
		cached, found, err := s.listCache.Get(ctx, "all")
		if err == nil && found {
			return cached, nil
		}
	}

	result, err := s.repo.List(ctx, organizationId)
	if err != nil {
		return []domain.Application{}, fmt.Errorf("get: %w", err)
	}

	if s.listCache != nil {
		_ = s.listCache.Set(ctx, "all", result, 10)
	}

	return result, nil
}

func (s *Service) Create(ctx context.Context, application domain.Application) (domain.Application, error) {
	if len(application.OrganizationId) < 5 {
		return domain.Application{}, fmt.Errorf("organization id is invalid")
	}
	if len(application.Name) < 5 {
		return domain.Application{}, fmt.Errorf("name is invalid")
	}
	if len(application.UserId) < 5 {
		return domain.Application{}, fmt.Errorf("user id is invalid")
	}
	if len(application.GroupId) < 5 {
		return domain.Application{}, fmt.Errorf("group id is invalid")
	}

	_, err := s.GetByUserId(ctx, application.UserId)

	switch {
	case err == nil:
		return domain.Application{}, fmt.Errorf(
			"already exist %s",
			application.UserId,
		)

	case errors.Is(err, gorm.ErrRecordNotFound):
		_, err = s.groupService.GetById(ctx, application.GroupId)
		if err != nil {
			return domain.Application{}, fmt.Errorf("validate group: %w", err)
		}

		application.Id = uuid.New().String()

		if err := s.repo.Create(ctx, &application); err != nil {
			return domain.Application{}, fmt.Errorf("create applicatiom: %w", err)
		}

		if s.cache != nil {
			_ = s.cache.Set(ctx, application.Id, application, 10)
		}

		return application, nil
	default:
		return domain.Application{}, fmt.Errorf(
			"check existing application: %w",
			err,
		)
	}

}
