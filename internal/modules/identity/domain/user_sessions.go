package domain

import (
	"time"

	"github.com/google/uuid"
)

type  User_Sessions struct {
 session_id    uuid.UUID `gorm:"type:uuid;primaryKey;"`
 user_id       uuid.UUID `gorm:"type:uuid;not null;"`
 device_name   string    `gorm:"type:varchar(255);"`
 os_version    string    `gorm:"type:varchar(50);"`
 browser       string    `gorm:"type:varchar(50);"`
 ip_address    string    `gorm:"type:varchar(45);"`
 location_city    string    `gorm:"type:varchar(100);"`
 location_country string    `gorm:"type:varchar(100);"`
 refresh_token string    `gorm:"type:text;"`
 is_active     bool      `gorm:"default:true;"`
 last_active_at time.Time
 CreatedAt     time.Time
}