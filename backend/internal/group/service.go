package student

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
	"uuid"

	"github.com/exexeeeex/fiks-api/internal/domain"
	encrypt "github.com/exexeeeex/fiks-api/internal/encrypt"
)

type StudentCache interface {
	Get(
		ctx context.Context,
		key string,
	) (domain.Student, bool, error)

	Set(
		ctx context.Context,
		key string,
		student domain.Student,
		ttl int,
	) error
}

type StudentListCache interface {
	Get(
		ctx context.Context,
		key string,
	) ([]domain.Student, bool, error)

	Set(
		ctx context.Context,
		key string,
		students []domain.Student,
		ttl int,
	) error
}

type GroupService interface {
	GetById(ctx context.Context, id string) (domain.Group, error)
}

type Service struct {
	repo         Repository
	cache        StudentCache
	listCache    StudentListCache
	groupService GroupService
}

func NewService(
	repo Repository,
	cache StudentCache,
	listCache StudentListCache,
	groupService GroupService,
) *Service {
	return &Service{
		repo:         repo,
		cache:        cache,
		listCache:    listCache,
		groupService: groupService,
	}
}

func (s *Service) GenerateStudentKey(ctx context.Context) (string, error) {
	key := uuid.New().String()

	encrypted, err := encrypt.Encrypt(key, os.Getenv("USER_KEY"))
	if err != nil {
		return "", fmt.Errorf("error encrypt key: %w", err)
	}

	result := base64.StdEncoding.EncodeToString(encrypted)

	return result, nil
}

func (s *Service) GetByGroupId(
	ctx context.Context,
	groupId string,
) ([]domain.Student, error) {
	if groupId == "" {
		return nil, errors.New("group_id cannot be empty")
	}

	if s.listCache != nil {
		cached, found, err := s.listCache.Get(ctx, groupId)
		if err == nil && found {
			return cached, nil
		}
	}

	result, err := s.repo.GetByGroupId(ctx, groupId)
	if err != nil {
		return nil, fmt.Errorf("get students: %w", err)
	}

	if s.listCache != nil {
		_ = s.listCache.Set(ctx, groupId, result, 10)
	}

	return result, nil
}

func (s *Service) GetById(ctx context.Context, id string) (domain.Student, error) {
	if id == "" {
		return domain.Student{}, errors.New("id cannot be empty")
	}

	if s.cache != nil {
		cached, found, err := s.cache.Get(ctx, id)
		if err == nil && found {
			return cached, nil
		}
	}

	result, err := s.repo.GetById(ctx, id)
	if err != nil {
		return domain.Student{}, fmt.Errorf("get student: %w", err)
	}

	if s.cache != nil {
		_ = s.cache.Set(ctx, id, result, 10)
	}

	return result, nil
}

func (s *Service) Create(ctx context.Context, student domain.Student) (domain.Student, error) {
	student.Name = strings.TrimSpace(student.Name)

	if student.GroupId == "" {
		return domain.Student{}, errors.New("group_id cannot be empty")
	}

	if student.Name == "" {
		return domain.Student{}, errors.New("student Name cannot be empty")
	}

	_, err := s.groupService.GetById(ctx, student.GroupId)
	if err != nil {
		return domain.Student{}, fmt.Errorf("validate student group: %w", err)
	}

	student.Id = uuid.NewV4().String()

	if err := s.repo.Create(ctx, &student); err != nil {
		return domain.Student{}, fmt.Errorf("create student: %w", err)
	}

	if s.cache != nil {
		_ = s.cache.Set(ctx, student.Id, student, 10)
	}

	return student, nil
}
