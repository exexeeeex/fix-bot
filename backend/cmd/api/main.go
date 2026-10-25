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
	"github.com/exexeeeex/fiks-api/internal/domain"
	"github.com/exexeeeex/fiks-api/internal/httpserver"

	groupapi "github.com/exexeeeex/fiks-api/internal/group"
	qrcodeapi "github.com/exexeeeex/fiks-api/internal/qr"
	studentapi "github.com/exexeeeex/fiks-api/internal/student"

	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	redis "github.com/redis/go-redis/v9"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx := context.Background()

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

	groupHandler, groupService := CreateGroup(db, redisClient)
	studentHandler, studentService := CreateStudent(db, redisClient, groupService)
	qrHander, _ := CreateQrCode(studentService, groupService)

	server := httpserver.New(
		studentHandler,
		groupHandler,
		qrHander,
	)
	server.Addr = ":" + cfg.Port

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("HTTP server starting", "address", server.Addr)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	stopCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

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

func CreateStudent(
	db *gorm.DB,
	redis *redis.Client,
	groupService studentapi.GroupService,
) (*studentapi.Handler, *studentapi.Service) {
	repo := studentapi.NewRepository(db)
	studentCache := cache.New[domain.Student](redis, "student:")
	studentListCache := cache.New[[]domain.Student](redis, "students:")

	service := studentapi.NewService(
		repo,
		studentCache,
		studentListCache,
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

func CreateQrCode(
	studentService *studentapi.Service,
	groupService *groupapi.Service,
) (*qrcodeapi.Handler, *qrcodeapi.Service) {
	service := qrcodeapi.NewService(studentService, groupService)
	handler := qrcodeapi.NewHandler(service)

	return handler, service
}
