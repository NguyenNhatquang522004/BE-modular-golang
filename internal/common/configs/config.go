package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	DB_SSLMODE string
	DB_TimeZone string
}

type ServerConfig struct {
	Port    string
	BaseURL string
}
func LoadConfig() (*Config, error) {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 4. MANUAL MAPPING (Quan trọng)
	// Vì file .env của bạn đặt tên biến lộn xộn (Host, POSTGRES_USER...)
	// nên ta phải lấy từng cái bỏ vào đúng chỗ trong Struct.
	cfg := &Config{
		Database: DatabaseConfig{
			Host:     viper.GetString("Host"),              // Lấy key "Host"
			Port:     viper.GetString("PORT"),              // Lấy key "PORT"
			User:     viper.GetString("POSTGRES_USER"),     // Lấy key "POSTGRES_USER"
			Password: viper.GetString("POSTGRES_PASSWORD"), // Lấy key "POSTGRES_PASSWORD"
			DBName:   viper.GetString("POSTGRES_DB"),       // Lấy key "POSTGRES_DB"
			DB_SSLMODE: viper.GetString("DB_SSLMODE"),
			DB_TimeZone: viper.GetString("DB_TIMEZONE"),
		},
		Server: ServerConfig{
			// Nếu trong env không có SERVER_PORT thì dùng mặc định 8080
			Port:    viper.GetString("SERVER_PORT"),
			BaseURL: viper.GetString("Base_URL"),
		},
	}

	return cfg, nil
}

