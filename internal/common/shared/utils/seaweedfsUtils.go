package utils

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type StorageType string

const (
	// User Assets
	BucketAvatar StorageType = "avatar"
	BucketCover  StorageType = "cover"
	BucketCV     StorageType = "cv"

	// Social Content (Chia theo ngày tháng)
	BucketPost  StorageType = "posts"
	BucketReel  StorageType = "reels"
	BucketStory StorageType = "stories"
	BucketLive  StorageType = "lives" // Video xem lại livestream

	// Group Assets
	BucketGroupAvatar StorageType = "groups/avatar"
	BucketGroupCover  StorageType = "groups/cover"
	BucketGroupFile   StorageType = "groups/files"
	BucketGroupStream StorageType = "groups/streams" // Dành cho file video livestream của nhóm, nếu có tính năng này trong tương lai

	//page assets
	BucketPageAvatar     StorageType = "pages/avatar"
	BucketPageCover      StorageType = "pages/cover"
	BucketPageLiveStream StorageType = "pages/streams" // Dành cho file video livestream của page, nếu có tính năng này trong tương lai
	// System
	BucketSystem StorageType = "system"
)

// Allowed MIME Types
var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

var AllowedVideoTypes = map[string]bool{
	"video/mp4":       true,
	"video/quicktime": true, // .mov
	"video/webm":      true,
}

var AllowedDocTypes = map[string]bool{
	"application/pdf":    true,
	"application/msword": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true, // .docx
}

func (r StorageType) String() string {
	return string(r)
}
func GenerateStoragePath(ownerID string, bucketType StorageType) string {
	// Base path: users/{uuid}/{type}
	// VD: users/550e8400.../avatar
	basePath := fmt.Sprintf("users/%s/%s", ownerID, bucketType.String())

	// Nhóm Group thì path khác một chút
	if bucketType == BucketGroupAvatar || bucketType == BucketGroupCover || bucketType == BucketGroupFile || bucketType == BucketGroupStream {
		basePath = fmt.Sprintf("%s/%s", bucketType.String(), ownerID) // groups/avatar/{group_id}
	}

	// Logic phân chia thư mục theo thời gian (Partitioning)
	// Chỉ áp dụng cho dữ liệu sinh ra liên tục (Posts, Reels, Messages)
	// Để tránh folder bị quá tải (Too many files in directory)
	switch bucketType {
	case BucketPost, BucketReel, BucketStory, BucketLive, BucketGroupFile:
		now := time.Now()
		// Thêm /YYYY/MM vào đuôi
		// VD: users/.../posts/2026/02
		return fmt.Sprintf("%s/%d/%02d", basePath, now.Year(), now.Month())
	}

	// Các loại tĩnh (Avatar, Cover, CV) không cần chia ngày tháng
	// User thay avatar mới thì xóa cái cũ hoặc ghi đè, số lượng ít.
	return basePath
}
func ValidateFileCheck(content io.ReadSeeker, fileName string, maxSize int64, allowedTypes map[string]bool) (string, error) {
	// 1. Check Extension (Bước lọc sơ bộ)
	ext := strings.ToLower(filepath.Ext(fileName))
	if ext == "" {
		return "", fmt.Errorf("file extension is missing")
	}

	// 2. Check Magic Bytes (Quan trọng nhất)
	// Đọc 512 bytes đầu tiên để xem file thực sự là gì
	buffer := make([]byte, 512)
	_, err := content.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read file header: %w", err)
	}

	// Detect Content-Type thật từ buffer
	realContentType := http.DetectContentType(buffer)

	// Reset con trỏ file về đầu (0) để lát nữa Upload còn đọc được
	// Nếu không reset, file upload lên sẽ bị mất 512 bytes đầu -> File lỗi
	_, err = content.Seek(0, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("failed to reset file pointer: %w", err)
	}

	// 3. Validate Content Type
	// realContentType có thể trả về "image/jpeg; charset=utf-8", cần cắt bỏ charset
	mimeType := strings.Split(realContentType, ";")[0]

	if allowedTypes != nil {
		if !allowedTypes[mimeType] {
			return "", fmt.Errorf("file type not allowed: detected %s", mimeType)
		}
	}

	// 4. Check Size (Nếu cần check chính xác stream size, nhưng thường check từ Header ở Controller rồi)
	// Logic check size nên làm ở Middleware hoặc Controller để reject sớm.

	return mimeType, nil
}

// IsImage: Helper nhanh
func IsImage(contentType string) bool {
	return strings.HasPrefix(contentType, "image/")
}

// IsVideo: Helper nhanh
func IsVideo(contentType string) bool {
	return strings.HasPrefix(contentType, "video/")
}
