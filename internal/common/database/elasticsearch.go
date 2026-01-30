package database

import (
	"fmt"
	"log"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"github.com/elastic/go-elasticsearch/v8"
)

// 1. Tạo Wrapper Struct
type ElasticConnection struct {
	elastic *elasticsearch.Client
}

func NewElasticConnection(cfg *configs.Config) (*ElasticConnection, error) {
	db := &ElasticConnection{}
	err := db.ConnectElastic(&cfg.ElasticDB)
	if err != nil {
		return nil, err
	}
	return db, nil
}
// 2. Hàm khởi tạo kết nối
// Receiver đặt là 'el' (viết tắt của Elastic)
func (el *ElasticConnection) ConnectElastic(cfg *configs.ElasticsearchConfig) error {
	// Cấu hình Client
	// Thư viện này yêu cầu Addresses là một mảng []string
	esCfg := elasticsearch.Config{
		Addresses: []string{cfg.ES_ADDRESS},
		// Vì trong Docker bạn tắt security (xpack=false), nên không cần username/password/cert
		// Nếu sau này bật security, bạn thêm Username: "elastic", Password: "..." vào đây.
	}

	// Tạo Client
	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return fmt.Errorf("error creating elastic client: %w", err)
	}

	// 3. Ping thử server (Quan trọng)
	// Elasticsearch client mặc định không kết nối ngay, nên ta gọi hàm Info() để test xem server có sống không
	res, err := client.Info()
	if err != nil {
		return fmt.Errorf("error getting response from elastic: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error pinging elastic: %s", res.String())
	}

	log.Println("Connected to Elasticsearch successfully!")

	// 4. Gán client vào struct (Để dùng lại)
	el.elastic = client

	return nil
}

// 3. Hàm Getter
func (el *ElasticConnection) GetClient() *elasticsearch.Client {
	return el.elastic
}
