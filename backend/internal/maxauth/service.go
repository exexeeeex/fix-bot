package maxauth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"uuid"

	"github.com/exexeeeex/fiks-api/internal/domain"
	jwt "github.com/exexeeeex/fiks-api/internal/jwt"
	jwtgo "github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type ApplicationService interface {
	GetByMaxId(ctx context.Context, maxId string) (domain.Application, error)
}

type RoleService interface {
	GetById(ctx context.Context, id string) (domain.Role, error)
}

type JwtService interface {
	WriteTokensPair(jwt.AccessRequest) (jwt.TokensPair, error)
	ValidateToken(tokenString string, tokenType string) (*jwtgo.Token, error)
}

type UserService interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	GetByMaxId(ctx context.Context, maxId string) (domain.User, error)
}

type StaffService interface {
	GetByUserId(ctx context.Context, userId string) (domain.Teacher, error)
}

var ErrInvalidRefreshToken = errors.New("invalid refresh token")
var ErrRefreshTokenRequired = errors.New("refresh token is required")
var ErrInvalidDemoRole = errors.New("demo role must be student or teacher")

type Service struct {
	testEnabled        bool
	userService        UserService
	jwtService         JwtService
	roleService        RoleService
	applicationService ApplicationService
	staffService       StaffService
}

func NewService(
	userService UserService,
	jwtService JwtService,
	roleService RoleService,
	applicationService ApplicationService,
	staffServices ...StaffService,
) *Service {
	service := &Service{
		testEnabled:        os.Getenv("APP_ENV") == "development",
		userService:        userService,
		jwtService:         jwtService,
		roleService:        roleService,
		applicationService: applicationService,
	}
	if len(staffServices) > 0 {
		service.staffService = staffServices[0]
	}
	return service
}

func (s *Service) Authenticate(ctx context.Context, data string) (AuthResponse, error) {
	if data == "" {
		return AuthResponse{}, errors.New("data is required")
	}

	params, err := s.ValidateInitData(data)
	if err != nil {
		return AuthResponse{}, err
	}

	authDate, err := strconv.ParseInt(
		params.Get("auth_date"),
		10,
		64,
	)
	if err != nil || authDate <= 0 {
		return AuthResponse{}, errors.New("invalid auth_date")
	}

	createdAt := time.Unix(authDate, 0)
	now := time.Now()

	if createdAt.After(now.Add(time.Minute)) {
		return AuthResponse{}, errors.New("auth_date is not the future")
	}

	if !s.testEnabled {
		if now.Sub(createdAt) > time.Hour {
			return AuthResponse{}, errors.New("init data expired")
		}
	}

	userJSON := params.Get("user")
	if userJSON == "" {
		return AuthResponse{}, errors.New("user not found")
	}

	var user MaxUser

	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return AuthResponse{}, errors.New("invalid user data")
	}

	if user.Id <= 0 {
		return AuthResponse{}, errors.New("invalid user id")
	}

	authenticatedUser, err := s.userService.Create(ctx, domain.User{
		Id:    uuid.NewV4().String(),
		MaxId: strconv.Itoa(user.Id),
	})

	if err != nil {
		return AuthResponse{}, fmt.Errorf("get or create user: %w", err)
	}

	role, err := s.roleService.GetById(ctx, authenticatedUser.RoleId)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("get user role: %w", err)
	}

	applicationStatus, organizationId, err := s.getApplicationData(ctx, authenticatedUser.MaxId)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("get application status: %w", err)
	}
	if role.Name == "admin" || role.Name == "teacher" {
		if s.staffService != nil {
			staff, err := s.staffService.GetByUserId(ctx, authenticatedUser.Id)
			if err != nil {
				return AuthResponse{}, fmt.Errorf("get staff organization: %w", err)
			}
			organizationId = staff.OrganizationId
		}
	}

	tokens, err := s.jwtService.WriteTokensPair(jwt.AccessRequest{
		Id:                authenticatedUser.Id,
		MaxId:             authenticatedUser.MaxId,
		ApplicationStatus: applicationStatus,
		Role:              role.Name,
		OrganizationId:    organizationId,
	})

	if err != nil {
		return AuthResponse{}, fmt.Errorf("write tokens: %w", err)
	}

	return AuthResponse{
		Success: true,
		Tokens:  tokens,
		User: AuthenticatedUser{
			Id:                authenticatedUser.Id,
			MaxId:             authenticatedUser.MaxId,
			FirstName:         user.FirstName,
			LastName:          user.LastName,
			Username:          user.Username,
			LanguageCode:      user.LanguageCode,
			PhotoURL:          user.PhotoURL,
			Role:              role.Name,
			OrganizationId:    organizationId,
			ApplicationStatus: applicationStatus,
		},
	}, nil
}

func (s *Service) AuthenticateDemo(ctx context.Context, requestedRole string) (AuthResponse, error) {
	var maxId, firstName, lastName string
	switch requestedRole {
	case "student":
		maxId, firstName, lastName = "demo-student", "Демо", "Студент"
	case "teacher":
		maxId, firstName, lastName = "demo-teacher", "Демо", "Преподаватель"
	default:
		return AuthResponse{}, ErrInvalidDemoRole
	}

	user, err := s.userService.GetByMaxId(ctx, maxId)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("get demo user: %w", err)
	}
	role, err := s.roleService.GetById(ctx, user.RoleId)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("get demo user role: %w", err)
	}
	if role.Name != requestedRole {
		return AuthResponse{}, fmt.Errorf("demo user role mismatch: got %q", role.Name)
	}

	applicationStatus, organizationId, err := s.getApplicationData(ctx, user.MaxId)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("get demo application status: %w", err)
	}
	if applicationStatus != domain.ApplicationStatusApproved {
		return AuthResponse{}, errors.New("demo application is not approved")
	}
	if requestedRole == "teacher" {
		if s.staffService == nil {
			return AuthResponse{}, errors.New("demo teacher service is not configured")
		}
		staff, err := s.staffService.GetByUserId(ctx, user.Id)
		if err != nil {
			return AuthResponse{}, fmt.Errorf("get demo teacher organization: %w", err)
		}
		organizationId = staff.OrganizationId
	}

	tokens, err := s.jwtService.WriteTokensPair(jwt.AccessRequest{
		Id:                user.Id,
		MaxId:             user.MaxId,
		ApplicationStatus: applicationStatus,
		Role:              role.Name,
		OrganizationId:    organizationId,
	})
	if err != nil {
		return AuthResponse{}, fmt.Errorf("write demo tokens: %w", err)
	}

	return AuthResponse{
		Success: true,
		Tokens:  tokens,
		User: AuthenticatedUser{
			Id:                user.Id,
			MaxId:             user.MaxId,
			FirstName:         firstName,
			LastName:          lastName,
			Username:          maxId,
			Role:              role.Name,
			OrganizationId:    organizationId,
			ApplicationStatus: applicationStatus,
		},
	}, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (jwt.TokensPair, error) {
	if refreshToken == "" {
		return jwt.TokensPair{}, ErrRefreshTokenRequired
	}

	token, err := s.jwtService.ValidateToken(refreshToken, "refresh")
	if err != nil {
		return jwt.TokensPair{}, fmt.Errorf("%w: %v", ErrInvalidRefreshToken, err)
	}

	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok || !token.Valid {
		return jwt.TokensPair{}, ErrInvalidRefreshToken
	}

	userId, ok := claims["id"].(string)
	if !ok || userId == "" {
		return jwt.TokensPair{}, fmt.Errorf("%w: user id is missing", ErrInvalidRefreshToken)
	}
	maxId, ok := claims["maxId"].(string)
	if !ok || maxId == "" {
		return jwt.TokensPair{}, fmt.Errorf("%w: max id is missing", ErrInvalidRefreshToken)
	}

	user, err := s.userService.GetByMaxId(ctx, maxId)
	if err != nil {
		return jwt.TokensPair{}, fmt.Errorf("get user for refresh: %w", err)
	}
	if user.Id != userId {
		return jwt.TokensPair{}, errors.New("refresh token user does not match")
	}

	role, err := s.roleService.GetById(ctx, user.RoleId)
	if err != nil {
		return jwt.TokensPair{}, fmt.Errorf("get user role: %w", err)
	}
	applicationStatus, organizationId, err := s.getApplicationData(ctx, user.MaxId)
	if err != nil {
		return jwt.TokensPair{}, fmt.Errorf("get application status: %w", err)
	}
	if role.Name == "admin" || role.Name == "teacher" {
		if s.staffService != nil {
			staff, err := s.staffService.GetByUserId(ctx, user.Id)
			if err != nil {
				return jwt.TokensPair{}, fmt.Errorf("get staff organization: %w", err)
			}
			organizationId = staff.OrganizationId
		}
	}

	tokens, err := s.jwtService.WriteTokensPair(jwt.AccessRequest{
		Id:                user.Id,
		MaxId:             user.MaxId,
		ApplicationStatus: applicationStatus,
		Role:              role.Name,
		OrganizationId:    organizationId,
	})
	if err != nil {
		return jwt.TokensPair{}, fmt.Errorf("write refreshed tokens: %w", err)
	}
	return tokens, nil
}

func (s *Service) getApplicationStatus(ctx context.Context, maxId string) (domain.ApplicationStatus, error) {
	status, _, err := s.getApplicationData(ctx, maxId)
	return status, err
}

func (s *Service) getApplicationData(
	ctx context.Context,
	maxId string,
) (domain.ApplicationStatus, string, error) {
	application, err := s.applicationService.GetByMaxId(ctx, maxId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ApplicationStatusNone, "", nil
	}
	if err != nil {
		return "", "", err
	}
	return application.Status, application.OrganizationId, nil
}

func (s *Service) ValidateInitData(
	data string,
) (url.Values, error) {
	params, err := url.ParseQuery(data)
	if err != nil {
		return nil, errors.New("invalid query string")
	}

	for _, values := range params {
		if len(values) != 1 {
			return nil, errors.New("duplicate parameter")
		}
	}

	hashValues, ok := params["hash"]
	if !ok || len(hashValues) != 1 || hashValues[0] == "" {
		return nil, errors.New("hash is missing")
	}

	originalHash := hashValues[0]

	if _, err := hex.DecodeString(originalHash); err != nil {
		return nil, errors.New("Invalid hash format")
	}

	keys := make([]string, 0, len(params))

	for key := range params {
		if key != "hash" {
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)

	lines := make([]string, 0, len(keys))

	for _, key := range keys {
		lines = append(lines, key+"="+params.Get(key))
	}

	secretMAC := hmac.New(sha256.New, []byte("WebAppData"))
	secretMAC.Write([]byte(os.Getenv("MAX_BOT_TOKEN")))
	secretKey := secretMAC.Sum(nil)

	dataMAC := hmac.New(sha256.New, secretKey)
	dataMAC.Write([]byte(strings.Join(lines, "\n")))
	expectedHash := dataMAC.Sum(nil)

	actualHash, err := hex.DecodeString(originalHash)
	if err != nil {
		return nil, errors.New("invalid hash")
	}

	if !hmac.Equal(expectedHash, actualHash) {
		return nil, errors.New("signature mismatch")
	}

	return params, nil
}

func (s *Service) GenerateTestInitialData(
	user MaxUser,
	createdAt time.Time,
) (string, error) {
	if os.Getenv("MAX_BOT_TOKEN") == "" {
		return "", errors.New("MAX_BOT_TOKEN is not configured")
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("marshal user: %w", err)
	}

	params := url.Values{}
	params.Set("auth_date", strconv.FormatInt(createdAt.Unix(), 10))
	params.Set("user", string(userJSON))

	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	lines := make([]string, 0, len(keys))
	for _, key := range keys {
		lines = append(lines, key+"="+params.Get(key))
	}

	dataCheckString := strings.Join(lines, "\n")

	secretMAC := hmac.New(sha256.New, []byte("WebAppData"))
	_, _ = secretMAC.Write([]byte(os.Getenv("MAX_BOT_TOKEN")))
	secretKey := secretMAC.Sum(nil)

	dataMAC := hmac.New(sha256.New, secretKey)
	_, _ = dataMAC.Write([]byte(dataCheckString))
	hash := hex.EncodeToString(dataMAC.Sum(nil))

	params.Set("hash", hash)

	return params.Encode(), nil
}
