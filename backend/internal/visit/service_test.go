package visit

import (
	"context"
	"testing"
	"time"

	"github.com/exexeeeex/fiks-api/internal/domain"
)

type testRepository struct {
	visits      []domain.Visit
	createCalls int
}

func (r *testRepository) GetById(context.Context, string) (domain.Visit, error) {
	return domain.Visit{}, nil
}

func (r *testRepository) ListByStudentId(context.Context, string) ([]domain.Visit, error) {
	return r.visits, nil
}

func (r *testRepository) ListByGroupId(context.Context, string) ([]domain.Visit, error) {
	return nil, nil
}

func (r *testRepository) ListByOrganizationId(context.Context, string) ([]domain.Visit, error) {
	return nil, nil
}

func (r *testRepository) Create(_ context.Context, visit *domain.Visit) error {
	r.createCalls++
	if visit.CreatedAt.IsZero() {
		visit.CreatedAt = time.Now()
	}
	r.visits = append([]domain.Visit{*visit}, r.visits...)
	return nil
}

func TestCreateFirstVisitAndAvoidDuplicateOnSameDay(t *testing.T) {
	repository := &testRepository{}
	service := NewService(repository)

	firstVisit, err := service.Create(context.Background(), domain.Visit{StudentId: "student-123"})
	if err != nil {
		t.Fatalf("Create() first visit error = %v", err)
	}
	if firstVisit.Id == "" {
		t.Fatal("Create() first visit returned an empty ID")
	}
	if repository.createCalls != 1 {
		t.Fatalf("create calls = %d, want 1", repository.createCalls)
	}

	secondVisit, err := service.Create(context.Background(), domain.Visit{StudentId: "student-123"})
	if err != nil {
		t.Fatalf("Create() same-day visit error = %v", err)
	}
	if secondVisit.Id != firstVisit.Id {
		t.Errorf("same-day visit ID = %q, want existing ID %q", secondVisit.Id, firstVisit.Id)
	}
	if repository.createCalls != 1 {
		t.Errorf("create calls after same-day visit = %d, want 1", repository.createCalls)
	}
}
