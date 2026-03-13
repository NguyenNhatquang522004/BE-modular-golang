package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
)

type ConsumerMeiaAssetHandleMetadata struct {
	events    events.EventBus
	pool      IRepositoryShare.IWorkerPool
	redisRepo IRepositoryShare.IRedis
	mediaRepo IRepositoryMongodb.IMediaAssetsRepository
	storage   IRepositoryShare.ISeaweedfs
}

func NewConsumerMeiaAssetHandleMetadata(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, mediaRepo IRepositoryMongodb.IMediaAssetsRepository, storage IRepositoryShare.ISeaweedfs) *ConsumerMeiaAssetHandleMetadata {
	return &ConsumerMeiaAssetHandleMetadata{
		events:    events,
		pool:      pool,
		redisRepo: redisRepo,
		mediaRepo: mediaRepo,
		storage:   storage,
	}
}
func (c *ConsumerMeiaAssetHandleMetadata) ConsumeMediaAsset(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicMediaAsset.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
		var wg sync.WaitGroup
		errchan := make(chan error, len(events))

		for _, event := range events {
			// 1. LẤY LOCK (Idempotency)
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
			ev := event // Tránh lỗi tham chiếu biến trong vòng lặp của Goroutine

			// 2. CHẠY WORKER POOL (Đã fix lỗi Race Condition)
			poolErr := c.pool.Run(ctx, func() {
				defer wg.Done()

				var processErr error
				switch ev.Type {
				case constants.Created.String(): // SỰ KIỆN MỚI CHO CÁCH 2
					processErr = c.handlerProcessMediaAsset(ctx, ev)
				case constants.Updated.String():
					processErr = c.handlerProcessMediaAsset(ctx, ev)
				default:
					processErr = fmt.Errorf("unknown event type: %s", ev.Type)
				}

				// Xử lý Lock và Error NGAY TRONG GOROUTINE để đảm bảo đồng bộ
				if processErr != nil {
					errchan <- fmt.Errorf("failed to process event %s: %w", ev.ID, processErr)
					c.redisRepo.Unlock(ctx, ev.ID) // Mở khóa để cho phép Retry
				} else {
					c.redisRepo.MarkCompleted(ctx, ev.ID)
					errchan <- nil
				}
			})

			// Lỗi khi xin slot từ Worker Pool (Quá tải)
			if poolErr != nil {
				c.redisRepo.Unlock(ctx, ev.ID)
				errchan <- fmt.Errorf("worker pool full/error for event %s: %w", ev.ID, poolErr)
				wg.Done()
			}
		}

		wg.Wait()
		close(errchan)

		// Gộp lỗi
		var combinedErr error
		for err := range errchan {
			if err != nil {
				combinedErr = errors.Join(combinedErr, err)
			}
		}
		return combinedErr
	})
	return err
}

// =========================================================================
// CÁCH 2: HÀM XỬ LÝ BACKGROUND CHO MEDIA (CẮT THUMBNAIL, CẬP NHẬT DB)
// =========================================================================

func (c *ConsumerMeiaAssetHandleMetadata) handlerProcessMediaAsset(ctx context.Context, event events.IntegrationEvent) error {
	// 1. PARSE PAYLOAD & KÉO ENTITY TỪ MONGODB
	data, err := utils.ParsePayload[mediaEvent.ProcessMediaPayload](event.Payload)
	if err != nil || data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("invalid payload"))
	}

	mediaAsset, err := c.mediaRepo.GetMediaAssetByID(ctx, data.MediaID)
	if err != nil || mediaAsset == nil {
		return fmt.Errorf("media asset not found: %s", data.MediaID)
	}

	if !utils.IsVideo(mediaAsset.Metadata.MimeType) && !utils.IsImage(mediaAsset.Metadata.MimeType) {
		return nil // Không hỗ trợ xử lý file rác/doc ở luồng media này
	}

	// 2. TẢI FILE TỪ SEAWEEDFS VỀ LOCAL (/TMP)
	stream, err := c.storage.Download(ctx, mediaAsset.StorageFileID)
	if err != nil {
		return fmt.Errorf("failed to download from seaweedfs: %w", err)
	}
	defer stream.Close()

	ext := filepath.Ext(mediaAsset.StorageFileID)
	localRawPath := fmt.Sprintf("/tmp/%s_raw%s", mediaAsset.ID.Hex(), ext)
	localThumbPath := fmt.Sprintf("/tmp/%s_thumb.jpg", mediaAsset.ID.Hex())

	outFile, err := os.Create(localRawPath)
	if err != nil {
		return fmt.Errorf("failed to create local temp file: %w", err)
	}

	// Tính toán SizeBytes ngay trong lúc copy (tối ưu tốc độ)
	writtenBytes, err := io.Copy(outFile, stream)
	outFile.Close()
	defer os.Remove(localRawPath)
	defer os.Remove(localThumbPath)

	if err != nil {
		return fmt.Errorf("failed to write stream to disk: %w", err)
	}

	// Cập nhật SizeBytes
	mediaAsset.Metadata.SizeBytes = writtenBytes

	// 3. ĐỌC METADATA (WIDTH, HEIGHT, DURATION) BẰNG FFPROBE
	width, height, duration, err := c.extractMediaMetadata(ctx, localRawPath)
	if err != nil {
		// Chỉ log cảnh báo, không làm chết worker vì có thể file bị lỗi header nhẹ
		fmt.Printf("[Warning] Có lỗi khi đọc metadata file %s: %v\n", localRawPath, err)
	} else {
		mediaAsset.Metadata.Width = width
		mediaAsset.Metadata.Height = height
		mediaAsset.Metadata.Duration = duration
	}

	// 4. TẠO THUMBNAIL
	if utils.IsVideo(mediaAsset.Metadata.MimeType) {
		err = c.generateVideoThumbnail(ctx, localRawPath, localThumbPath)
	} else {
		err = c.generateImageThumbnail(ctx, localRawPath, localThumbPath)
	}

	if err != nil {
		return fmt.Errorf("failed to generate thumbnail: %w", err)
	}

	// 5. UPLOAD THUMBNAIL LÊN SEAWEEDFS
	thumbFile, err := os.Open(localThumbPath)
	if err != nil {
		return fmt.Errorf("failed to open generated thumbnail: %w", err)
	}
	defer thumbFile.Close()
	thumbInfo, _ := thumbFile.Stat()

	uploadInput := &dto.FileUploadInput{
		FileName:    fmt.Sprintf("%s_thumb.jpg", mediaAsset.ID.Hex()),
		OwnerID:     mediaAsset.UserID,
		Storage:     utils.BucketSystem,
		Folder:      fmt.Sprintf("/thumbnails/users/%s", mediaAsset.UserID),
		Content:     thumbFile,
		Size:        thumbInfo.Size(),
		ContentType: "image/jpeg",
	}

	uploadResult, err := c.storage.Upload(ctx, uploadInput)
	if err != nil {
		return fmt.Errorf("failed to upload thumbnail: %w", err)
	}

	// 6. CẬP NHẬT TOÀN BỘ DATA VÀO MONGODB
	mediaAsset.ThumbnailURL = uploadResult.PublicURL
	mediaAsset.UpdatedAt = time.Now()
	// Tùy chọn: Set Status = Ready nếu bạn có trường Status
	// mediaAsset.Status = sharedEnums.ProcessingStatusReady

	err = c.mediaRepo.UpdateMediaAsset(ctx, mediaAsset)
	if err != nil {
		_ = c.storage.Delete(ctx, uploadResult.FilePath) // Rollback
		return fmt.Errorf("failed to update db with metadata & thumbnail: %w", err)
	}

	// 7. (TÙY CHỌN) BẮN SỰ KIỆN QUA SOCKET ĐỂ FRONTEND RENDER ẢNH
	// socketMessage := socket.Message{ Type: "media_ready", Payload: mediaAsset.ID }
	// c.socketHub.SendToUser(mediaAsset.UserID, socketMessage)

	return nil
}

// --- HELPER FUNC: TẠO THUMBNAIL TỪ VIDEO ---
func (c *ConsumerMeiaAssetHandleMetadata) generateVideoThumbnail(ctx context.Context, inputPath, outputPath string) error {
	// Lệnh lấy frame tại giây thứ 1 (hoặc giữa video) làm ảnh thumb
	args := []string{
		"-y",            // Ghi đè nếu có
		"-i", inputPath, // File video đầu vào
		"-ss", "00:00:01.000", // Tại giây số 1
		"-vframes", "1", // Chỉ lấy 1 frame
		"-vf", "scale=-1:720", // Resize height về 720p cho nhẹ, giữ nguyên tỉ lệ width
		outputPath, // Đầu ra (jpg)
	}
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg error: %s", stderr.String())
	}
	return nil
}

// --- HELPER FUNC: TẠO THUMBNAIL TỪ ẢNH GỐC ---
func (c *ConsumerMeiaAssetHandleMetadata) generateImageThumbnail(ctx context.Context, inputPath, outputPath string) error {
	// Dùng ffmpeg hoặc ImageMagick để resize ảnh lớn thành ảnh nhỏ
	args := []string{
		"-y",
		"-i", inputPath,
		"-vf", "scale=-1:480", // Nén về 480p làm thumbnail
		outputPath,
	}
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

type FFProbeOutput struct {
	Streams []struct {
		Width     int    `json:"width"`
		Height    int    `json:"height"`
		Duration  string `json:"duration"` // ffprobe trả duration dạng chuỗi "12.345"
		CodecType string `json:"codec_type"`
	} `json:"streams"`
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
}

// extractMediaMetadata: Đọc Width, Height, Duration của Ảnh/Video
func (c *ConsumerMeiaAssetHandleMetadata) extractMediaMetadata(ctx context.Context, filePath string) (width int, height int, duration float64, err error) {
	args := []string{
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		filePath,
	}

	cmd := exec.CommandContext(ctx, "ffprobe", args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return 0, 0, 0, fmt.Errorf("ffprobe execution failed: %w", err)
	}

	var output FFProbeOutput
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse ffprobe json: %w", err)
	}

	// Ưu tiên tìm stream video hoặc hình ảnh đầu tiên để lấy Width/Height
	for _, stream := range output.Streams {
		if stream.CodecType == "video" || stream.CodecType == "image" {
			width = stream.Width
			height = stream.Height

			// Đôi khi video lưu duration ở trong stream
			if stream.Duration != "" {
				parsedDuration, _ := strconv.ParseFloat(stream.Duration, 64)
				if parsedDuration > duration {
					duration = parsedDuration
				}
			}
			break
		}
	}

	// Thông thường video format lưu duration chuẩn xác nhất ở root format
	if output.Format.Duration != "" {
		parsedDuration, _ := strconv.ParseFloat(output.Format.Duration, 64)
		if parsedDuration > 0 {
			duration = parsedDuration
		}
	}

	return width, height, duration, nil
}
