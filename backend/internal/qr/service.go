package qr

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/exexeeeex/fiks-api/internal/domain"
	"github.com/exexeeeex/fiks-api/internal/encrypt"
	"github.com/skip2/go-qrcode"
)

type StudentService interface {
	GetById(ctx context.Context, id string) (domain.Student, error)
	GetByUserId(ctx context.Context, userId string) (domain.Student, error)
	GenerateStudentKey(ctx context.Context, studentId string) (string, error)
	ValidateStudentKey(ctx context.Context, key string) (bool, error)
}

type GroupService interface {
	GetById(ctx context.Context, id string) (domain.Group, error)
}

type VisitService interface {
	Create(ctx context.Context, visit domain.Visit) (domain.Visit, error)
}

type Service struct {
	studentService StudentService
	groupService   GroupService
	visitService   VisitService
}

func NewService(
	studentService StudentService,
	groupService GroupService,
	visitService VisitService,
) *Service {
	return &Service{
		studentService: studentService,
		groupService:   groupService,
		visitService:   visitService,
	}
}

func (s *Service) ReadQrCode(ctx context.Context, data string) (ReadQrResponse, error) {
	if data == "" {
		return ReadQrResponse{}, fmt.Errorf("data cannot be empty")
	}

	encrypted, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return ReadQrResponse{}, fmt.Errorf("decode data: %w", err)
	}

	decrypted, err := encrypt.Decrypt(encrypted, os.Getenv("QR_CODE_KEY"))
	if err != nil {
		return ReadQrResponse{}, fmt.Errorf("decrypt data: %w", err)
	}

	result := ReadQrResponse{}

	if err := json.Unmarshal([]byte(decrypted), &result); err != nil {
		return ReadQrResponse{}, fmt.Errorf("unmarshal: %w", err)
	}

	validated, err := s.studentService.ValidateStudentKey(ctx, result.Key)
	if err != nil {
		return ReadQrResponse{}, fmt.Errorf("validation key: %w", err)
	}
	if !validated {
		return ReadQrResponse{}, fmt.Errorf("student key is invalid")
	}

	if _, err := s.visitService.Create(ctx, domain.Visit{StudentId: result.Id}); err != nil {
		return ReadQrResponse{}, fmt.Errorf("create visit: %w", err)
	}

	return result, nil
}

func (s *Service) GenerateQrCode(ctx context.Context, studentId string) ([]byte, error) {
	if studentId == "" {
		return nil, fmt.Errorf("student_id cannot be empty")
	}

	student, err := s.studentService.GetById(ctx, studentId)
	if err != nil {
		return []byte("Error"), fmt.Errorf("get student by id: %w", err)
	}

	group, err := s.groupService.GetById(ctx, student.GroupId)
	if err != nil {
		return []byte("Error"), fmt.Errorf("get group by id: %w", err)
	}

	key, err := s.studentService.GenerateStudentKey(ctx, studentId)
	if err != nil {
		return []byte("Error"), fmt.Errorf("generate key: %w", err)
	}

	model := ReadQrResponse{
		Id:    student.Id,
		Group: group,
		Key:   key,
	}

	data, err := json.Marshal(model)
	if err != nil {
		return []byte("Error"), fmt.Errorf("marshal student model: %w", err)
	}

	encrypted, err := encrypt.Encrypt(string(data), os.Getenv("QR_CODE_KEY"))
	if err != nil {
		return []byte(""), fmt.Errorf("encrypt student model: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(encrypted)

	png, err := qrcode.Encode(encoded, qrcode.Medium, 256)
	if err != nil {
		return []byte(""), fmt.Errorf("failed to generate student QR code: %w", err)
	}

	return png, nil
}

func (s *Service) GenerateQrCodeForUser(ctx context.Context, userId string) ([]byte, error) {
	student, err := s.studentService.GetByUserId(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("get student for user: %w", err)
	}

	return s.GenerateQrCode(ctx, student.Id)
}
