package search

import (
	"context"
	"fmt"
	"strings"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/search/domain/entity"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type ModuleSearch struct {
	client *elasticsearch.Client
}

func NewModuleSearch(client *elasticsearch.Client) *ModuleSearch {
	m := &ModuleSearch{
		client: client,
	}
	err := m.InitElasticPostES()
	if err != nil {
		panic("Failed to init ElasticSearch for Search Module: " + err.Error())
	}
	return m
}

// InitElasticPostES: Khởi tạo Index và Mapping cho Search Module
func (m *ModuleSearch) InitElasticPostES() error {
	ctx := context.Background()
	indexName := entity.PostES{}.IndexName()

	// 1. Kiểm tra xem Index đã tồn tại chưa
	// Dùng esapi để gọi API check exists
	reqExists := esapi.IndicesExistsRequest{
		Index: []string{indexName},
	}
	resExists, err := reqExists.Do(ctx, m.client)
	if err != nil {
		return fmt.Errorf("check index exists error: %w", err)
	}
	defer resExists.Body.Close()

	// Nếu Index đã tồn tại (StatusCode 200) -> Bỏ qua (hoặc xử lý migration nếu cần)
	if resExists.StatusCode == 200 {
		return nil
	}

	// 2. Định nghĩa Mapping & Settings (JSON)
	// Đây là phần QUAN TRỌNG NHẤT
	mapping := `{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0,
			"analysis": {
				"analyzer": {
					"vietnamese_folding": {
						"tokenizer": "standard",
						"filter": ["lowercase", "asciifolding"]
					}
				}
			}
		},
		"mappings": {
			"properties": {
				"id": { "type": "keyword" },

				"content": {
					"type": "text",
					"analyzer": "vietnamese_folding",
					"search_analyzer": "vietnamese_folding"
				},

				"hashtags": { "type": "keyword" },
				"author_id": { "type": "keyword" },
				"group_id": { "type": "keyword" },
				"page_id": { "type": "keyword" },

				"media_types": { "type": "keyword" },
				"privacy": { "type": "keyword" },

				"created_at": { "type": "date" },

				"likes_count": { "type": "integer" },
				"comments_count": { "type": "integer" }
			}
		}
	}`

	// 3. Tạo Index
	reqCreate := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  strings.NewReader(mapping),
	}

	resCreate, err := reqCreate.Do(ctx, m.client)
	if err != nil {
		return fmt.Errorf("create index error: %w", err)
	}
	defer resCreate.Body.Close()

	if resCreate.IsError() {
		return fmt.Errorf("create index failed: %s", resCreate.String())
	}

	fmt.Printf(">>> Elastic Index [%s] initialized successfully with Vietnamese Analyzer\n", indexName)
	return nil
}
