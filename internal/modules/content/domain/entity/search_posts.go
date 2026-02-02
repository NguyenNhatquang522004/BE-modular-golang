package entity

import (
	"time"

	// Import Enum từ module Content để tái sử dụng -> BEST PRACTICE (DRY)
	contentEnum "github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/enum"
)
type (
	
)
// SearchPost đại diện cho document trong index "search_posts"
// Index này dùng để tìm kiếm Full-text và lọc nâng cao
type SearchPost struct {
	// 1. ĐỊNH DANH
	// Trong ES, _id là string. Ta map field này vào _id khi bulk index.
	ID string `json:"id"`

	// 2. NỘI DUNG & KEYWORDS
	// Dùng text analyzer cho tìm kiếm nội dung
	Content string `json:"content"`
	// Dùng keyword analyzer cho việc filter chính xác
	Hashtags []string `json:"hashtags"`

	// 3. QUAN HỆ (Lưu ID dạng String)
	AuthorID string `json:"author_id"`

	// Dùng Pointer (*string) và omitempty.
	// Lý do: Nếu bài viết ở tường nhà -> GroupID = null.
	// ES sẽ không index field này nếu nó null -> Tiết kiệm dung lượng index.
	GroupID *string `json:"group_id,omitempty"`
	PageID  *string `json:"page_id,omitempty"`

	// 4. PHÂN LOẠI (FILTERS)
	// Tái sử dụng Enum của Content.
	// Nhờ 'enumer -json', nó sẽ tự lưu là ["image", "video"] thay vì [0, 1]
	MediaTypes []contentEnum.MediaType `json:"media_types"`

	// Chỉ index bài public. Nhưng vẫn lưu field này để double-check nếu cần.
	Privacy contentEnum.PrivacyScope `json:"privacy"`

	// 5. TIMESTAMPS
	CreatedAt time.Time `json:"created_at"`

	// 6. RANKING METRICS (Dùng để sort)
	LikesCount    int `json:"likes_count"`
	CommentsCount int `json:"comments_count"`
}

// IndexName trả về tên index trong Elasticsearch
func (SearchPost) IndexName() string {
	return "search_posts"
}
