package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

// GroupFileRes - Trả về 100% dữ liệu của Entity cho Client
type GroupFileRes struct {
	ID            string                `json:"id"`       // ObjectID -> String
	GroupID       string                `json:"group_id"` // ObjectID -> String
	UploaderID    string                `json:"uploader_id"`
	FileName      string                `json:"file_name"`
	FileType      sharedEnums.MediaType `json:"file_type"`
	FileSize      int64                 `json:"file_size"`
	StorageFileID string                `json:"seaweedfs_file_id"`
	DownloadCount int                   `json:"download_count"`
	CreatedAt     time.Time             `json:"created_at"`
}
