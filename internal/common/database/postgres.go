package database

import (
	"fmt"
	"log"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConnection struct {
	Postgres *gorm.DB
}
func NewPostgresConnection(cfg *configs.Config) (*PostgresConnection, error) {
	db := &PostgresConnection{}
	err := db.ConnectPostgres(cfg)
	if err != nil {
		return nil, err
	}
	return db, nil
}
func (p *PostgresConnection) ConnectPostgres(config *configs.Config) error {
	log.Printf("Connecting to Postgres at %s:%s with user %s to database %s",
		config.PostgresDB.Host,
		config.PostgresDB.Port,
		config.PostgresDB.User,
		config.PostgresDB.DBName,
	)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.PostgresDB.Host,
		config.PostgresDB.User,
		config.PostgresDB.Password,
		config.PostgresDB.DBName,
		config.PostgresDB.Port,
		config.PostgresDB.DB_SSLMODE,
		config.PostgresDB.DB_TimeZone)

	// Mở kết nối
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to Postgres: %w", err)
	}

	// Cấu hình Connection Pool
	sqlDB, err := db.DB()
	if err != nil {
		// Lưu ý: Đừng dùng log.Fatal ở đây, hãy return error để main xử lý
		return fmt.Errorf("failed to get generic database object: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)
	sqlDB.SetConnMaxIdleTime(time.Minute)

	log.Println("Connected to PostgreSQL successfully")

	// 3. Gán vào struct receiver (Quan trọng)
	p.Postgres = db

	return nil
}

// Hàm lấy DB instance (Helper)
func (p *PostgresConnection) GetDB() *gorm.DB {
	return p.Postgres
}
func ProvideGormDB(conn *PostgresConnection) *gorm.DB {
	return conn.GetDB()
}