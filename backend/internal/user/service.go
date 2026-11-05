package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/exexeeeex/fiks-api/internal/domain"
	"gorm.io/gorm"
)

type RoleService interface {
	GetByName(ctx context.Context, name string) (domain.Role, error)
}

type UserCache interface {
	Get(
		ctx context.Context,
		key string,
	) (domain.User, bool, error)
	Set(
		ctx context.Context,
		key string,
		user domain.User,
		ttl int,
	) error
}

type UserListCache interface {
	Get(
		ctx context.Context,
		key string,
	) ([]domain.User, bool, error)
	Set(
		ctx context.Context,
		key string,
		students []domain.User,
		ttl int,
	) error
}

type Service struct {
	repo        Repository
	cache       UserCache
	listCache   UserListCache
	roleService RoleService
}

func NewService(
	repo Repository,
	cache UserCache,
	listCache UserListCache,
	roleService RoleService,
) *Service {
	return &Service{
		repo:        repo,
		cache:       cache,
		listCache:   listCache,
		roleService: roleService,
	}
}

func (s *Service) GetByMaxId(ctx context.Context, maxId string) (domain.User, error) {
	if maxId == "" {
		return domain.User{}, errors.New("max_id cannot be empty")
	}

	if s.cache != nil {
		cached, found, err := s.cache.Get(ctx, maxId)
		if err == nil && found {
			return cached, nil
		}
	}

	result, err := s.repo.GetByMaxId(ctx, maxId)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user by max_id: %w", err)
	}

	if s.cache != nil {
		_ = s.cache.Set(ctx, maxId, result, 10)
	}

	return result, nil
}

func (s *Service) GetById(ctx context.Context, id string) (domain.User, error) {
	if id == "" {
		return domain.User{}, errors.New("id cannot be empty")
	}

	if s.cache != nil {
		cached, found, err := s.cache.Get(ctx, id)
		if err == nil && found {
			return cached, nil
		}
	}

	result, err := s.repo.GetById(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user by id: %w", err)
	}

	if s.cache != nil {
		_ = s.cache.Set(ctx, id, result, 10)
	}

	return result, nil
}

func (s *Service) List(ctx context.Context, organizationId string) ([]domain.User, error) {
	if organizationId == "" {
		return nil, errors.New("organization_id cannot be empty")
	}

	if s.listCache != nil {
		cached, found, err := s.listCache.Get(ctx, organizationId)
		if err == nil && found {
			return cached, nil
		}
	}

	result, err := s.repo.List(ctx, organizationId)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	if s.listCache != nil {
		_ = s.listCache.Set(ctx, organizationId, result, 10)
	}

	return result, nil
}

func (s *Service) Create(ctx context.Context, user domain.User) (domain.User, error) {
	if user.MaxId == "" {
		return domain.User{}, errors.New("max_id cannot be empty")
	}

	existing, err := s.repo.GetByMaxId(ctx, user.MaxId)
	if err == nil {
		return existing, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.User{}, fmt.Errorf("get user by max_id: %w", err)
	}

	role, err := s.roleService.GetByName(ctx, "guest")
	if err != nil {
		return domain.User{}, fmt.Errorf("get role by name: %w", err)
	}

	user.RoleId = role.Id

	if err := s.repo.Create(ctx, &user); err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	if s.cache != nil {
		_ = s.cache.Set(ctx, user.Id, user, 10)
	}

	return user, nil
}
