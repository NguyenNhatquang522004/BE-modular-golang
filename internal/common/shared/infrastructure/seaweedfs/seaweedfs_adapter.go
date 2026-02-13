package seaweedfs

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/google/uuid"
	"github.com/linxGnu/goseaweedfs"
)

type SeaweedfsAdapter struct {
	client *goseaweedfs.Filer
}

func NewSeaweedfsAdapter(client *goseaweedfs.Filer) IRepositoryShare.ISeaweedfs {
	return &SeaweedfsAdapter{
		client: client,
	}
}

// Upload: Upload file hoàn chỉnh (Avatar, Post, CV)
func (r *SeaweedfsAdapter) Upload(ctx context.Context, input *dto.FileUploadInput) (*dto.FileUploadOutput, error) {
	// 1. TẠO TÊN FILE DUY NHẤT (UUID)
	// Tránh trùng lặp file và bảo mật tên file gốc của người dùng.
	ext := filepath.Ext(input.FileName)
	if ext == "" {
		ext = ".bin" // Fallback nếu không xác định được đuôi file
	}
	uniqueName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// 2. XÁC ĐỊNH THƯ MỤC LƯU TRỮ (FOLDER STRATEGY)
	// Nếu input.Folder trống, tự động sinh path dựa trên logic phân cấp Tầng 1-2-3 bạn đã định nghĩa.
	folderPath := input.Folder
	if folderPath == "" {
		folderPath = utils.GenerateStoragePath(input.OwnerID, input.Storage)
	}

	// 3. CHUẨN HÓA ĐƯỜNG DẪN ĐẦY ĐỦ (FULL PATH)
	// Đảm bảo đường dẫn bắt đầu bằng / và không có khoảng trắng thừa.
	cleanFolder := "/" + strings.Trim(folderPath, "/")
	fullPath := fmt.Sprintf("%s/%s", cleanFolder, uniqueName)

	// 4. THỰC HIỆN UPLOAD QUA FILER
	// Sử dụng io.Reader để stream trực tiếp, tối ưu RAM cho các file nặng như video/reels.
	// r.client là *goseaweedfs.Filer đã được khởi tạo trong database layer.
	result, err := r.client.Upload(input.Content, input.Size, fullPath, "", "")
	if err != nil {
		return nil, fmt.Errorf("seaweedfs_adapter: upload failed: %w", err)
	}

	// 5. TRẢ VỀ OUTPUT CHUẨN DTO
	// FilePath dùng để lưu vào MongoDB, PublicURL dùng để hiển thị lên UI.
	return &dto.FileUploadOutput{
		OwnerID:   input.OwnerID,
		Storage:   input.Storage,
		FilePath:  fullPath,
		PublicURL: r.GetPublicURL(fullPath),
		FileID:    result.FileID,
	}, nil
}

// Download: Lấy nội dung file (Dùng khi cần xử lý ảnh/video ở Backend)
func (r *SeaweedfsAdapter) Download(ctx context.Context, filePath string) (io.ReadCloser, error) {
	// Download implementation here
	return nil, nil
}

// Delete: Xóa file đơn lẻ (Ví dụ: Xóa ảnh cũ khi đổi avatar)
func (r *SeaweedfsAdapter) Delete(ctx context.Context, filePath string) error {
	// Delete implementation here
	return nil
}

// DeleteFolder: Xóa toàn bộ thư mục (Dùng khi Hard Delete User/Group)
func (r *SeaweedfsAdapter) DeleteFolder(ctx context.Context, folderPath string) error {
	// DeleteFolder implementation here
	return nil
}

// Exists: Kiểm tra file đã tồn tại chưa (Tránh upload trùng)
func (r *SeaweedfsAdapter) Exists(ctx context.Context, filePath string) (bool, error) {
	// Exists implementation here
	return false, nil
}

// --- LIVE STREAM & REALTIME OPERATIONS ---

// UploadStreamSegment: Upload từng mảnh video (ts/m4s) của Live Stream
// Sử dụng cho cơ chế HLS/DASH để đạt hiệu suất realtime
func (r *SeaweedfsAdapter) UploadStreamSegment(ctx context.Context, sessionID string, segmentName string, content io.Reader) error {
	// UploadStreamSegment implementation here
	return nil
}

// UpdateStreamManifest: Cập nhật file chỉ mục (.m3u8 hoặc .mpd)
// Để trình phát (Player) biết segment nào mới nhất để load
func (r *SeaweedfsAdapter) UpdateStreamManifest(ctx context.Context, sessionID string, manifestContent []byte) error {
	// UpdateStreamManifest implementation here
	return nil
}

// --- HELPER METHODS ---

// GetPublicURL: Trả về URL ảnh/video để hiển thị trên UI
func (r *SeaweedfsAdapter) GetPublicURL(filePath string) string {
	// GetPublicURL implementation here
	return ""
}

// GetStreamURL: Trả về URL của file manifest để xem Live Stream
func (r *SeaweedfsAdapter) GetStreamURL(sessionID string) string {
	// GetStreamURL implementation here
	return ""
}
