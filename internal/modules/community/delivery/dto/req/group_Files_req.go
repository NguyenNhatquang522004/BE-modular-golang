package req

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/enum"
)

// CreateGroupFileReq - Ánh xạ đủ 100% các trường cần thiết để tạo File
type CreateGroupFileReq struct {
	GroupID       string        `json:"group_id" binding:"required"`
	UploaderID    string        `json:"uploader_id" binding:"required"`
	FileName      string        `json:"file_name" binding:"required"`
	FileType      enum.FileType `json:"file_type" binding:"required"`
	FileSize      int64         `json:"file_size" binding:"required,gt=0"` // File size nên lớn hơn 0
	StorageFileID string        `json:"seaweedfs_file_id" binding:"required"`
	// DownloadCount và CreatedAt không nên để Client gửi lên khi Create mà do hệ thống tự gán giá trị mặc định (0 và time.Now)
}

// UpdateGroupFileReq - Sử dụng 100% con trỏ để cập nhật linh hoạt (Partial Update)
type UpdateGroupFileReq struct {
	GroupID       *string        `json:"group_id,omitempty"`
	UploaderID    *string        `json:"uploader_id,omitempty"`
	FileName      *string        `json:"file_name,omitempty"`
	FileType      *enum.FileType `json:"file_type,omitempty"`
	FileSize      *int64         `json:"file_size,omitempty"`
	StorageFileID *string        `json:"seaweedfs_file_id,omitempty"`
	DownloadCount *int           `json:"download_count,omitempty"` // Hữu ích nếu có API riêng để tăng lượt tải
}