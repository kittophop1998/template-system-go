package main

import (
	"context"
	"file_storage_service/infrastructure/config"
	"file_storage_service/internal/app/handler/http"
	"file_storage_service/internal/app/usecase"
	"file_storage_service/internal/platform/database"
	"file_storage_service/internal/platform/storage"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func initializeApp(cfg *config.Config) (*gin.Engine, error) {
	// ===== Setup Database =====
	db, err := setupDatabase(cfg)
	if err != nil {
		return nil, err
	}

	// ===== Setup Object Storage =====
	s3Client := storage.NewS3Client(cfg)
	if err := ensureBucket(s3Client, cfg.Storage.BucketName); err != nil {
		return nil, err
	}

	// ===== Initialize Repositories =====
	userRepo := database.NewUserPostGres(db)

	// ===== Initialize UseCases =====
	userUC := usecase.NewUserUseCase(userRepo)
	attachmentUC := &usecase.AttachmentUseCase{
		S3:     s3Client,
		Bucket: cfg.Storage.BucketName,
	}

	// ===== Initialize Handlers =====
	handlers := http.InitializeHandlers(&http.HandlerDependency{
		UserUC:       userUC,
		AttachmentUC: attachmentUC,
	})

	// ===== Setup Router =====
	router := setupRouter(cfg, handlers)

	return router, nil
}

func setupDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Bangkok",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
	)

	db, err := database.InitializePostgres(dsn)
	if err != nil {
		return nil, err
	}

	if err := database.Migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func ensureBucket(s3Client *s3.Client, bucketName string) error {
	_, err := s3Client.CreateBucket(context.Background(), &s3.CreateBucketInput{
		Bucket: &bucketName,
	})
	if err != nil {
		// ถ้า bucket มีอยู่แล้ว ให้ ignore
		if strings.Contains(err.Error(), "BucketAlreadyOwnedByYou") || strings.Contains(err.Error(), "BucketAlreadyExists") {
			log.Printf("⚠️ bucket %s already exists, skip creating", bucketName)
			return nil
		}
		return err // return error จริง
	}

	log.Printf("✅ bucket %s created successfully", bucketName)
	return nil
}

func setupRouter(cfg *config.Config, handlers *http.Handlers) *gin.Engine {
	gin.SetMode(cfg.Server.GinMode)
	router := gin.New()

	// ===== Middleware =====
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept", "Authorization"},
		MaxAge:          12 * time.Hour,
	}))
	router.Use(gin.Recovery())

	// ===== Setup Routes =====
	http.SetupRoutes(router, handlers)

	return router
}
