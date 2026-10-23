package student

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"uuid"

	"github.com/exexeeeex/fiks-api/internal/domain"
	encrypt "github.com/exexeeeex/fiks-api/internal/encrypt"
	"github.com/skip2/go-qrcode"
)

type Cache interface {
	Get(ctx context.Context, key string) (domain.Student, bool, error)
	GetList(ctx context.Context, key string) ([]domain.Student, bool, error)
	Set(ctx context.Context, key string, student domain.Student) error
	SetList(ctx context.Context, key string, students []domain.Student) error
}

type GroupService interface {
	GetById(ctx context.Context, id string) (domain.Group, error)
}

type Service struct {
	repo         Repository
	cache        Cache
	groupService GroupService
}

func NewService(
	repo Repository,
	cache Cache,
	groupService GroupService,
) *Service {
	return &Service{
		repo:         repo,
		cache:        cache,
		groupService: groupService,
	}
}

func (s *Service) GenerateStudentQrCode(ctx context.Context, studentId string) ([]byte, error) {
	if studentId == "" {
		return nil, fmt.Errorf("student_id cannot be empty")
	}

	request, err := s.GetById(ctx, studentId)
	if err != nil {
		return []byte("error"), fmt.Errorf("get student by id: %w", err)
	}

	student := domain.Student{
		Id:      studentId,
		Name:    request.Name,
		GroupId: request.GroupId,
	}

	data, err := json.Marshal(student)
	if err != nil {
		return []byte("error"), fmt.Errorf("marshal student: %w", err)
	}

	encrypted, err := encrypt.Encrypt(string(data), "studentId")
	if err != nil {
		return []byte("error"), fmt.Errorf("encrypt error: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(encrypted)

	png, err := qrcode.Encode(encoded, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate student QR code: %w", err)
	}

	return png, nil
}

func (s *Service) GetByGroupId(ctx context.Context, groupId string) ([]domain.Student, error) {
	if groupId == "" {
		return []domain.Student{}, errors.New("groupd_id cannot be empty")
	}

	if s.cache != nil {
		cached, found, err := s.cache.GetList(ctx, "students")
		if err == nil && found {
			return cached, nil
		}
	}

	result, err := s.repo.GetByGroupId(ctx, groupId)
	if err != nil {
		return []domain.Student{}, fmt.Errorf("get students: %w", err)
	}

	if s.cache != nil {
		_ = s.cache.SetList(ctx, "all", result)
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
		_ = s.cache.Set(ctx, id, result)
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
		_ = s.cache.Set(ctx, student.Id, student)
	}

	return student, nil
}
