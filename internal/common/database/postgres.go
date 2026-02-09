package database

import (
	"fmt"
	"log"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(cfg *configs.Config) (*gorm.DB, func(), error) {
	conf := cfg.PostgresDB

	log.Printf("Connecting to Postgres at %s:%s...", conf.Host, conf.Port)

	// 1. DSN String
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		conf.Host, conf.User, conf.Password, conf.DBName, conf.Port, conf.DB_SSLMODE, conf.DB_TimeZone)

	// 2. Open Connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to Postgres: %w", err)
	}

	// 3. Config Connection Pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)
	sqlDB.SetConnMaxIdleTime(time.Minute)

	// Ping thử phát cho chắc
	if err := sqlDB.Ping(); err != nil {
		return nil, nil, fmt.Errorf("failed to ping Postgres: %w", err)
	}

	log.Println("✅ Connected to PostgreSQL successfully!")

	// 4. Cleanup Function
	cleanup := func() {
		log.Println("⚠️ Closing PostgreSQL connection...")
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing postgres: %v", err)
		}
	}

	return db, cleanup, nil
}
