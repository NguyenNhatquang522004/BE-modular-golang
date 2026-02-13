package database

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"github.com/linxGnu/goseaweedfs"
)

// NewSeaweedFSClient khởi tạo kết nối đến SeaweedFS Filer
func NewSeaweedFSClient(cfg *configs.Config) (*goseaweedfs.Filer, func(), error) {
	// 1. Config
	filerUrl := cfg.SEAWEEDFS.SEAWEEDFS_FILER_URL

	// 2. Create Client
	// Sử dụng http.Client với timeout để tránh treo app
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	client, err := goseaweedfs.NewFiler(
		filerUrl,
		httpClient,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing seaweedfs filer client: %w", err)
	}

	// 3. Ping / Info Check (Thủ công)
	// Vì thư viện không có hàm Ping, ta tự gửi 1 request nhẹ vào root của Filer để check health
	resp, err := httpClient.Head(filerUrl) // Dùng HEAD cho nhẹ
	if err != nil {
		// Thử lại bằng GET nếu HEAD bị chặn, hoặc log warning
		// Ở đây log warning để không crash app nếu storage chưa lên kịp
		log.Printf("⚠️ Warning: Could not connect to SeaweedFS Filer at %s: %v", filerUrl, err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			log.Println("✅ Connected to SeaweedFS Filer successfully!")
		} else {
			log.Printf("⚠️ Warning: SeaweedFS Filer responded with status: %d", resp.StatusCode)
		}
	}

	// 4. Cleanup Function
	// goseaweedfs.Filer không có hàm Close(), http.Client tự quản lý connection pool.
	cleanup := func() {
		log.Println("⚠️ SeaweedFS client cleanup (No-op)")
	}

	return client, cleanup, nil
}
