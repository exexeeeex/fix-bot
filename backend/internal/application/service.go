package application

import (
	"context"
	"errors"
	"fmt"
	"uuid"

	"github.com/exexeeeex/fiks-api/internal/domain"
	"gorm.io/gorm"
)

type ApplicationsCache interface {
	Get(
		ctx context.Context,
		key string,
	) (domain.Application, bool, error)

	Set(
		ctx context.Context,
		key string,
		application domain.Application,
		ttl int,
	) error
}

type ApplicationsListCache interface {
	Get(
		ctx context.Context,
		key string,
	) ([]domain.Application, bool, error)

	Set(
		ctx context.Context,
		key string,
		applications []domain.Application,
		ttl int,
	) error
}

type GroupService interface {
	GetById(ctx context.Context, id string) (domain.Group, error)
}

type StudentService interface {
	Create(ctx context.Context, student domain.Student) (domain.Student, error)
	GetByUserId(ctx context.Context, userId string) (domain.Student, error)
}

type TeacherService interface {
	Create(ctx context.Context, teacher domain.Teacher) (domain.Teacher, error)
	GetByUserId(ctx context.Context, userId string) (domain.Teacher, error)
}

type RoleService interface {
	GetByName(ctx context.Context, name string) (domain.Role, error)
}

type UserService interface {
	Update(ctx context.Context, user domain.User) (domain.User, error)
	GetById(ctx context.Context, id string) (domain.User, error)
}

type Service struct {
	repo           Repository
	cache          ApplicationsCache
	listCache      ApplicationsListCache
	groupService   GroupService
	studentService StudentService
	teacherService TeacherService
	roleService    RoleService
	userService    UserService
}

func NewService(
	repo Repository,
	cache ApplicationsCache,
	listCache ApplicationsListCache,
	groupService GroupService,
	studentService StudentService,
	teacherService TeacherService,
	roleService RoleService,
	userService UserService,
) *Service {
	return &Service{
		repo:           repo,
		cache:          cache,
		listCache:      listCache,
		groupService:   groupService,
		studentService: studentService,
		teacherService: teacherService,
		roleService:    roleService,
		userService:    userService,
	}
}

func (s *Service) GetByMaxId(ctx context.Context, maxId string) (domain.Application, error) {
	if maxId == "" {
		return domain.Application{}, fmt.Errorf("max id cannot be empty")
	}

	if s.cache != nil {
		cached, found, err := s.cache.Get(ctx, "max-id:"+maxId)
		if err == nil && found {
			return cached, nil
		}
	}

	result, err := s.repo.GetByMaxId(ctx, maxId)
	if err != nil {
		return domain.Application{}, fmt.Errorf("get by max id: %w", err)
	}

	if s.cache != nil {
		_ = s.cache.Set(ctx, "max-id:"+maxId, result, 10)
	}

	return result, nil
}

func (s *Service) GetByUserId(ctx context.Context, userId string) (domain.Application, error) {
	if userId == "" {
		return domain.Application{}, fmt.Errorf("user id cannot be empty")
	}

	if s.cache != nil {
		cached, found, err := s.cache.Get(ctx, "user-id:"+userId)
		if err == nil && found {
			return cached, nil
		}
	}

	result, err := s.repo.GetByUserId(ctx, userId)
	if err != nil {
		return domain.Application{}, fmt.Errorf("get by user id: %w", err)
	}

	if s.cache != nil {
		_ = s.cache.Set(ctx, "user-id:"+userId, result, 10)
	}

	return result, nil
}

func (s *Service) GetById(ctx context.Context, id string) (domain.Application, error) {
	if id == "" {
		return domain.Application{}, fmt.Errorf("id cannot be empty")
	}

	if s.cache != nil {
		cached, found, err := s.cache.Get(ctx, id)
		if err == nil && found {
			return cached, nil
		}
	}

	result, err := s.repo.GetById(ctx, id)
	if err != nil {
		return domain.Application{}, fmt.Errorf("get application by id: %w", err)
	}

	if s.cache != nil {
		_ = s.cache.Set(ctx, id, result, 10)
	}

	return result, nil
}

func (s *Service) GetByGroupId(ctx context.Context, groupId string) ([]*domain.Application, error) {
	if groupId == "" {
		return []*domain.Application{}, fmt.Errorf("group_id cannot be empty")
	}

	if s.listCache != nil {
		cached, found, err := s.listCache.Get(ctx, "group:"+groupId)
		if err == nil && found {
			applications := make([]*domain.Application, 0, len(cached))
			for i := range cached {
				application := cached[i]
				applications = append(applications, &application)
			}
			return applications, nil
		}
	}

	result, err := s.repo.GetByGroupId(ctx, groupId)
	if err != nil {
		return []*domain.Application{}, fmt.Errorf("get by group id: %w", err)
	}

	if s.listCache != nil {
		applications := make([]domain.Application, 0, len(result))
		for _, application := range result {
			if application != nil {
				applications = append(applications, *application)
			}
		}
		_ = s.listCache.Set(ctx, "group:"+groupId, applications, 10)
	}

	return result, nil
}

func (s *Service) List(ctx context.Context, organizationId string) ([]domain.Application, error) {
	if organizationId == "" {
		return []domain.Application{}, fmt.Errorf("organization id cannot be empty")
	}

	if s.listCache != nil {
		cacheKey := "organization:" + organizationId
		cached, found, err := s.listCache.Get(ctx, cacheKey)
		if err == nil && found {
			return cached, nil
		}
	}

	result, err := s.repo.List(ctx, organizationId)
	if err != nil {
		return []domain.Application{}, fmt.Errorf("get: %w", err)
	}

	if s.listCache != nil {
		_ = s.listCache.Set(ctx, "organization:"+organizationId, result, 10)
	}

	return result, nil
}

func (s *Service) GetByOrganizationId(ctx context.Context, organizationId string) ([]*domain.Application, error) {
	result, err := s.repo.GetByOrganizationId(ctx, organizationId)
	if err != nil {
		return []*domain.Application{}, err
	}

	return result, nil
}

func (s *Service) Create(ctx context.Context, application domain.Application) (domain.Application, error) {
	if len(application.OrganizationId) < 5 {
		return domain.Application{}, fmt.Errorf("organization id is invalid")
	}
	if len(application.Name) < 5 {
		return domain.Application{}, fmt.Errorf("name is invalid")
	}
	if len(application.UserId) < 5 {
		return domain.Application{}, fmt.Errorf("user id is invalid")
	}
	switch application.Type {
	case domain.ApplicationTypeStudent:
		if application.GroupId == nil || len(*application.GroupId) < 5 {
			return domain.Application{}, fmt.Errorf("group id is invalid")
		}
	case domain.ApplicationTypeTeacher:
		application.GroupId = nil
	default:
		return domain.Application{}, fmt.Errorf("unsupported application type: %s", application.Type)
	}

	_, err := s.GetByUserId(ctx, application.UserId)

	switch {
	case err == nil:
		return domain.Application{}, fmt.Errorf(
			"already exist %s",
			application.UserId,
		)

	case errors.Is(err, gorm.ErrRecordNotFound):
		if application.Type == domain.ApplicationTypeStudent {
			group, err := s.groupService.GetById(ctx, *application.GroupId)
			if err != nil {
				return domain.Application{}, fmt.Errorf("validate group: %w", err)
			}
			if group.OrganizationId != application.OrganizationId {
				return domain.Application{}, fmt.Errorf("group does not belong to application organization")
			}
		}

		application.Id = uuid.New().String()

		if err := s.repo.Create(ctx, &application); err != nil {
			return domain.Application{}, fmt.Errorf("create applicatiom: %w", err)
		}

		if s.cache != nil {
			_ = s.cache.Set(ctx, application.Id, application, 10)
		}

		return application, nil
	default:
		return domain.Application{}, fmt.Errorf(
			"check existing application: %w",
			err,
		)
	}

}

func (s *Service) Update(
	ctx context.Context,
	application domain.Application,
) (domain.Application, error) {
	if application.Id == "" {
		return domain.Application{}, fmt.Errorf("id cannot be empty")
	}

	current, err := s.repo.GetById(ctx, application.Id)
	if err != nil {
		return domain.Application{}, fmt.Errorf("get application by id: %w", err)
	}

	if application.Name != "" {
		current.Name = application.Name
	}
	if application.Status != "" {
		current.Status = application.Status
	}
	if application.Type != "" {
		current.Type = application.Type
	}
	if application.UserId != "" {
		current.UserId = application.UserId
	}
	if application.OrganizationId != "" {
		current.OrganizationId = application.OrganizationId
	}
	if application.GroupId != nil {
		current.GroupId = application.GroupId
	}

	if err := s.repo.Update(ctx, &current); err != nil {
		return domain.Application{}, fmt.Errorf("update application: %w", err)
	}

	if s.cache != nil {
		_ = s.cache.Set(ctx, current.Id, current, 10)
	}

	return current, nil
}

func (s *Service) UpdateApplicationStatus(
	ctx context.Context,
	application domain.Application,
	status domain.ApplicationStatus,
) (domain.Application, error) {
	if application.Id == "" {
		return domain.Application{}, fmt.Errorf("id cannot be empty")
	}
	if status == "" {
		return domain.Application{}, fmt.Errorf("status cannot be empty")
	}

	current, err := s.repo.GetById(ctx, application.Id)
	if err != nil {
		return domain.Application{}, fmt.Errorf("get application by id: %w", err)
	}

	switch status {
	case domain.ApplicationStatusApproved:
		if s.userService == nil {
			return domain.Application{}, errors.New("user service is not configured")
		}
		user, err := s.userService.GetById(ctx, current.UserId)
		if err != nil {
			return domain.Application{}, fmt.Errorf("get user by id: %w", err)
		}

		switch current.Type {
		case domain.ApplicationTypeStudent:
			if current.GroupId == nil {
				return domain.Application{}, errors.New("student application has no group")
			}
			if s.studentService == nil || s.roleService == nil || s.userService == nil {
				return domain.Application{}, errors.New("student approval services are not configured")
			}

			role, err := s.roleService.GetByName(ctx, "student")
			if err != nil {
				return domain.Application{}, fmt.Errorf("get role by name: %w", err)
			}

			_, err = s.studentService.GetByUserId(ctx, current.UserId)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				student := domain.Student{
					Name:     current.Name,
					IsActive: true,
					GroupId:  *current.GroupId,
					UserId:   current.UserId,
				}
				if _, err := s.studentService.Create(ctx, student); err != nil {
					return domain.Application{}, fmt.Errorf("create student: %w", err)
				}
			} else if err != nil {
				return domain.Application{}, fmt.Errorf("check existing student: %w", err)
			}
			user.RoleId = role.Id
			if _, err := s.userService.Update(ctx, user); err != nil {
				return domain.Application{}, fmt.Errorf("update user: %w", err)
			}
		case domain.ApplicationTypeTeacher:
			if s.teacherService == nil || s.roleService == nil || s.userService == nil {
				return domain.Application{}, errors.New("teacher approval services are not configured")
			}

			role, err := s.roleService.GetByName(ctx, "teacher")
			if err != nil {
				return domain.Application{}, fmt.Errorf("get role by name: %w", err)
			}

			_, err = s.teacherService.GetByUserId(ctx, current.UserId)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				teacher := domain.Teacher{
					Name:           current.Name,
					OrganizationId: current.OrganizationId,
					UserId:         current.UserId,
				}
				if _, err := s.teacherService.Create(ctx, teacher); err != nil {
					return domain.Application{}, fmt.Errorf("create teacher: %w", err)
				}
			} else if err != nil {
				return domain.Application{}, fmt.Errorf("check existing teacher: %w", err)
			}
			user.RoleId = role.Id
			if _, err := s.userService.Update(ctx, user); err != nil {
				return domain.Application{}, fmt.Errorf("update user: %w", err)
			}
		default:
			return domain.Application{}, fmt.Errorf("unsupported application type: %s", current.Type)
		}
	case domain.ApplicationStatusRejected,
		domain.ApplicationStatusPending,
		domain.ApplicationStatusNone:
	default:
		return domain.Application{}, fmt.Errorf("unsupported application status: %s", status)
	}

	current.Status = status

	if err := s.repo.Update(ctx, &current); err != nil {
		return domain.Application{}, fmt.Errorf("update application: %w", err)
	}

	if s.cache != nil {
		_ = s.cache.Set(ctx, current.Id, current, 10)
	}

	return current, nil
}
