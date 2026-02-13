package dto

import (
	"io"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
)

// FileUploadInput: Dữ liệu đầu vào khi muốn upload
type FileUploadInput struct {
	// FileName: Tên file gốc (ví dụ: "avatar.jpg").
	// Dùng để lấy extension và set Content-Type.
	FileName string
	// OwnerID: ID của entity sở hữu file (UserID, PostID, GroupID).
	OwnerID string
	// Storage: Loại storage (Avatar, Post, CV) để xác định folder lưu trữ.
	Storage utils.StorageType
	// Folder: Thư mục chứa file (ví dụ: "avatars", "posts/2024", "cvs").
	// Giúp tổ chức file gọn gàng trên SeaweedFS.
	Folder string
	// Content: Dùng io.Reader thay vì []byte để tối ưu RAM.
	// Cho phép stream trực tiếp từ HTTP Request Body sang SeaweedFS
	// mà không cần load toàn bộ file 500MB vào RAM.
	Content io.Reader

	// Size: Kích thước file (bytes). SeaweedFS cần biết size để cấp phát chunk.
	Size int64

	// ContentType: MIME type (image/jpeg, video/mp4).
	// Quan trọng để Browser hiển thị đúng thay vì download.
	ContentType string
}

// FileUploadOutput: Kết quả trả về sau khi upload
type FileUploadOutput struct {
	// OwnerID: ID của entity sở hữu file (UserID, PostID, GroupID).
	OwnerID string
	// Storage: Loại storage (Avatar, Post, CV) để xác định folder lưu trữ.
	Storage utils.StorageType
	// FilePath: Đường dẫn nội bộ trong SeaweedFS (ví dụ: "/avatars/abc-123.jpg")
	// Dùng để lưu vào DB.
	FilePath string
	// PublicURL: Đường dẫn đầy đủ để Frontend truy cập
	// (ví dụ: "https://cdn.mysocial.com/avatars/abc-123.jpg")
	PublicURL string

	// FileID: ID nội bộ của SeaweedFS (nếu cần dùng cho advanced logic)
	FileID string
}
