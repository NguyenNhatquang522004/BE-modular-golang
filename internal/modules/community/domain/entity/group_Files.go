package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/enum"
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
	UploaderID string `bson:"uploader_id" json:"uploader_id"`

	// 2. FILE INFO
	FileName string        `bson:"file_name" json:"file_name"` // "Bao_cao.xlsx"
	FileType enum.FileType `bson:"file_type" json:"file_type"` // 'xlsx', 'pdf'...

	// Kích thước file (Bytes). Dùng int64 để an toàn với file lớn > 2GB
	FileSize int64 `bson:"file_size" json:"file_size"`

	// 3. STORAGE
	// ID file trong hệ thống lưu trữ (SeaweedFS / S3)
	StorageFileID string `bson:"seaweedfs_file_id" json:"seaweedfs_file_id"`

	// 4. STATS
	// Index: { group_id: 1, download_count: -1 } -> File tải nhiều nhất
	DownloadCount int `bson:"download_count" json:"download_count"`

	// 5. TIMESTAMPS
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
func (GroupFile) CollectionName() string {
	return CollectionGroupFiles
}