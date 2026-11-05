package maxauth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/exexeeeex/fiks-api/internal/domain"
	jwt "github.com/exexeeeex/fiks-api/internal/jwt"
	jwtgo "github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type testApplicationService struct {
	application domain.Application
	err         error
	requested   string
}

func (s *testApplicationService) GetByMaxId(_ context.Context, maxId string) (domain.Application, error) {
	s.requested = maxId
	return s.application, s.err
}

type testUserService struct {
	user  domain.User
	input domain.User
}

func (s *testUserService) Create(_ context.Context, user domain.User) (domain.User, error) {
	s.input = user
	return s.user, nil
}

func (s *testUserService) GetByMaxId(_ context.Context, _ string) (domain.User, error) {
	return s.user, nil
}

type testRoleService struct {
	role      domain.Role
	requested string
}

func (s *testRoleService) GetById(_ context.Context, id string) (domain.Role, error) {
	s.requested = id
	return s.role, nil
}

type testJwtService struct {
	request      jwt.AccessRequest
	token        *jwtgo.Token
	validateErr  error
	validatedTyp string
}

func (s *testJwtService) WriteTokensPair(request jwt.AccessRequest) (jwt.TokensPair, error) {
	s.request = request
	return jwt.TokensPair{AccessToken: "access", RefreshToken: "refresh"}, nil
}

func (s *testJwtService) ValidateToken(_ string, tokenType string) (*jwtgo.Token, error) {
	s.validatedTyp = tokenType
	return s.token, s.validateErr
}

func TestAuthenticateUsesUserReturnedByMaxIDLookup(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("MAX_BOT_TOKEN", "test-bot-token")

	user := MaxUser{Id: 12345, FirstName: "Test"}
	initData, err := NewService(nil, nil, nil, nil).GenerateTestInitialData(user, time.Now())
	if err != nil {
		t.Fatal(err)
	}

	applicationService := &testApplicationService{
		application: domain.Application{
			Id:             "application-id",
			Name:           "Egor egorov egorovich",
			Status:         domain.ApplicationStatusPending,
			Type:           domain.ApplicationTypeStudent,
			OrganizationId: "organization-id",
			GroupId:        "group-id",
		},
	}

	userService := &testUserService{
		user: domain.User{
			Id:     "existing-user-id",
			MaxId:  "12345",
			RoleId: "admin-role-id",
		},
	}
	jwtService := &testJwtService{}
	roleService := &testRoleService{
		role: domain.Role{Id: "admin-role-id", Name: "admin"},
	}
	service := NewService(userService, jwtService, roleService, applicationService)

	response, err := service.Authenticate(context.Background(), initData)
	if err != nil {
		t.Fatal(err)
	}

	if jwtService.request.Id != userService.user.Id {
		t.Errorf("token user id = %q, want %q", jwtService.request.Id, userService.user.Id)
	}
	if jwtService.request.MaxId != userService.user.MaxId {
		t.Errorf("token max id = %q, want %q", jwtService.request.MaxId, userService.user.MaxId)
	}
	if jwtService.request.Role != roleService.role.Name {
		t.Errorf("token role = %q, want %q", jwtService.request.Role, roleService.role.Name)
	}
	if roleService.requested != userService.user.RoleId {
		t.Errorf("role lookup id = %q, want %q", roleService.requested, userService.user.RoleId)
	}
	if userService.input.MaxId != "12345" {
		t.Errorf("lookup max id = %q, want %q", userService.input.MaxId, "12345")
	}
	if !response.Success || response.Tokens.AccessToken != "access" {
		t.Errorf("unexpected auth response: %+v", response)
	}
	if response.User.Id != userService.user.Id {
		t.Errorf("response user_id = %q, want %q", response.User.Id, userService.user.Id)
	}
	if response.User.MaxId != userService.user.MaxId {
		t.Errorf("response max_id = %q, want %q", response.User.MaxId, userService.user.MaxId)
	}
	if jwtService.request.ApplicationStatus != domain.ApplicationStatusPending ||
		response.User.ApplicationStatus != domain.ApplicationStatusPending {
		t.Errorf("application status was not propagated: token=%q response=%q",
			jwtService.request.ApplicationStatus, response.User.ApplicationStatus)
	}
}

func TestAuthenticateAllowsUserWithoutApplication(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("MAX_BOT_TOKEN", "test-bot-token")

	initData, err := NewService(nil, nil, nil, nil).GenerateTestInitialData(
		MaxUser{Id: 12345},
		time.Now(),
	)
	if err != nil {
		t.Fatal(err)
	}

	userService := &testUserService{
		user: domain.User{Id: "existing-user-id", MaxId: "12345", RoleId: "guest-role-id"},
	}
	jwtService := &testJwtService{}
	roleService := &testRoleService{
		role: domain.Role{Id: "guest-role-id", Name: "guest"},
	}
	applicationService := &testApplicationService{err: gorm.ErrRecordNotFound}
	service := NewService(userService, jwtService, roleService, applicationService)

	response, err := service.Authenticate(context.Background(), initData)
	if err != nil {
		t.Fatalf("Authenticate() error = %v, want nil", err)
	}
	if !response.Success {
		t.Fatal("Authenticate() returned success=false")
	}
	if response.User.ApplicationStatus != domain.ApplicationStatusNone {
		t.Errorf("response application status = %q, want %q",
			response.User.ApplicationStatus, domain.ApplicationStatusNone)
	}
	if jwtService.request.ApplicationStatus != domain.ApplicationStatusNone {
		t.Errorf("token application status = %q, want %q",
			jwtService.request.ApplicationStatus, domain.ApplicationStatusNone)
	}
}

func TestAuthenticateReturnsApplicationLookupErrors(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("MAX_BOT_TOKEN", "test-bot-token")

	initData, err := NewService(nil, nil, nil, nil).GenerateTestInitialData(
		MaxUser{Id: 12345},
		time.Now(),
	)
	if err != nil {
		t.Fatal(err)
	}
	applicationService := &testApplicationService{err: errors.New("database unavailable")}
	service := NewService(
		&testUserService{user: domain.User{Id: "user-id", MaxId: "12345", RoleId: "role-id"}},
		&testJwtService{},
		&testRoleService{role: domain.Role{Id: "role-id", Name: "guest"}},
		applicationService,
	)

	if _, err := service.Authenticate(context.Background(), initData); err == nil {
		t.Fatal("Authenticate() error = nil, want application lookup error")
	}
}

func TestRefreshUpdatesApplicationStatusInAccessToken(t *testing.T) {
	jwtService := &testJwtService{
		token: &jwtgo.Token{
			Valid: true,
			Claims: jwtgo.MapClaims{
				"id":    "user-id",
				"maxId": "12345",
			},
		},
	}
	userService := &testUserService{
		user: domain.User{Id: "user-id", MaxId: "12345", RoleId: "role-id"},
	}
	roleService := &testRoleService{role: domain.Role{Id: "role-id", Name: "student"}}
	applicationService := &testApplicationService{
		application: domain.Application{Status: domain.ApplicationStatusApproved},
	}
	service := NewService(userService, jwtService, roleService, applicationService)

	tokens, err := service.Refresh(context.Background(), "refresh-token")
	if err != nil {
		t.Fatalf("Refresh() error = %v, want nil", err)
	}
	if tokens.AccessToken != "access" {
		t.Errorf("access token = %q, want %q", tokens.AccessToken, "access")
	}
	if jwtService.validatedTyp != "refresh" {
		t.Errorf("validated token type = %q, want %q", jwtService.validatedTyp, "refresh")
	}
	if jwtService.request.ApplicationStatus != domain.ApplicationStatusApproved {
		t.Errorf("refreshed application status = %q, want %q",
			jwtService.request.ApplicationStatus, domain.ApplicationStatusApproved)
	}
}

func TestRefreshRejectsInvalidRefreshToken(t *testing.T) {
	jwtService := &testJwtService{validateErr: errors.New("expired")}
	service := NewService(nil, jwtService, nil, nil)

	_, err := service.Refresh(context.Background(), "expired-token")
	if !errors.Is(err, ErrInvalidRefreshToken) {
		t.Errorf("Refresh() error = %v, want ErrInvalidRefreshToken", err)
	}
}
