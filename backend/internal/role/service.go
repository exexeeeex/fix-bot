package role

import (
	"context"
	"fmt"

	"github.com/exexeeeex/fiks-api/internal/domain"
)

type RoleCache interface {
	Get(
		ctx context.Context,
		key string,
	) (domain.Role, bool, error)
	Set(
		ctx context.Context,
		key string,
		role domain.Role,
		ttl int,
	) error
}

type RoleListCache interface {
	Get(
		ctx context.Context,
		key string,
	) ([]domain.Role, bool, error)
	Set(
		ctx context.Context,
		key string,
		roles []domain.Role,
		ttl int,
	) error
}

type Service struct {
	repo      Repository
	cache     RoleCache
	listCache RoleListCache
}

func NewService(
	repo Repository,
	cache RoleCache,
	listCache RoleListCache,
) *Service {
	return &Service{
		repo:      repo,
		cache:     cache,
		listCache: listCache,
	}
}

func (s *Service) List(ctx context.Context) ([]domain.Role, error) {
	if s.listCache != nil {
		roles, found, err := s.listCache.Get(ctx, "all")
		if err == nil && found {
			return roles, nil
		}
	}

	roles, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list roles: %w", err)
	}

	if s.listCache != nil {
		_ = s.listCache.Set(ctx, "all", roles, 10)
	}

	return roles, nil
}

func (s *Service) GetById(ctx context.Context, id string) (domain.Role, error) {
	if id == "" {
		return domain.Role{}, fmt.Errorf("id cannot be empty")
	}

	if s.cache != nil {
		role, found, err := s.cache.Get(ctx, id)
		if err == nil && found {
			return role, nil
		}
	}

	role, err := s.repo.GetById(ctx, id)
	if err != nil {
		return domain.Role{}, fmt.Errorf("get role by id: %w", err)
	}

	if s.cache != nil {
		_ = s.cache.Set(ctx, id, role, 10)
	}

	return role, nil
}

func (s *Service) GetByName(ctx context.Context, name string) (domain.Role, error) {
	if name == "" {
		return domain.Role{}, fmt.Errorf("name cannot be empty")
	}

	if s.cache != nil {
		role, found, err := s.cache.Get(ctx, name)
		if err == nil && found {
			return role, nil
		}
	}

	role, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return domain.Role{}, fmt.Errorf("get role by name: %w", err)
	}

	if s.cache != nil {
		_ = s.cache.Set(ctx, name, role, 10)
	}

	return role, nil
}
