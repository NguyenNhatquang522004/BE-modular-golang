package database

import (
	"fmt"
	"log"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgres(config *configs.Config) (db *gorm.DB, err error) {
	log.Printf("Connecting to Postgres at %s:%s with user %s to database %s",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.DBName,
	)
	// Thêm logic kết nối PostgreSQL ở đây
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.Database.Host, config.Database.User, config.Database.Password, config.Database.DBName,
		config.Database.Port,
		config.Database.DB_SSLMODE,
		config.Database.DB_TimeZone)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Postgres: %w", err)
	}
	log.Println("Connected to PostgreSQL successfully")
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Error getting generic database object:", err)
	}

	// Configure connection pool settings
	sqlDB.SetMaxIdleConns(10)                  // Maximum number of connections in the idle pool
	sqlDB.SetMaxOpenConns(100)                 // Maximum number of open connections (idle + active)
	sqlDB.SetConnMaxLifetime(10 * time.Minute) // Maximum amount of time a connection may be reused
	sqlDB.SetConnMaxIdleTime(time.Minute)      // Maximum amount of time a connection can be idle	
	return db, nil

}
