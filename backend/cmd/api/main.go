package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/exexeeeex/fiks-api/internal/cache"
	"github.com/exexeeeex/fiks-api/internal/config"
	configmiddleware "github.com/exexeeeex/fiks-api/internal/config/middlewares"
	"github.com/exexeeeex/fiks-api/internal/domain"
	"github.com/exexeeeex/fiks-api/internal/httpserver"

	applicationapi "github.com/exexeeeex/fiks-api/internal/application"
	groupapi "github.com/exexeeeex/fiks-api/internal/group"
	jwtapi "github.com/exexeeeex/fiks-api/internal/jwt"
	maxauthapi "github.com/exexeeeex/fiks-api/internal/maxauth"
	organizationapi "github.com/exexeeeex/fiks-api/internal/organization"
	qrcodeapi "github.com/exexeeeex/fiks-api/internal/qr"
	roleapi "github.com/exexeeeex/fiks-api/internal/role"
	statisticapi "github.com/exexeeeex/fiks-api/internal/statistic"
	studentapi "github.com/exexeeeex/fiks-api/internal/student"
	teacherapi "github.com/exexeeeex/fiks-api/internal/teacher"
	userapi "github.com/exexeeeex/fiks-api/internal/user"
	visitapi "github.com/exexeeeex/fiks-api/internal/visit"

	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	redis "github.com/redis/go-redis/v9"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	go bot(ctx, cfg.MaxBotToken)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	db, err := gorm.Open(
		gormpostgres.Open(cfg.DatabaseURL),
		&gorm.Config{
			Logger: gormlogger.Default.LogMode(gormlogger.Warn),
		},
	)
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	pingCtx, cancelPing := context.WithTimeout(ctx, 5*time.Second)
	defer cancelPing()

	if err := sqlDB.PingContext(pingCtx); err != nil {
		return err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	defer redisClient.Close()

	redisCtx, cancelRedis := context.WithTimeout(ctx, 5*time.Second)
	err = redisClient.Ping(redisCtx).Err()
	cancelRedis()
	if err != nil {
		return err
	}

	// role
	roleService := CreateRole(db, redisClient)

	// jwt
	jwtService := jwtapi.NewService()

	// protected routes middleware
	protectedRoutes := configmiddleware.NewProtectedRoutes(jwtService)

	// group
	groupHandler, groupService := CreateGroup(db, redisClient)

	// student
	studentHandler, studentService := CreateStudent(db, redisClient, groupService)

	// visit
	visitHandler, visitService := CreateVisit(db, protectedRoutes, groupService)

	// qr
	qrHander, _ := CreateQrCode(studentService, groupService, visitService, protectedRoutes)

	// organization
	organizationService, organizationHandler := CreateOrganization(db, redisClient, protectedRoutes)

	// teacher
	teacherService := CreateTeacher(db, redisClient, organizationService)

	// user
	_, userService := CreateUser(db, redisClient, roleService)

	// statistic
	statisticHandler, _ := CreateStatistic(db, protectedRoutes, visitService, studentService)

	// application
	applicationHandler, applicationService := CreateApplication(
		db,
		redisClient,
		groupService,
		studentService,
		teacherService,
		roleService,
		userService,
		protectedRoutes,
	)

	// auth
	maxauthHandler := CreateMaxAuth(
		userService,
		jwtService,
		roleService,
		applicationService,
		teacherService,
	)

	server := httpserver.New(
		studentHandler,
		groupHandler,
		qrHander,
		applicationHandler,
		maxauthHandler,
		organizationHandler,
		visitHandler,
		statisticHandler,
	)
	server.Addr = ":" + cfg.Port

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("HTTP server starting", "address", server.Addr)
		if os.Getenv("APP_ENV") == "development" {
			logger.Info("development mode")
		}

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	stopCtx := ctx

	select {
	case err := <-serverErrors:
		return err
	case <-stopCtx.Done():
		logger.Info("shutdown signal received")
	}

	shutdownCtx, cancelShutdown := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancelShutdown()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	logger.Info("HTTP server stopped")
	return nil
}

func CreateRole(
	db *gorm.DB,
	redis *redis.Client,
) *roleapi.Service {
	repo := roleapi.NewRepository(db)
	roleCache := cache.New[domain.Role](redis, "role:")
	roleListCache := cache.New[[]domain.Role](redis, "roles:")
	service := roleapi.NewService(repo, roleCache, roleListCache)

	return service
}

func CreateStudent(
	db *gorm.DB,
	redis *redis.Client,
	groupService studentapi.GroupService,
) (*studentapi.Handler, *studentapi.Service) {
	repo := studentapi.NewRepository(db)
	studentCache := cache.New[domain.Student](redis, "student:")
	studentListCache := cache.New[[]domain.Student](redis, "students:")
	studentKeyCache := cache.NewStudentKeyCache(redis)

	service := studentapi.NewService(
		repo,
		studentCache,
		studentListCache,
		studentKeyCache,
		groupService,
	)

	handler := studentapi.NewHandler(service)
	return handler, service
}

func CreateGroup(
	db *gorm.DB,
	redis *redis.Client,
) (*groupapi.Handler, *groupapi.Service) {
	repo := groupapi.NewRepository(db)
	groupCache := cache.New[domain.Group](redis, "group:")
	groupListCache := cache.New[[]domain.Group](redis, "groups:")
	service := groupapi.NewService(repo, groupCache, groupListCache)
	handler := groupapi.NewHandler(service)

	return handler, service
}

func CreateApplication(
	db *gorm.DB,
	redis *redis.Client,
	groupService applicationapi.GroupService,
	studentService applicationapi.StudentService,
	teacherService applicationapi.TeacherService,
	roleService applicationapi.RoleService,
	userService applicationapi.UserService,
	protectedRoutes applicationapi.ProtectedRoutes,
) (*applicationapi.Handler, *applicationapi.Service) {
	repo := applicationapi.NewRepository(db)
	applicatiomCache := cache.New[domain.Application](redis, "application:")
	applicationListCache := cache.New[[]domain.Application](redis, "applications:")
	service := applicationapi.NewService(
		repo,
		applicatiomCache,
		applicationListCache,
		groupService,
		studentService,
		teacherService,
		roleService,
		userService,
	)
	handler := applicationapi.NewHandler(service, protectedRoutes)

	return handler, service
}

func CreateUser(
	db *gorm.DB,
	redis *redis.Client,
	roleService *roleapi.Service,
) (*userapi.Handler, *userapi.Service) {
	repo := userapi.NewRepository(db)
	userCache := cache.New[domain.User](redis, "user:")
	userListCache := cache.New[[]domain.User](redis, "users:")
	service := userapi.NewService(repo, userCache, userListCache, roleService)
	handler := userapi.NewHandler(service)

	return handler, service
}

func CreateQrCode(
	studentService *studentapi.Service,
	groupService *groupapi.Service,
	visitService *visitapi.Service,
	protectedRoutes applicationapi.ProtectedRoutes,
) (*qrcodeapi.Handler, *qrcodeapi.Service) {
	service := qrcodeapi.NewService(studentService, groupService, visitService)
	handler := qrcodeapi.NewHandler(service, protectedRoutes)

	return handler, service
}

func CreateOrganization(
	db *gorm.DB,
	redis *redis.Client,
	protectedRoutes organizationapi.ProtectedRoutes,
) (*organizationapi.Service, *organizationapi.Handler) {
	repository := organizationapi.NewRepository(db)

	organizationCache := cache.New[domain.Organization](redis, "organization:")
	organizationListCache := cache.New[[]domain.Organization](redis, "organizations:")

	service := organizationapi.NewService(
		repository,
		organizationCache,
		organizationListCache,
	)
	handler := organizationapi.NewHandler(service, protectedRoutes)

	return service, handler
}

func CreateTeacher(
	db *gorm.DB,
	redis *redis.Client,
	organizationService teacherapi.OrganizationService,
) *teacherapi.Service {
	repository := teacherapi.NewRepository(db)

	teacherCache := cache.New[domain.Teacher](redis, "teacher:")
	teacherListCache := cache.New[[]domain.Teacher](redis, "teachers:")

	service := teacherapi.NewService(
		repository,
		teacherCache,
		teacherListCache,
		organizationService,
	)

	return service
}

func CreateStatistic(
	db *gorm.DB,
	protectedRoutes statisticapi.ProtectedRoutes,
	visitService statisticapi.VisitService,
	studentSerive statisticapi.StudentService,
) (*statisticapi.Handler, *statisticapi.Service) {
	service := statisticapi.NewService(visitService, studentSerive)
	handler := statisticapi.NewHandler(service, protectedRoutes)

	return handler, service
}

func CreateVisit(
	db *gorm.DB,
	protectedRoutes visitapi.ProtectedRoutes,
	groupService visitapi.GroupService,
) (*visitapi.Handler, *visitapi.Service) {
	repo := visitapi.NewRepository(db)
	service := visitapi.NewService(repo)
	handler := visitapi.NewHandler(service, protectedRoutes, groupService)

	return handler, service
}

func CreateMaxAuth(
	userService *userapi.Service,
	jwtService *jwtapi.Service,
	roleService *roleapi.Service,
	applicationService *applicationapi.Service,
	staffService maxauthapi.StaffService,
) *maxauthapi.Handler {
	service := maxauthapi.NewService(
		userService,
		jwtService,
		roleService,
		applicationService,
		staffService,
	)
	handler := maxauthapi.NewHandler(service)

	return handler
}
