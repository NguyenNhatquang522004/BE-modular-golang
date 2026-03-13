package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionGroupFiles = "GroupFiles"
)

type GroupFile struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// 1. LINKING
	// File thuộc về nhóm nào (Mongo ID)
	// Index: { group_id: 1, created_at: -1 } -> List file mới nhất trong nhóm
	GroupID primitive.ObjectID `bson:"group_id" json:"group_id"`

	// Người upload (Postgres UUID -> String)
	UploaderID string              `bson:"uploader_id" json:"uploader_id"`
	PostID     *primitive.ObjectID `bson:"post_id,omitempty" json:"post_id,omitempty"`
	CommentID  *primitive.ObjectID `bson:"comment_id,omitempty" json:"comment_id,omitempty"`
	MessageID  *primitive.ObjectID `bson:"message_id,omitempty" json:"message_id,omitempty"`
	// 2. FILE INFO
	FileName string                `bson:"file_name" json:"file_name"` // "Bao_cao.xlsx"
	FileType sharedEnums.MediaType `bson:"file_type" json:"file_type"` // 'xlsx', 'pdf'...

	// Kích thước file (Bytes). Dùng int64 để an toàn với file lớn > 2GB
	FileSize int64 `bson:"file_size" json:"file_size"`

	// 3. STORAGE
	// ID file trong hệ thống lưu trữ (SeaweedFS / S3)
	StorageFileID string `bson:"seaweedfs_file_id" json:"seaweedfs_file_id"`
	OriginalURL   string `bson:"original_url" json:"original_url"`                       // Link tải file gốc
	ThumbnailURL  string `bson:"thumbnail_url,omitempty" json:"thumbnail_url,omitempty"` // Link ảnh thumbnail (nếu có, VD: PDF có thể có thumbnail)
	// 4. STATS
	// Index: { group_id: 1, download_count: -1 } -> File tải nhiều nhất
	DownloadCount int          `bson:"download_count" json:"download_count"`
	Metadata      FileMetadata `bson:"metadata" json:"metadata"`
	Caption       string       `bson:"caption,omitempty" json:"caption,omitempty"` // Mô tả file
	// 5. TIMESTAMPS
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}
type FileMetadata struct {
	Extension string `bson:"extension" json:"extension"`   // VD: ".xlsx", ".pdf"
	MimeType  string `bson:"mime_type" json:"mime_type"`   // VD: "application/vnd.ms-excel"
	SizeBytes int64  `bson:"size_bytes" json:"size_bytes"` // Dùng int64 cho file > 2GB

	// Các trường dưới đây dùng omitempty vì chỉ có ý nghĩa nếu File là Image/Video
	Width    *int     `bson:"width,omitempty" json:"width,omitempty"`
	Height   *int     `bson:"height,omitempty" json:"height,omitempty"`
	Duration *float64 `bson:"duration,omitempty" json:"duration,omitempty"` // Seconds
}

func (GroupFile) CollectionName() string {
	return CollectionGroupFiles
}
