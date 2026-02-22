package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/enum"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	CollectionAds = "ads"
)

// Ad đại diện cho bảng 'ads' trong Postgres.
// Một Chiến dịch (Campaign) có thể chứa nhiều Quảng cáo (Ads/Ad Sets).
type Ad struct {
	// 1. PRIMARY KEY
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`

	// 2. RELATIONS
	// Thuộc về chiến dịch nào?
	// Index: { campaign_id: 1 } -> Lấy danh sách quảng cáo của chiến dịch
	CampaignID uuid.UUID `gorm:"type:uuid;not null;index" json:"campaign_id"`

	// 3. CROSS-DB LINKING (Postgres -> Mongo)
	// ID bài viết bên MongoDB muốn chạy quảng cáo.
	// Lưu dưới dạng String (ObjectId hex).
	TargetPostID string `gorm:"type:varchar(50);not null" json:"target_post_id"`

	// 4. BIDDING (Đấu thầu)
	// Giá thầu (VD: 5000 VNĐ cho 1 click)
	BidAmount float64 `gorm:"type:decimal(10,2);not null" json:"bid_amount"`

	// 5. STATUS & MODERATION
	// Mặc định là 'reviewing' khi vừa tạo
	Status enum.AdStatus `gorm:"type:varchar(20);index;default:'reviewing'" json:"status"`

	// Lý do từ chối (Chỉ có giá trị khi Status = Rejected)
	// Dùng Pointer (*string) để cho phép NULL trong DB.
	RejectionReason *string `gorm:"type:text" json:"rejection_reason,omitempty"`

	// 6. TIMESTAMPS
	// (Nên thêm để tracking thời gian tạo/sửa quảng cáo)
	CreatedAt time.Time      `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamp;default:now()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName ghi đè tên bảng
func (Ad) TableName() string {
	return CollectionAds
}
