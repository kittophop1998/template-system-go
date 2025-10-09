package database

import (
	"context"
	"file_storage_service/internal/domain/model"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBStore *gorm.DB

func InitializePostgres(dsn string) (*gorm.DB, error) {
	env := os.Getenv("ENV_NAME")
	dsnRailway := os.Getenv("DATABASE_URL")

	if env == "railway" && dsnRailway != "" {
		dsn = dsnRailway
	}

	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	sqlDB := stdlib.OpenDB(*cfg)

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// ===== Connection Pool Settings =====
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// ===== Ping Database =====
	if err := sqlDB.PingContext(context.Background()); err != nil {
		log.Printf("Database ping failed: %v", err)
		return nil, err
	}

	DBStore = db

	return db, nil
}

// ===== Database Migration =====
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.User{},
		&model.Attachment{},
	); err != nil {
		log.Printf("Database migration failed: %v", err)
	}
	return nil
}
