package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/gocql/gocql"
	// Thay your_module_name bằng tên module của project
)

type LiveCommentReq struct {
	StreamID      gocql.UUID               `json:"stream_id" validate:"required"`
	CreatedAt     time.Time                `json:"created_at"`
	CommentID     gocql.UUID               `json:"comment_id"`
	UserID        gocql.UUID               `json:"user_id" validate:"required"`
	UserNickname  string                   `json:"user_nickname" validate:"required"`
	UserAvatarURL string                   `json:"user_avatar_url"`
	UserBadges    []*sharedEnums.UserBadge `json:"user_badges"` // Sử dụng enum ở DTO
	Content       string                   `json:"content" validate:"required"`
	IsPinned      bool                     `json:"is_pinned"`
}
