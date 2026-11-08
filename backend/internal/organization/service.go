package organization

import (
	"context"
	"fmt"

	"github.com/exexeeeex/fiks-api/internal/domain"
)

type OrganizationCache interface {
	Get(
		ctx context.Context,
		key string,
	) (domain.Organization, bool, error)

	Set(
		ctx context.Context,
		key string,
		organization domain.Organization,
		ttl int,
	) error
}

type OrganizationListCache interface {
	Get(
		ctx context.Context,
		key string,
	) ([]domain.Organization, bool, error)

	Set(
		ctx context.Context,
		key string,
		organizations []domain.Organization,
		ttl int,
	) error
}

type Service struct {
	repo      Repository
	cache     OrganizationCache
	listCache OrganizationListCache
}

func NewService(
	repo Repository,
	cache OrganizationCache,
	listCache OrganizationListCache,
) *Service {
	return &Service{
		repo:      repo,
		cache:     cache,
		listCache: listCache,
	}
}

func (s *Service) GetById(ctx context.Context, id string) (domain.Organization, error) {
	if id == "" {
		return domain.Organization{}, fmt.Errorf("id cannot be empty")
	}

	if s.cache != nil {
		result, found, err := s.cache.Get(ctx, id)
		if found && err == nil {
			return result, nil
		}
	}

	result, err := s.repo.GetById(ctx, id)
	if err != nil {
		return domain.Organization{}, fmt.Errorf("get by id: %w", err)
	}

	if s.cache != nil {
		_ = s.cache.Set(ctx, id, result, 50)
	}

	return result, nil
}

func (s *Service) GetByName(ctx context.Context, name string) ([]domain.Organization, error) {
	if name == "" {
		return []domain.Organization{}, fmt.Errorf("name connot be empty")
	}
	if len(name) < 2 {
		return []domain.Organization{}, fmt.Errorf("name invalid")
	}

	if s.listCache != nil {
		result, found, err := s.listCache.Get(ctx, name)
		if found && err == nil {
			return result, nil
		}
	}

	result, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return []domain.Organization{}, fmt.Errorf("get by name: %w", err)
	}

	if s.listCache != nil {
		_ = s.listCache.Set(ctx, name, result, 10)
	}

	return result, nil
}

func (s *Service) List(ctx context.Context) ([]domain.Organization, error) {
	if s.listCache != nil {
		result, found, err := s.listCache.Get(ctx, "all")
		if found && err == nil {
			return result, nil
		}
	}

	result, err := s.repo.List(ctx)
	if err != nil {
		return []domain.Organization{}, fmt.Errorf("list organization: %w", err)
	}

	if s.listCache != nil {
		_ = s.listCache.Set(ctx, "all", result, 10)
	}

	return result, nil
}
