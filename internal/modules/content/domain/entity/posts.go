package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	collectionnamepost = "Post"
)

type Post struct {
	// 1. ĐỊNH DANH
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// UserID: Đây là UUID từ Postgres (Identity Module).
	// Bắt buộc lưu String. Không dùng ObjectId vì không tương thích.
	UserID string `bson:"user_id" json:"user_id"`

	// 2. PHÂN LOẠI & NGỮ CẢNH
	Type    sharedEnums.PostType `bson:"type" json:"type"`
	Context *PostContext         `bson:"context,omitempty" json:"context,omitempty"`

	// 3. NỘI DUNG
	Content string       `bson:"content" json:"content"`
	Slug    string       `bson:"slug" json:"slug"` // Dùng để share link (SEO)
	Summary *PostSummary `bson:"summary,omitempty" json:"summary,omitempty"`

	// 4. LOGIC & TRẠNG THÁI
	// Privacy nên khởi tạo mặc định, không để null
	Privacy  PostPrivacy                  `bson:"privacy" json:"privacy"`
	Status   sharedEnums.ProcessingStatus `bson:"status" json:"status"`
	IsPinned bool                         `bson:"is_pinned" json:"is_pinned"`
	IsEdited bool                         `bson:"is_edited" json:"is_edited"`
	IsShared bool                         `bson:"is_shared" json:"is_shared"` // Nếu là bài viết được share lại từ bài khác thì true

	// 5. COUNTERS (Thường xuyên update -> Tách struct giúp code rõ ràng)
	Stats PostStats `bson:"stats" json:"stats"`

	// 6. SEARCH & TAGGING
	// Lưu mảng String cho Hashtag và Mention (Mention là UserID - UUID)
	Hashtags []string `bson:"hashtags,omitempty" json:"hashtags,omitempty"`
	Mentions []string `bson:"mentions,omitempty" json:"mentions,omitempty"`

	// 7. TIMESTAMPS
	PublishedAt *time.Time `bson:"published_at,omitempty" json:"published_at,omitempty"` // Pointer vì có thể là Draft (chưa publish)
	CreatedAt   time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// --- 1. CONTEXT (Ngữ cảnh bài viết) ---
type PostContext struct {
	Type sharedEnums.ContextType `bson:"type" json:"type"`

	// TargetID có thể là GroupID (ObjectId) hoặc UserID (UUID) tùy context.
	// Lưu String là an toàn nhất để chứa cả 2 loại.
	TargetID string `bson:"target_id,omitempty" json:"target_id,omitempty"`
}

// --- 2. SUMMARY (Tóm tắt hiển thị) ---
type PostSummary struct {
	FeelingIcon  string `bson:"feeling_icon,omitempty" json:"feeling_icon,omitempty"`   // "😄"
	FeelingName  string `bson:"feeling_name,omitempty" json:"feeling_name,omitempty"`   // "hạnh phúc"
	LocationName string `bson:"location_name,omitempty" json:"location_name,omitempty"` // "tại Starbucks"

	HasMedia     bool   `bson:"has_media" json:"has_media"`
	MediaCount   int    `bson:"media_count" json:"media_count"`
	ThumbnailURL string `bson:"thumbnail_url,omitempty" json:"thumbnail_url,omitempty"`

	// ID của background theme (nếu là post dạng text nền màu)
	BackgroundThemeID string `bson:"background_theme_id,omitempty" json:"background_theme_id,omitempty"`
}

// --- 3. PRIVACY (Quyền riêng tư) ---
type PostPrivacy struct {
	Scope        sharedEnums.PrivacyScope `bson:"scope" json:"scope"`
	AllowComment bool                     `bson:"allow_comment" json:"allow_comment"`
	AllowShare   bool                     `bson:"allow_share" json:"allow_share"`
}

// --- 4. STATS (Counters - Cache) ---
type PostStats struct {
	TotalReactions int `bson:"total_reactions" json:"total_reactions"`
	Comments       int `bson:"comments" json:"comments"`
	Shares         int `bson:"shares" json:"shares"`
	Views          int `bson:"views" json:"views"`
	Like           int `bson:"like" json:"like"`
	Love           int `bson:"love" json:"love"`
	Haha           int `bson:"haha" json:"haha"`
	Wow            int `bson:"wow" json:"wow"`
	Sad            int `bson:"sad" json:"sad"`
	Angry          int `bson:"angry" json:"angry"`

	// Cache top 2 reaction icon nhiều nhất để hiển thị (VD: ["👍", "❤️"])
	TopReactionTypes []string `bson:"top_reaction_types" json:"top_reaction_types"`
}

func (Post) CollectionNamePost() string {
	return collectionnamepost
}
