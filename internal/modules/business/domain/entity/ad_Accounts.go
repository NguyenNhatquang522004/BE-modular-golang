package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/enum"
	"github.com/google/uuid"
)
const (
	CollectionAdAccounts = "ad_accounts"
)
// AdAccount đại diện cho bảng 'ad_accounts' trong Postgres.
type AdAccount struct {
	// 1. PRIMARY KEY
	// Dùng uuid.UUID cho Postgres là chuẩn nhất
	AccountID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"account_id"`

	// 2. OWNERSHIP
	// Khóa ngoại tham chiếu tới bảng Users
	OwnerUserID uuid.UUID `gorm:"type:uuid;not null;index" json:"owner_user_id"`

	// 3. SETTINGS
	// Mặc định là 'VND'
	Currency string `gorm:"type:varchar(3);default:'VND'" json:"currency"`
	Timezone string `gorm:"type:varchar(50)" json:"timezone"`

	// 4. FINANCIALS (DECIMAL)
	// Balance: Số dư (trả trước)
	Balance float64 `gorm:"type:decimal(15,2);default:0" json:"balance"`

	// CreditLimit: Hạn mức tín dụng (trả sau)
	CreditLimit float64 `gorm:"type:decimal(15,2);default:0" json:"credit_limit"`

	// 5. STATUS (ENUM)
	// Nhờ implement interface Valuer/Scanner ở file hooks, GORM tự động map string <-> int
	Status enum.AccountStatus `gorm:"type:varchar(20);index" json:"status"`

	// 6. TIMESTAMPS
	CreatedAt time.Time `gorm:"type:timestamp;default:now()" json:"created_at"`
}

func (AdAccount) TableName() string {
	return CollectionAdAccounts
}
