package group

import (
	"context"
	"errors"
	"fmt"

	"github.com/exexeeeex/fiks-api/internal/domain"
)

type Cache interface {
	Get(ctx context.Context, key string) (domain.Group, bool, error)
	GetList(ctx context.Context, key string) ([]domain.Group, bool, error)
	Set(ctx context.Context, key string, group domain.Group) error
	SetList(ctx context.Context, key string, groups []domain.Group) error
}

type Service struct {
	repo  Repository
	cache Cache
}

func NewService(repo Repository, cache Cache) *Service {
	return &Service{
		repo:  repo,
		cache: cache,
	}
}

func (s *Service) GetByName(
	ctx context.Context, 
	name string,
	organizationId string,
) ([]domain.Group, error) {
	if name == "" {
		return []domain.Group{}, errors.New("name cannot be empty")
	}

	if s.cache != nil {
		cached, found, err := s.cache.GetList(ctx, name)
		if err == nil && found {
			return cached, nil
		}
	}

	result, err := s.repo.GetByName(ctx, name, organizationId)
	if err != nil {
		return  []domain.Group{}, fmt.Errorf("get group by name: %w", err)
	}

	if s.cache != nil {
		_ = s.cache.SetList(ctx, name, result)
	}

	return result, nil
}

func (s *Service) GetById(ctx context.Context, id string) (domain.Group, error) {
	if id == "" {
		return domain.Group{}, errors.New("id cannot be empty")
	}

	if s.cache != nil {
		cached, found, err := s.cache.Get(ctx, id)
		if err == nil && found {
			return cached, nil
		}
	}

	result, err := s.repo.GetById(ctx, id)
	if err != nil {
		return domain.Group{}, fmt.Errorf("get group by id: %w", err)
	}

	if s.cache != nil {
		_ = s.cache.Set(ctx, id, result)
	}

	return result, nil
}

func (s *Service) List(ctx context.Context, organizationId string) ([]domain.Group, error) {
	if s.cache != nil {
		cached, found, err := s.cache.GetList(ctx, "all")
		if err == nil && found {
			return cached, nil
		}
	}

	result, err := s.repo.List(ctx, organizationId)
	if err != nil {
		return []domain.Group{}, fmt.Errorf("list group: %w", err)
	}

	if s.cache != nil {
		_ = s.cache.SetList(ctx, "all", result)
	}
	return result, nil
}
