package req

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/enum"
)

// CreatePostEntityEditLogReq: Dùng để ghi lại một hành động sửa đổi
type CreatePostEntityEditLogReq struct {
	TargetCollection enum.TargetCollection `json:"target_collection" binding:"required"` // 'posts' or 'comments'
	TargetID         string                `json:"target_id" binding:"required,mongoId"` // ID của bài viết bị sửa
	
	Version          int                   `json:"version" binding:"required,min=1"`     // Version mới là số mấy
	EditorID         string                `json:"editor_id" binding:"required,uuid"`    // Ai là người sửa

	Diff             LogDiffReq            `json:"diff" binding:"required"`              // Chi tiết thay đổi

	// Thông tin audit thường được lấy từ Context (Middleware), nhưng DTO vẫn cần để truyền vào
	IPAddress        string                `json:"ip_address,omitempty"`
	UserAgent        string                `json:"user_agent,omitempty"`
}

// UpdatePostEntityEditLogReq: (Ít dùng) Dùng để patch lại log nếu ghi sai
type UpdatePostEntityEditLogReq struct {
	// Các field này cho phép sửa nếu cần thiết
	Version   *int        `json:"version,omitempty" binding:"omitempty,min=1"`
	Diff      *LogDiffReq `json:"diff,omitempty"`
	IPAddress string      `json:"ip_address,omitempty"`
	UserAgent string      `json:"user_agent,omitempty"`
}

// --- Nested Struct ---

type LogDiffReq struct {
	OldContent    string   `json:"old_content"`
	NewContent    string   `json:"new_content"`
	ChangedFields []string `json:"changed_fields" binding:"required,dive,min=1"` // Ít nhất phải có tên 1 trường bị thay đổi
}