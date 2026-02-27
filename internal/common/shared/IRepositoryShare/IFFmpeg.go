package IRepositoryShare

import "context"

type TranscodeConfig struct {
	SessionID  string
	InputURL   string // Ví dụ: rtmp://localhost:1935/live/stream_key
	OutputDir  string // Thư mục tạm local (VD: /tmp/lives/{sessionID})
	SegmentLen int    // Thời lượng mỗi file .ts (giây) - Thường là 2-4 giây cho Low Latency
}

type IMediaTranscoder interface {
	// StartTranscoding chạy quá trình băm video.
	// Dùng context.Context để có thể hủy (stop live) từ xa một cách an toàn.
	StartTranscoding(ctx context.Context, config TranscodeConfig) error
	// Thêm hàm gộp VOD
	// m3u8URL: Link file index.m3u8 trên SeaweedFS
	// outputPath: Nơi lưu file mp4 tạm thời trên server để chuẩn bị upload
	MergeToMP4(ctx context.Context, m3u8URL string, outputPath string) error
}
	