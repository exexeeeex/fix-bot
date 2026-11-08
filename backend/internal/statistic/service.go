package statistic

import (
	"context"
	"fmt"
	"time"

	"github.com/exexeeeex/fiks-api/internal/domain"
)

type VisitService interface {
	ListByGroupId(ctx context.Context, groupId string) ([]domain.Visit, error)
	ListByOrganizationId(ctx context.Context, organizationId string) ([]domain.Visit, error)
}

type StudentService interface {
	GetByGroupId(ctx context.Context, groupId string) ([]domain.Student, error)
	List(ctx context.Context, organizationId string) ([]domain.Student, error)
}

type Service struct {
	visitService       VisitService
	studentService     StudentService
	startTimeRepository StartTimeRepository
}

func NewService(
	visitService VisitService,
	studentService StudentService,
	startTimeRepository ...StartTimeRepository,
) *Service {
	var repo StartTimeRepository
	if len(startTimeRepository) > 0 {
		repo = startTimeRepository[0]
	}

	return &Service{
		visitService:        visitService,
		studentService:      studentService,
		startTimeRepository: repo,
	}
}

func (s *Service) GetOrganizationAbsences(
	ctx context.Context,
	organizationId string,
) (Absence, error) {
	result, err := s.GetOrganizationStatistics(ctx, organizationId)
	if err != nil {
		return Absence{}, err
	}

	return Absence{
		AllStudents:       result.AllStudents,
		AbsencesStudents: result.AbsencesStudents,
	}, nil
}

func (s *Service) GetOrganizationLates(
	ctx context.Context,
	organizationId string,
) ([]*domain.Student, error) {
	result, err := s.GetOrganizationStatistics(ctx, organizationId)
	if err != nil {
		return nil, err
	}

	return result.LateStudents, nil
}

func (s *Service) GetOrganizationStatistics(
	ctx context.Context,
	organizationId string,
) (Statistics, error) {
	students, err := s.studentService.List(ctx, organizationId)
	if err != nil {
		return Statistics{}, fmt.Errorf("get students by organization id: %w", err)
	}

	visits, err := s.visitService.ListByOrganizationId(ctx, organizationId)
	if err != nil {
		return Statistics{}, fmt.Errorf("get visits by organization id: %w", err)
	}

	result := Statistics{
		AllStudents: students,
	}

	if s.startTimeRepository == nil {
		return s.getLegacyOrganizationStatistics(result, visits), nil
	}

	stats, err := s.getScheduledOrganizationStatistics(ctx, students, visits)
	if err != nil {
		return Statistics{}, err
	}

	return stats, nil
}

func (s *Service) getLegacyOrganizationStatistics(
	result Statistics,
	visits []domain.Visit,
) Statistics {
	todayVisitsIds := make(map[string]struct{})
	now := time.Now()
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, now.Location())

	for _, visit := range visits {
		visitDate := visit.CreatedAt.In(now.Location())
		if visitDate.Year() == today.Year() &&
			visitDate.Month() == today.Month() &&
			visitDate.Day() == today.Day() {
			todayVisitsIds[visit.StudentId] = struct{}{}
		}
	}

	for i := range result.AllStudents {
		student := &result.AllStudents[i]
		if _, visited := todayVisitsIds[student.Id]; !visited {
			result.AbsencesStudents = append(result.AbsencesStudents, student)
		}
	}

	return result
}

func (s *Service) getScheduledOrganizationStatistics(
	ctx context.Context,
	students []domain.Student,
	visits []domain.Visit,
) (Statistics, error) {
	result := Statistics{AllStudents: students}
	groupIds := make([]string, 0, len(students))
	seenGroups := make(map[string]struct{}, len(students))
	for _, student := range students {
		if student.GroupId == "" {
			continue
		}
		if _, exists := seenGroups[student.GroupId]; !exists {
			seenGroups[student.GroupId] = struct{}{}
			groupIds = append(groupIds, student.GroupId)
		}
	}

	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endDate := startDate.Add(24 * time.Hour)
	startTimes, err := s.startTimeRepository.ListByGroupIds(ctx, groupIds, startDate, endDate)
	if err != nil {
		return Statistics{}, err
	}

	startTimeByGroup := make(map[string]time.Time, len(startTimes))
	for _, startTime := range startTimes {
		if _, exists := startTimeByGroup[startTime.GroupId]; !exists || startTime.StartTime.Before(startTimeByGroup[startTime.GroupId]) {
			startTimeByGroup[startTime.GroupId] = startTime.StartTime
		}
	}

	todayVisitsByStudent := make(map[string]time.Time)
	for _, visit := range visits {
		visitDate := visit.CreatedAt.In(now.Location())
		if visitDate.Year() == startDate.Year() &&
			visitDate.Month() == startDate.Month() &&
			visitDate.Day() == startDate.Day() {
			todayVisitsByStudent[visit.StudentId] = visit.CreatedAt
		}
	}

	for i := range students {
		student := &students[i]
		startTime, hasStartTime := startTimeByGroup[student.GroupId]
		if !hasStartTime {
			continue
		}

		visitAt, visited := todayVisitsByStudent[student.Id]
		if !visited {
			result.AbsencesStudents = append(result.AbsencesStudents, student)
			continue
		}

		if visitAt.After(startTime) {
			result.LateStudents = append(result.LateStudents, student)
		}
	}

	return result, nil
}
