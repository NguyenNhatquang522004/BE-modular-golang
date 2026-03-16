package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

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
	seaweedfsRepo IRepositoryShare.ISeaweedfs
	aiRepo        IRepositoryShare.IAI
}

func NewConsumerAIAnalysisMediaAsset(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, mediaRepo IRepositoryMongodb.IMediaAssetsRepository, seaweedfsRepo IRepositoryShare.ISeaweedfs, aiRepo IRepositoryShare.IAI) *ConsumerAIAnalysisMediaAsset {

	return &ConsumerAIAnalysisMediaAsset{
		events:        events,
		pool:          pool,
		redisRepo:     redisRepo,
		mediaRepo:     mediaRepo,
		httpClient:    &http.Client{Timeout: 30 * time.Second},
		seaweedfsRepo: seaweedfsRepo,
		aiRepo:        aiRepo,
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
					processErr = c.handlecreatedAIAnalysisMediaAsset(ctx, event)
				case constants.Deleted.String():
					processErr = c.handlecreatedAIAnalysisMediaAsset(ctx, event)
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
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse event payload: %w", err))
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

	// --- 1. Gọi Ollama Model VỚI PROMPT MỚI ---
	systemPrompt := fmt.Sprintf(`
Bạn là một hệ thống AI phân tích nội dung mạng xã hội (bao gồm Text và Hình ảnh) để trích xuất các Topic (chủ đề) cho cơ sở dữ liệu đồ thị.

Nội dung bài viết:
"""%s"""

Hướng dẫn phân loại:
1. ĐÁNH GIÁ TỔNG THỂ: Kết hợp ý nghĩa của cả text và ảnh.
2. LỌC RÁC: Nếu chỉ là cảm xúc cá nhân, selfie, check-in vô thưởng vô phạt -> BẮT BUỘC trả về Topic là "BỎ_QUA".
3. TRÍCH XUẤT: Nếu có kiến thức/sở thích rõ ràng, trích xuất 1 đến 5 từ khóa (dạng snake_case, tiếng Việt không dấu).

RÀNG BUỘC ĐẦU RA (QUAN TRỌNG NHẤT):
- BẠN CHỈ ĐƯỢC PHÉP TRẢ VỀ 1 MẢNG JSON OBJECT HỢP LỆ.
- Mỗi object phải chứa đúng 2 trường: "topic" (string) và "confidence_score" (float, từ 0.0 đến 1.0 thể hiện độ chắc chắn của bạn).
- KHÔNG giải thích. KHÔNG dùng markdown (không dùng `+"```json"+`).

Ví dụ 1 (Bài rác):
Input: "Cuối tuần lười biếng ra quán quen chill" + Ảnh ly cafe
Output: [{"topic": "BỎ_QUA", "confidence_score": 1.0}]

Ví dụ 2 (Bài hợp lệ):
Input: "Vừa setup xong server test thử Kafka với Golang, chạy mượt phết." + Ảnh code.
Output: [{"topic": "golang", "confidence_score": 0.95}, {"topic": "kafka", "confidence_score": 0.90}, {"topic": "backend", "confidence_score": 0.85}]
`, data.Content)

	result, err := c.aiRepo.CallOllamaModel(ctx, systemPrompt, imagesBase64, c.httpClient)
	if err != nil {
		return fmt.Errorf("failed to call Ollama model: %w", err)
	}

	// Dọn dẹp chuỗi trả về
	cleanResult := strings.TrimSpace(result)
	cleanResult = strings.TrimPrefix(cleanResult, "```json")
	cleanResult = strings.TrimPrefix(cleanResult, "```")
	cleanResult = strings.TrimSuffix(cleanResult, "```")
	cleanResult = strings.TrimSpace(cleanResult)

	// --- 2. STRUCT ĐỂ HỨNG DATA CÓ ĐIỂM SỐ ---
	type AITopicResult struct {
		Topic           string  `json:"topic"`
		ConfidenceScore float64 `json:"confidence_score"`
	}

	var aiResults []AITopicResult
	if err := json.Unmarshal([]byte(cleanResult), &aiResults); err != nil {
		return fmt.Errorf("không thể parse kết quả từ AI thành mảng JSON Object. Raw: %s, err: %w", result, err)
	}

	// --- 3. Xử lý logic BỎ_QUA ---
	if len(aiResults) == 1 && aiResults[0].Topic == "BỎ_QUA" {
		fmt.Println("Hệ thống xác nhận đây là bài post rác/cá nhân. Không trích xuất Topic.")
		return nil
	}

	// --- 4. Tách mảng tên Topic (dành cho Payload cũ) và In ra kết quả ---
	var topicNames []string
	fmt.Printf("Trích xuất thành công %d topics:\n", len(aiResults))
	for i, item := range aiResults {
		topicNames = append(topicNames, item.Topic)
		fmt.Printf("%d. %s (Độ tự tin: %.2f)\n", i+1, item.Topic, item.ConfidenceScore)

		// GỢI Ý CHO BẠN: Nếu ở đây bạn gọi thẳng GraphRepository thì sẽ như thế này:
		// c.graphRepo.UpsertTopicNode(ctx, &entity.TopicNode{Name: item.Topic})
		// c.graphRepo.LinkPostToTopic(ctx, data.PostID, item.Topic, item.ConfidenceScore)
	}

	// --- 5. Publish Event ---
	// Lưu ý: Hiện tại Payload của bạn (`Hashtags: &topics`) đang nhận mảng []string.
	// Tôi trích xuất mảng `topicNames` đưa vào đây để code của bạn không bị lỗi.
	for _, mediaID := range data.MediaID {
		payload := &mediaEvent.UpdateMediaAssetsPayload{
			MediaID:  mediaID,
			Hashtags: &topicNames,
		}
		err = c.events.Publish(ctx, constants.TopicMeiaAssetHandleMetadata.String(), payload.MediaID, constants.Updated.String(), payload)
		if err != nil {
			return fmt.Errorf("failed to publish media asset update: %w", err)
		}
	}

	payloadpost := &contentEvent.UpdatePostPayload{
		PostID:   data.PostID,
		Hashtags: &topicNames,
	}
	err = c.events.Publish(ctx, constants.TopicPost.String(), payloadpost.PostID, constants.Updated.String(), payloadpost)
	if err != nil {
		return fmt.Errorf("failed to publish post update: %w", err)
	}

	return nil
}

func (c *ConsumerAIAnalysisMediaAsset) ConsumerFailedAIAnalysisMediaAsset(ctx context.Context) error {

	return nil
}
