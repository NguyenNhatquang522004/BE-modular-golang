package consumer

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
	"github.com/fsnotify/fsnotify"
)

type ConsumerLiveSession struct {
	liveComment IRepositoryCassandra.ILiveCommentsRepository
	livesession IRepositoryMongodb.ILiveSessionRepository
	events      events.EventBus
	seaweedfs   IRepositoryShare.ISeaweedfs
	ffampeg     IRepositoryShare.IMediaTranscoder
	pool        IRepositoryShare.IWorkerPool
}

func NewConsumerLiveSession(liveComment IRepositoryCassandra.ILiveCommentsRepository, livesession IRepositoryMongodb.ILiveSessionRepository, events events.EventBus, seaweedfs IRepositoryShare.ISeaweedfs, ffampeg IRepositoryShare.IMediaTranscoder, pool IRepositoryShare.IWorkerPool) *ConsumerLiveSession {
	return &ConsumerLiveSession{
		liveComment: liveComment,
		livesession: livesession,
		events:      events,
		seaweedfs:   seaweedfs,
		ffampeg:     ffampeg,
		pool:        pool,
	}
}
func (c *ConsumerLiveSession) ConsumeStartLiveStream(ctx context.Context) {
	err := c.events.Subscribe(ctx, constants.TopicStartStopLive.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data, ok := event.Payload.(*mediaEvent.StartStopVideoLiveStreamPayload)
		if !ok {
			return fmt.Errorf("invalid event payload")
		}
		switch event.Type {
		case string(constants.Created):
			var outputdir string
			if data.PageID != "" {
				outputdir = c.seaweedfs.GetStreamURLDIR(data.OwnerID, data.LiveSessionID, utils.BucketPageLiveStream)
			}
			if data.GroupID != "" {
				outputdir = c.seaweedfs.GetStreamURLDIR(data.OwnerID, data.LiveSessionID, utils.BucketGroupStream)
			}
			if data.PageID == "" && data.GroupID == "" {
				outputdir = c.seaweedfs.GetStreamURLDIR(data.OwnerID, data.LiveSessionID, utils.BucketLive)
			}
			log.Printf("Stream URL DIR: %s", outputdir)
			go c.watchAndUploadSegments(ctx, data.OwnerID, data.LiveSessionID, outputdir)
			transcodeConfig := &IRepositoryShare.TranscodeConfig{
				SessionID:  data.LiveSessionID,
				OutputDir:  outputdir,
				SegmentLen: 5, // Ví dụ: 5 giây mỗi segment
				InputURL:   "",
			}

			err := c.ffampeg.StartTranscoding(ctx, *transcodeConfig)
			if err != nil {
				return fmt.Errorf("failed to start transcoding: %w", err)
			}
			return nil
		case string(constants.Deleted):
			datalivesession, err := c.livesession.GetLiveSessionByID(ctx, data.LiveSessionID)
			if err != nil {
				return fmt.Errorf("failed to get live session: %w", err)
			}
			if datalivesession == nil {
				return fmt.Errorf("live session not found")
			}
			if datalivesession.RecordingSetting.IsRecorded == false {
				return nil
			}
			datafile, err := c.seaweedfs.GenerateVOD(ctx, data.LiveSessionID, data.OwnerID, data.Name)
			if err != nil {
				return fmt.Errorf("failed to generate VOD: %w", err)
			}
			log.Printf("Generated VOD file: %s", datafile.PublicURL)

			datalivesession.PlaybackURL = datafile.PublicURL
			datalivesession.RecordingSetting.ArchiveURL = datafile.PublicURL
			err = c.livesession.UpdateLiveSession(ctx, datalivesession)
			if err != nil {
				return fmt.Errorf("failed to update live session with VOD URL: %w", err)
			}

		}
		return nil
	})
	if err != nil {
		// Log lỗi hoặc xử lý theo yêu cầu
		log.Printf("Error subscribing to events: %v", err)
	}
}
func (s *ConsumerLiveSession) watchAndUploadSegments(ctx context.Context, OwnerID, sessionID, outputDir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("[Worker] Error creating watcher for %s: %v", sessionID, err)
		return
	}
	defer watcher.Close()

	if err := watcher.Add(outputDir); err != nil {
		log.Printf("[Worker] Error watching directory %s: %v", outputDir, err)
		return
	}

	log.Printf("[Worker] Started watching: %s", outputDir)

	for {
		select {
		case <-ctx.Done(): // FFmpeg dừng -> Hủy worker theo dõi
			log.Printf("[Worker] Stopped watching for session: %s", sessionID)
			return

		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Lọc sự kiện Create hoặc Write
			if event.Has(fsnotify.Create) || event.Has(fsnotify.Write) {
				fileName := filepath.Base(event.Name)

				// 1. XỬ LÝ FILE MANIFEST (.m3u8)
				if strings.HasSuffix(fileName, ".m3u8") {
					// Best Practice: Đọc file m3u8 ĐỒNG BỘ ngay lập tức trên luồng chính để lấy "snapshot" mới nhất,
					// vì FFmpeg ghi đè file này rất nhanh. Nếu để Worker tự mở file sau, có thể bị lỗi Race Condition.
					content, err := os.ReadFile(event.Name)
					if err == nil && len(content) > 0 {

						// Đẩy tác vụ Upload sang WorkerPool
						err := s.pool.Run(ctx, func() {
							_ = s.seaweedfs.UpdateStreamManifest(ctx, OwnerID, sessionID, content)
						})
						if err != nil {
							log.Printf("[WorkerPool] Failed to assign m3u8 upload task: %v", err)
						}
					}
				}

				// 2. XỬ LÝ FILE SEGMENT (.ts)
				if strings.HasSuffix(fileName, ".ts") {
					// Best Practice: Copy biến ra scope cục bộ để tránh lỗi "Closure capture" trong Go
					// (Đặc biệt quan trọng khi truyền vào func() của WorkerPool)
					evtName := event.Name
					fName := fileName

					// Đẩy tác vụ mở file và upload sang WorkerPool để tránh block luồng theo dõi
					err := s.pool.Run(ctx, func() {
						fileInfo, err := os.Stat(evtName)
						// Chỉ upload khi file đã có dung lượng (FFmpeg đã bắt đầu ghi)
						if err == nil && fileInfo.Size() > 0 {
							file, err := os.Open(evtName)
							if err == nil {
								// BẮT BUỘC: Đóng file ngay trong Worker này khi xong tác vụ
								defer file.Close()

								_ = s.seaweedfs.UploadStreamSegment(ctx, OwnerID, sessionID, fName, file)
							}
						}
					})
					if err != nil {
						log.Printf("[WorkerPool] Failed to assign TS segment upload task: %v", err)
					}
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("[Worker] Watcher error for %s: %v", sessionID, err)
		}
	}
}
