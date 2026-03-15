package IRepositoryShare

import (
	"context"
	"io"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
)

type ISeaweedfs interface {
	// --- BASIC STORAGE OPERATIONS ---

	// Upload: Upload file hoàn chỉnh (Avatar, Post, CV)
	Upload(ctx context.Context, input *dto.FileUploadInput) (*dto.FileUploadOutput, error)

	// Download: Lấy nội dung file (Dùng khi cần xử lý ảnh/video ở Backend)
	Download(ctx context.Context, filePath string) (io.ReadCloser, error)

	// Delete: Xóa file đơn lẻ (Ví dụ: Xóa ảnh cũ khi đổi avatar)
	Delete(ctx context.Context, filePath string) error

	// DeleteFolder: Xóa toàn bộ thư mục (Dùng khi Hard Delete User/Group)
	DeleteFolder(ctx context.Context, folderPath string) error

	// Exists: Kiểm tra file đã tồn tại chưa (Tránh upload trùng)
	Exists(ctx context.Context, filePath string) (bool, error)

	// --- LIVE STREAM & REALTIME OPERATIONS ---

	// UploadStreamSegment: Upload từng mảnh video (ts/m4s) của Live Stream
	// Sử dụng cho cơ chế HLS/DASH để đạt hiệu suất realtime
	UploadStreamSegment(ctx context.Context, OwnerID string, sessionID string, segmentName string, content io.Reader) error

	MergeToMP4(ctx context.Context, OwnerID string, sessionID string, outputPath string) error

	// UpdateStreamManifest: Cập nhật file chỉ mục (.m3u8 hoặc .mpd)
	// Để trình phát (Player) biết segment nào mới nhất để load
	UpdateStreamManifest(ctx context.Context, OwnerID string, sessionID string, manifestContent []byte) error

	GenerateVOD(ctx context.Context, OwnerID string, sessionID string, name string) (*dto.FileUploadOutput, error)

	// --- HELPER METHODS ---

	// GetPublicURL: Trả về URL ảnh/video để hiển thị trên UI
	GetPublicURL(filePath string) string

	GetStreamSegmentURL(ownerID string, sessionID string, segmentName string, Storage utils.StorageType) string

	GetStreamURLDIR(ownerID string, sessionID string, Storage utils.StorageType) string

	// GetStreamURL: Trả về URL của file manifest để xem Live Stream
	GetStreamURL(ownerID string, sessionID string, Storage utils.StorageType) string

	// GetUploadPresignedUrl: Lấy URL tạm thời để upload trực tiếp từ Frontend đến SeaweedFS
	GetUploadPresignedUrl(ctx context.Context, input *dto.FileUploadInput) (*dto.PresignedURLResponse, error)

	// GetDownloadPresignedUrl: Sinh link tải file trực tiếp với tốc độ cực cao, bypass Backend
	GetDownloadPresignedUrl(ctx context.Context, filePath string, forceDownload bool) (string, error)
}
