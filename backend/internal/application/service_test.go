package application

import (
	"context"
	"testing"

	"github.com/exexeeeex/fiks-api/internal/domain"
	"gorm.io/gorm"
)

type organizationListCache struct {
	values map[string][]domain.Application
}

func (c *organizationListCache) Get(
	_ context.Context,
	key string,
) ([]domain.Application, bool, error) {
	value, found := c.values[key]
	return value, found, nil
}

func (c *organizationListCache) Set(
	_ context.Context,
	key string,
	value []domain.Application,
	_ int,
) error {
	c.values[key] = value
	return nil
}

type organizationListRepository struct {
	Repository
	calls int
}

func (r *organizationListRepository) List(
	_ context.Context,
	organizationId string,
) ([]domain.Application, error) {
	r.calls++
	return []domain.Application{{OrganizationId: organizationId}}, nil
}

func TestListCachesResultsPerOrganization(t *testing.T) {
	repository := &organizationListRepository{}
	cache := &organizationListCache{values: make(map[string][]domain.Application)}
	service := NewService(repository, nil, cache, nil, nil, nil, nil, nil)

	for _, organizationId := range []string{"organization-1", "organization-2", "organization-1"} {
		applications, err := service.List(context.Background(), organizationId)
		if err != nil {
			t.Fatalf("List(%q) error = %v", organizationId, err)
		}
		if len(applications) != 1 || applications[0].OrganizationId != organizationId {
			t.Errorf("List(%q) returned %#v", organizationId, applications)
		}
	}

	if repository.calls != 2 {
		t.Errorf("repository calls = %d, want 2", repository.calls)
	}
}

type createApplicationRepository struct {
	Repository
	created *domain.Application
}

func (r *createApplicationRepository) GetByUserId(context.Context, string) (domain.Application, error) {
	return domain.Application{}, gorm.ErrRecordNotFound
}

func (r *createApplicationRepository) Create(_ context.Context, application *domain.Application) error {
	r.created = application
	return nil
}

type createApplicationGroupService struct {
	group domain.Group
	calls int
}

func (s *createApplicationGroupService) GetById(context.Context, string) (domain.Group, error) {
	s.calls++
	return s.group, nil
}

func TestCreateTeacherApplicationStoresNoGroup(t *testing.T) {
	repository := &createApplicationRepository{}
	groups := &createApplicationGroupService{}
	service := NewService(repository, nil, nil, groups, nil, nil, nil, nil)
	groupId := "group-id"

	created, err := service.Create(context.Background(), domain.Application{
		Name:           "Teacher Applicant",
		Type:           domain.ApplicationTypeTeacher,
		UserId:         "user-id",
		OrganizationId: "organization-id",
		GroupId:        &groupId,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.GroupId != nil || repository.created.GroupId != nil {
		t.Fatalf("teacher group id = %v, want nil", created.GroupId)
	}
	if groups.calls != 0 {
		t.Fatalf("group lookup calls = %d, want 0", groups.calls)
	}
}

func TestCreateStudentApplicationRequiresGroupInOrganization(t *testing.T) {
	tests := []struct {
		name    string
		groupId *string
		group   domain.Group
		wantErr bool
	}{
		{
			name:    "missing group",
			wantErr: true,
		},
		{
			name:    "group from another organization",
			groupId: stringPointer("group-id"),
			group: domain.Group{
				Id:             "group-id",
				OrganizationId: "other-organization",
			},
			wantErr: true,
		},
		{
			name:    "valid group",
			groupId: stringPointer("group-id"),
			group: domain.Group{
				Id:             "group-id",
				OrganizationId: "organization-id",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &createApplicationRepository{}
			groups := &createApplicationGroupService{group: tt.group}
			service := NewService(repository, nil, nil, groups, nil, nil, nil, nil)

			created, err := service.Create(context.Background(), domain.Application{
				Name:           "Student Applicant",
				Type:           domain.ApplicationTypeStudent,
				UserId:         "user-id",
				OrganizationId: "organization-id",
				GroupId:        tt.groupId,
			})
			if (err != nil) != tt.wantErr {
				t.Fatalf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && (created.GroupId == nil || *created.GroupId != *tt.groupId) {
				t.Fatalf("created group id = %v, want %q", created.GroupId, *tt.groupId)
			}
		})
	}
}

func stringPointer(value string) *string {
	return &value
}
