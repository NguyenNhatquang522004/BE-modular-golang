package res

import (
	"time"

	"github.com/gocql/gocql"
)

type LiveCommentRes struct {
	StreamID      gocql.UUID `json:"stream_id"`
	CreatedAt     time.Time  `json:"created_at"`
	CommentID     gocql.UUID `json:"comment_id"`
	UserID        gocql.UUID `json:"user_id"`
	UserNickname  string     `json:"user_nickname"`
	UserAvatarURL string     `json:"user_avatar_url"`
	UserBadges    []string   `json:"user_badges"` // Sử dụng enum ở DTO
	Content       string     `json:"content"`
	IsPinned      bool       `json:"is_pinned"`
}
