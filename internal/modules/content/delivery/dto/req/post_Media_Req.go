package req

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/enum"

// CreatePostMediaReq: Dùng khi tạo mới (thường là upload xong thì gọi API này để lưu metadata)
type CreatePostMediaReq struct {
	PostID string          `json:"post_id" binding:"required,mongoId"`
	Items  []*MediaItemReq `json:"items" binding:"required,dive"` // Validate từng phần tử trong mảng
}

// UpdatePostMediaReq: Dùng khi muốn sửa danh sách ảnh (xóa ảnh, thêm ảnh, tag thêm người)
type UpdatePostMediaReq struct {
	// Chỉ cần gửi danh sách Items mới. Logic BE sẽ thay thế list cũ bằng list này.
	Items []*MediaItemReq `json:"items" binding:"required,dive"`
}

// --- Nested Structs (Dùng chung cho cả Create và Update) ---

type MediaItemReq struct {
	// ID có thể để trống khi tạo mới. Khi update nếu có ID thì giữ nguyên, không có thì tạo mới.
	ID string `json:"id,omitempty" binding:"omitempty,mongoId"`

	MediaType    enum.MediaType `json:"media_type" binding:"required,oneof=image video"` // Validate enum
	URL          string         `json:"url" binding:"required,url"`
	ThumbnailURL string         `json:"thumbnail_url" binding:"omitempty,url"`

	Metadata MediaMetadataReq `json:"metadata"`
	Order    int              `json:"order" binding:"min=0"`

	TaggedUsers []TaggedUserReq `json:"tagged_users,omitempty" binding:"omitempty,dive"`
}

type MediaMetadataReq struct {
	Width     int     `json:"width,omitempty" binding:"min=0"`
	Height    int     `json:"height,omitempty" binding:"min=0"`
	Duration  float64 `json:"duration,omitempty" binding:"min=0"`
	SizeBytes int64   `json:"size_bytes" binding:"min=0"`
	MimeType  string  `json:"mime_type" binding:"required"` // vd: image/jpeg
}

type TaggedUserReq struct {
	UserID string  `json:"user_id" binding:"required,uuid"` // Validate UUID Postgres
	Name   string  `json:"name" binding:"required"`
	X      float64 `json:"x" binding:"min=0,max=1"` // Tọa độ phải nằm trong ảnh (0 -> 1)
	Y      float64 `json:"y" binding:"min=0,max=1"`
}
