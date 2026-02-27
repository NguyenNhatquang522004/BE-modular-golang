package FFmpeg

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
)

type FFmpegAdapter struct {
}

func NewFFmpegAdapter() *FFmpegAdapter {
	return &FFmpegAdapter{}
}
func (f *FFmpegAdapter) StartTranscoding(ctx context.Context, config IRepositoryShare.TranscodeConfig) error {
	// 1. Chuẩn bị thư mục tạm (Temp folder)
	err := os.MkdirAll(config.OutputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	// Xóa thư mục tạm khi kết thúc phiên Live để dọn rác server
	defer os.RemoveAll(config.OutputDir)

	// Đường dẫn tới file m3u8
	manifestPath := filepath.Join(config.OutputDir, "index.m3u8")
	segmentPath := filepath.Join(config.OutputDir, "segment%04d.ts")

	// 2. Cấu hình FFmpeg Best Practice cho Livestream HLS
	args := []string{
		"-y", // Ghi đè file nếu đã tồn tại
		"-i", config.InputURL,

		// Cấu hình Video
		"-c:v", "libx264", // Dùng codec H.264 phổ biến nhất
		"-preset", "veryfast", // Tối ưu cho realtime, giảm tải CPU
		"-tune", "zerolatency", // Bắt buộc cho Livestream để giảm độ trễ
		"-profile:v", "main",
		"-g", "60", // Keyframe interval (ví dụ: 30fps * 2s = 60)

		// Cấu hình Audio
		"-c:a", "aac",
		"-ar", "44100",
		"-b:a", "128k",

		// Cấu hình băm HLS
		"-f", "hls",
		"-hls_time", fmt.Sprintf("%d", config.SegmentLen), // Độ dài 1 mảnh
		"-hls_list_size", "5", // Chỉ giữ 5 mảnh mới nhất trong m3u8 (cửa sổ live)
		"-hls_flags", "delete_segments+append_list", // FFmpeg tự xóa file .ts cũ ở local
		"-hls_segment_filename", segmentPath,
		manifestPath,
	}

	// 3. Khởi tạo Command với Context
	// Khi ctx.Done(), Go sẽ tự động gửi tín hiệu SIGKILL để tắt tiến trình FFmpeg này
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)

	// Bắt log của FFmpeg (FFmpeg thường in log ra Stderr)
	// Trong thực tế bạn nên redirect cái này vào một file log hoặc hệ thống logging (ELK/Loki)
	cmd.Stderr = os.Stderr

	log.Printf("Starting FFmpeg for session: %s", config.SessionID)

	// Chạy FFmpeg và chờ kết quả
	if err := cmd.Run(); err != nil {
		// Bỏ qua lỗi nếu do user chủ động cancel context (tắt live)
		if ctx.Err() == context.Canceled {
			log.Printf("Stream %s stopped gracefully by context", config.SessionID)
			return nil
		}
		return fmt.Errorf("ffmpeg transcoding failed: %w", err)
	}

	return nil
}

