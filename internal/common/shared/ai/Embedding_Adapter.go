package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
)

type EmbeddingAdapter struct {
	// Có thể thêm các trường cấu hình nếu cần, VD: ModelName, APIEndpoint, v.v.
	ollamaURL          string
	modelName          string
	modelEmbeddingName string
	cfg                *configs.Config
}

func NewEmbeddingAdapter(cfg *configs.Config) IRepositoryShare.IAI {
	ollamaURL := cfg.AIAnalysisMediaAsset.OllamaURL
	modelName := cfg.AIAnalysisMediaAsset.ModelName
	modelEmbeddingName := cfg.AIAnalysisMediaAsset.ModelEmbeddingName
	return &EmbeddingAdapter{
		ollamaURL:          ollamaURL,
		modelName:          modelName,
		modelEmbeddingName: modelEmbeddingName,
		cfg:                cfg,
	}
}

func (e *EmbeddingAdapter) GenerateBgeM3Embedding(text string) ([]float32, error) {
	url := e.ollamaURL + "/api/embeddings"

	// Khai báo đúng model bge-m3
	reqBody, _ := json.Marshal(map[string]string{
		"model":  e.modelEmbeddingName,
		"prompt": text,
	})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to call ollama: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Embedding []float32 `json:"embedding"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode embedding: %w", err)
	}

	return result.Embedding, nil
}
func (e *EmbeddingAdapter) CallOllamaModel(ctx context.Context, prompt string, imageBase64 []string, httpClient *http.Client) (string, error) {

	jsonValue, err := json.Marshal(mediaEvent.OllamaRequest{
		Model:  e.modelName,
		Prompt: prompt,
		Images: imageBase64,
		Stream: false,
	})
	if err != nil {
		return "", fmt.Errorf("lỗi marshal JSON: %w", err)
	}

	// Tạo request có kèm Context (để dễ dàng cancel/timeout từ phía trên)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.ollamaURL, bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", fmt.Errorf("lỗi tạo HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Gửi request
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("lỗi kết nối HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama trả về status code lỗi: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("lỗi đọc response body: %w", err)
	}

	var responseData mediaEvent.OllamaResponse
	if err := json.Unmarshal(body, &responseData); err != nil {
		return "", fmt.Errorf("lỗi unmarshal response: %w", err)
	}

	return responseData.Response, nil
}
