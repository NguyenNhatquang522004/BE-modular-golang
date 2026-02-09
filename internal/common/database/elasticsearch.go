package database

import (
	"fmt"
	"log"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"github.com/elastic/go-elasticsearch/v8"
)

// 1. Tạo Wrapper Struct
func NewElasticClient(cfg *configs.Config) (*elasticsearch.Client, func(), error) {
	// 1. Config
	esCfg := elasticsearch.Config{
		Addresses: []string{cfg.ElasticDB.ES_ADDRESS},
		// Username: cfg.ElasticDB.User,
		// Password: cfg.ElasticDB.Password,
	}

	// 2. Create Client
	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating elastic client: %w", err)
	}

	// 3. Ping / Info Check
	res, err := client.Info()
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to elastic: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, nil, fmt.Errorf("elastic returned error: %s", res.String())
	}

	log.Println("✅ Connected to Elasticsearch successfully!")

	// 4. Cleanup Function
	// Elastic client quản lý connection pool qua HTTP, không có hàm Close() cứng.
	// Nhưng vẫn trả về func rỗng để Wire dễ xử lý đồng bộ.
	cleanup := func() {
		// Log để biết app đang tắt
		log.Println("⚠️ Closing Elasticsearch client (handled by http transport)")
	}

	return client, cleanup, nil
}
