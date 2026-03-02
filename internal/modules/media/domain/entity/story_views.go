package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/gocql/gocql"
)

const (
	TableStoryViews = "story_views"
)

// StoryView đại diện cho bảng 'story_views' trong Cassandra.
// Bảng này tối ưu cho việc ghi (Write Heavy) và đọc danh sách người xem (Read List).
type StoryView struct {
	// 1. PRIMARY KEY COMPONENT 1: PARTITION KEY
	// Gom tất cả lượt xem của 1 Story vào 1 Partition Node
	StoryID gocql.UUID `cql:"story_id" json:"story_id"`

	// 2. PRIMARY KEY COMPONENT 2: CLUSTERING KEY
	// Sắp xếp thời gian xem mới nhất lên đầu (DESC)
	ViewedAt time.Time `cql:"viewed_at" json:"viewed_at"`

	// 3. PRIMARY KEY COMPONENT 3: CLUSTERING KEY
	// Đảm bảo tính duy nhất: 1 User + 1 Thời điểm -> 1 Record
	ViewerID gocql.UUID `cql:"viewer_id" json:"viewer_id"`

	// 4. DENORMALIZATION (Snapshot Info)
	// Lưu chết thông tin user tại thời điểm xem.
	// Nếu user đổi avatar sau này, log cũ vẫn giữ avatar cũ (Chấp nhận được với Story)
	ViewerName      string `cql:"viewer_name" json:"viewer_name"`
	ViewerAvatarURL string `cql:"viewer_avatar_url" json:"viewer_avatar_url"`

	// 5. INTERACTION DETAILS
	// Enum: view, reaction, poll_vote
	InteractionType sharedEnums.StoryInteractionType `cql:"interaction_type" json:"interaction_type"`

	// Reaction Code: "❤️", "😂" hoặc ID sticker
	ReactionCode sharedEnums.ReactionCode `cql:"reaction_code" json:"reaction_code"`

	Content string `cql:"content" json:"content"`

	// Poll Vote: Dùng pointer (*int)
	// - nil: Không vote
	// - 0: Vote option đầu tiên
	// - 1: Vote option thứ hai
	PollOptionIndex *int `cql:"poll_option_index" json:"poll_option_index"`
}

// TableName trả về tên bảng trong Cassandra
func (StoryView) TableName() string {
	return TableStoryViews
}
