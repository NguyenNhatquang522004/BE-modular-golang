package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/contentEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
)

type ConsumerAIAnalysisMediaAsset struct {
	// This struct can be used to implement the IConsumerAIAnalysisMediaAsset interface
	events        events.EventBus
	pool          IRepositoryShare.IWorkerPool
	redisRepo     IRepositoryShare.IRedis
	mediaRepo     IRepositoryMongodb.IMediaAssetsRepository
	httpClient    *http.Client
	ollamaURL     string
	modelName     string
	cfg           configs.Config
	seaweedfsRepo IRepositoryShare.ISeaweedfs
}

func NewConsumerAIAnalysisMediaAsset(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, mediaRepo IRepositoryMongodb.IMediaAssetsRepository, seaweedfsRepo IRepositoryShare.ISeaweedfs, cfg configs.Config) *ConsumerAIAnalysisMediaAsset {
	ollamaURL := cfg.AIAnalysisMediaAsset.OllamaURL
	modelName := cfg.AIAnalysisMediaAsset.ModelName
	return &ConsumerAIAnalysisMediaAsset{
		events:        events,
		pool:          pool,
		redisRepo:     redisRepo,
		mediaRepo:     mediaRepo,
		httpClient:    &http.Client{Timeout: 30 * time.Second},
		ollamaURL:     ollamaURL,
		modelName:     modelName,
		cfg:           cfg,
		seaweedfsRepo: seaweedfsRepo,
	}
}
func (c *ConsumerAIAnalysisMediaAsset) ConsumerAIAnalysisMediaAsset(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicAIAnalysisMediaAsset.String(), 100, time.Duration(20)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
		var wg sync.WaitGroup
		errchan := make(chan error, len(events))
		for _, event := range events {
			status, can, err := c.redisRepo.Lock(ctx, event.ID)
			if err != nil {
				errchan <- errors.New("failed to acquire lock for event " + event.ID + ": " + err.Error())
				continue
			}
			if !can {
				if status == constants.StatusProcessing {
					errchan <- errors.New("event " + event.ID + " is currently being processed by another worker. Skipping.")
				} else {
					errchan <- errors.New("event " + event.ID + " has already been processed with status " + status.String() + ". Skipping.")
				}
				continue
			}
			wg.Add(1)
			var processErr error
			err = c.pool.Run(ctx, func() {
				defer wg.Done()
				switch event.Type {
				case constants.Created.String():
					processErr = c.handlecreatedAIAnalysisMediaAsset(ctx, event)
				case constants.Updated.String():
					processErr = c.handleupdatedAIAnalysisMediaAsset(ctx, event)
				case constants.Deleted.String():
					processErr = c.handledeletedAIAnalysisMediaAsset(ctx, event)
				default:
					errchan <- errors.New("event " + event.ID + " has an unknown type " + event.Type + ". Skipping.")
					return
				}
			})
			if err != nil {
				errchan <- errors.New("failed to submit event " + event.ID + " to worker pool: " + err.Error())
				c.redisRepo.Unlock(ctx, event.ID)
				wg.Done()
			}
			if processErr != nil {
				errchan <- errors.New("failed to process event " + event.ID + ": " + processErr.Error())
				c.redisRepo.Unlock(ctx, event.ID)
			} else {
				errchan <- nil
				c.redisRepo.MarkCompleted(ctx, event.ID)
			}
		}
		wg.Wait()
		close(errchan)
		var finalErr error
		for err := range errchan {
			if err != nil {
				finalErr = errors.Join(finalErr, err)
			}
		}
		return finalErr
	})
	if err != nil {
		return err
	}
	return nil
}
func (c *ConsumerAIAnalysisMediaAsset) handlecreatedAIAnalysisMediaAsset(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.AIAnalysisMediaAssetPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse event payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	var filenames []string
	for _, media := range data.MediaID {
		datamedia, err := c.mediaRepo.GetMediaAssetByID(ctx, media)
		if err != nil {
			return fmt.Errorf("failed to get media asset by ID %s: %w", media, err)
		}
		if datamedia == nil {
			return fmt.Errorf("media asset with ID %s not found", media)
		}
		filenames = append(filenames, datamedia.StorageFileID)
	}
	imagesBase64, err := c.seaweedfsRepo.DownloadMultipleImagesAsBase64(ctx, filenames)
	if err != nil {
		return fmt.Errorf("failed to download images from SeaweedFS: %w", err)
	}
	// --- Gọi Ollama Model ---
	systemPrompt := fmt.Sprintf(`
		Bạn là một hệ thống AI phân tích nội dung mạng xã hội (bao gồm Text và Hình ảnh) để trích xuất các Topic (chủ đề) cho cơ sở dữ liệu đồ thị.

Nội dung bài viết:
"""%s"""

Hướng dẫn phân loại và trích xuất:
1. ĐÁNH GIÁ TỔNG THỂ: Kết hợp ý nghĩa của cả đoạn text và(các) hình ảnh đính kèm.
2. LỌC RÁC (SPAM/PERSONAL): Nếu tổng thể bài viết chỉ mang tính chất cá nhân, cảm xúc nhất thời, vô thưởng vô phạt (ví dụ: ảnh selfie, check-in đi cafe, than thở, cap thả thính) và KHÔNG mang giá trị thông tin, chuyên môn hay sở thích chung -> BẮT BUỘC trả về duy nhất mảng: ["BỎ_QUA"].
3. TRÍCH XUẤT TOPIC (INTEREST GRAPH): Nếu bài viết chứa thông tin, kiến thức, sở thích rõ ràng (ví dụ: lập trình, xe cộ, thể thao, ẩm thực, review sản phẩm), hãy trích xuất 1 đến 5 từ khóa chủ đề cốt lõi nhất.

Quy định format Topic:
- Là danh từ hoặc cụm từ ngắn.
- Chuyển thành chữ thường, không dấu tiếng Việt, thay dấu cách bằng gạch dưới (ví dụ: "cong_nghe", "golang", "review_sach", "bong_da").

RÀNG BUỘC ĐẦU RA (QUAN TRỌNG NHẤT):
- BẠN CHỈ ĐƯỢC PHÉP TRẢ VỀ 1 MẢNG JSON HỢP LỆ.
- KHÔNG giải thích. KHÔNG thêm bất kỳ từ ngữ nào khác. KHÔNG bọc trong markdown (không dùng `+"```json"+`).

Ví dụ 1:
Input text: "Cuối tuần lười biếng ra góc quán quen ngồi chill chill một chút" + Ảnh ly cafe
Output: ["BỎ_QUA"]

Ví dụ 2:
Input text: "Vừa setup xong con server test thử Kafka với Golang, chạy mượt phết anh em ạ." + Ảnh màn hình code.
Output: ["golang", "kafka", "backend", "devops"].
	`, data.Content)

	result, err := c.callOllamaModel(ctx, systemPrompt, imagesBase64) // Tạm thời chỉ gửi ảnh đầu tiên để phân tích
	if err != nil {
		return fmt.Errorf("failed to call Ollama model: %w", err)
	}
	cleanResult := strings.TrimSpace(result)
	cleanResult = strings.TrimPrefix(cleanResult, "```json")
	cleanResult = strings.TrimPrefix(cleanResult, "```")
	cleanResult = strings.TrimSuffix(cleanResult, "```")
	cleanResult = strings.TrimSpace(cleanResult)

	// 2. Ép kiểu (Unmarshal) từ chuỗi JSON sang mảng Go (Slice)
	var topics []string
	if err := json.Unmarshal([]byte(cleanResult), &topics); err != nil {
		// Log lại giá trị gốc để bạn dễ debug xem con AI đã "nói bậy" cái gì
		return fmt.Errorf("không thể parse kết quả từ AI thành mảng JSON. Raw response: %s, err: %w", result, err)
	}

	// 3. Xử lý logic nghiệp vụ với data lấy được
	if len(topics) == 1 && topics[0] == "BỎ_QUA" {
		fmt.Println("Hệ thống xác nhận đây là bài post rác/cá nhân. Không trích xuất Topic.")
		return nil // Hoặc cập nhật status bài viết rồi return
	}

	// Tới đây data của bạn đã là mảng []string chuẩn, bạn có thể loop qua nó
	fmt.Printf("Trích xuất thành công %d topics:\n", len(topics))

	for _, data := range data.MediaID {
		payload := &mediaEvent.UpdateMediaAssetsPayload{
			MediaID:  data,
			Hashtags: &topics,
		}
		err = c.events.Publish(ctx, constants.TopicMeiaAssetHandleMetadata.String(), payload.MediaID, constants.Updated.String(), payload)
		if err != nil {
			return fmt.Errorf("failed to publish media asset update event for media ID %s: %w", data, err)
		}
	}
	payloadpost := &contentEvent.UpdatePostPayload{
		PostID:   data.PostID,
		Hashtags: &topics,
	}
	err = c.events.Publish(ctx, constants.TopicPost.String(), payloadpost.PostID, constants.Updated.String(), payloadpost)
	if err != nil {
		return fmt.Errorf("failed to publish post update event for post ID %s: %w", data.PostID, err)
	}

	for i, topic := range topics {
		fmt.Printf("%d. %s\n", i+1, topic)
	}
	fmt.Printf("Ollama Model Response: %s\n", result)
	return nil
}
func (c *ConsumerAIAnalysisMediaAsset) callOllamaModel(ctx context.Context, prompt string, imageBase64 []string) (string, error) {

	jsonValue, err := json.Marshal(mediaEvent.OllamaRequest{
		Model:  c.modelName,
		Prompt: prompt,
		Images: imageBase64,
		Stream: false,
	})
	if err != nil {
		return "", fmt.Errorf("lỗi marshal JSON: %w", err)
	}

	// Tạo request có kèm Context (để dễ dàng cancel/timeout từ phía trên)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.ollamaURL, bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", fmt.Errorf("lỗi tạo HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Gửi request
	resp, err := c.httpClient.Do(req)
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
func (c *ConsumerAIAnalysisMediaAsset) handleupdatedAIAnalysisMediaAsset(ctx context.Context, event events.IntegrationEvent) error {

	return nil
}
func (c *ConsumerAIAnalysisMediaAsset) handledeletedAIAnalysisMediaAsset(ctx context.Context, event events.IntegrationEvent) error {

	return nil
}
func (c *ConsumerAIAnalysisMediaAsset) ConsumerFailedAIAnalysisMediaAsset(ctx context.Context) error {

	return nil
}
