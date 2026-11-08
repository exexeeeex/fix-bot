package statistic

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/exexeeeex/fiks-api/internal/domain"
)

type testVisitService struct {
	visits []domain.Visit
}

func (s testVisitService) ListByGroupId(context.Context, string) ([]domain.Visit, error) {
	return nil, nil
}

func (s testVisitService) ListByOrganizationId(context.Context, string) ([]domain.Visit, error) {
	return s.visits, nil
}

type testStudentService struct {
	students []domain.Student
}

func (s testStudentService) GetByGroupId(context.Context, string) ([]domain.Student, error) {
	return nil, nil
}

func (s testStudentService) List(context.Context, string) ([]domain.Student, error) {
	return s.students, nil
}

type testStartTimeRepository struct {
	startTimes []domain.StartTime
	err        error
}

func (r testStartTimeRepository) ListByGroupIds(
	context.Context,
	[]string,
	time.Time,
	time.Time,
) ([]domain.StartTime, error) {
	return r.startTimes, r.err
}

func TestGetOrganizationStatisticsSeparatesAbsencesAndLates(t *testing.T) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	start := today.Add(8*time.Hour + 30*time.Minute)

	service := NewService(
		testVisitService{visits: []domain.Visit{
			{StudentId: "on-time", CreatedAt: start},
			{StudentId: "late", CreatedAt: start.Add(time.Minute)},
		}},
		testStudentService{students: []domain.Student{
			{Id: "on-time", GroupId: "group-1"},
			{Id: "late", GroupId: "group-1"},
			{Id: "missing", GroupId: "group-1"},
			{Id: "unscheduled", GroupId: "group-2"},
		}},
		testStartTimeRepository{startTimes: []domain.StartTime{
			{GroupId: "group-1", StartTime: start},
		}},
	)

	result, err := service.GetOrganizationStatistics(context.Background(), "organization-1")
	if err != nil {
		t.Fatalf("GetOrganizationStatistics() error = %v", err)
	}

	absent := make(map[string]bool, len(result.AbsencesStudents))
	for _, student := range result.AbsencesStudents {
		absent[student.Id] = true
	}
	late := make(map[string]bool, len(result.LateStudents))
	for _, student := range result.LateStudents {
		late[student.Id] = true
	}

	if absent["on-time"] || late["on-time"] {
		t.Error("on-time student was classified as absent or late")
	}
	if absent["late"] || !late["late"] {
		t.Error("late student was not classified as late only")
	}
	if !absent["missing"] {
		t.Error("student without a visit was not marked absent")
	}
	if absent["unscheduled"] || late["unscheduled"] {
		t.Error("student in a group without a start time was marked absent")
	}
}

func TestGetOrganizationAbsencesReturnsStartTimeError(t *testing.T) {
	startTimeErr := errors.New("database unavailable")
	service := NewService(
		testVisitService{},
		testStudentService{},
		testStartTimeRepository{err: startTimeErr},
	)

	_, err := service.GetOrganizationStatistics(context.Background(), "organization-1")
	if !errors.Is(err, startTimeErr) {
		t.Fatalf("GetOrganizationStatistics() error = %v, want %v", err, startTimeErr)
	}
}

func TestGetOrganizationLatesReturnsOnlyLateStudents(t *testing.T) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	start := today.Add(8*time.Hour + 30*time.Minute)
	service := NewService(
		testVisitService{visits: []domain.Visit{
			{StudentId: "late", CreatedAt: start.Add(time.Minute)},
		}},
		testStudentService{students: []domain.Student{
			{Id: "late", GroupId: "group-1"},
			{Id: "absent", GroupId: "group-1"},
		}},
		testStartTimeRepository{startTimes: []domain.StartTime{
			{GroupId: "group-1", StartTime: start},
		}},
	)

	lateStudents, err := service.GetOrganizationLates(context.Background(), "organization-1")
	if err != nil {
		t.Fatalf("GetOrganizationLates() error = %v", err)
	}
	if len(lateStudents) != 1 || lateStudents[0].Id != "late" {
		t.Fatalf("GetOrganizationLates() = %v, want only late student", lateStudents)
	}
}
