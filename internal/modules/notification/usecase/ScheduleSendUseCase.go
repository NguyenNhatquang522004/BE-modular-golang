package usecase

import (
	"context"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/notificationEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/entity"
	"github.com/google/uuid"
)

type ScheduleSendUseCase struct {
	cron                 IRepositoryShare.IScheduler
	userNotificationRepo IRepositoryMongodb.IUserNotificationSettingsRepository
	eventbus             events.EventBus
	pool                 IRepositoryShare.IWorkerPool
}

func NewScheduleSendUseCase(cron IRepositoryShare.IScheduler, userNotificationSettingsRepo IRepositoryMongodb.IUserNotificationSettingsRepository, eventbus events.EventBus, pool IRepositoryShare.IWorkerPool) *ScheduleSendUseCase {
	return &ScheduleSendUseCase{
		cron:                 cron,
		userNotificationRepo: userNotificationSettingsRepo,
		eventbus:             eventbus,
		pool:                 pool,
	}
}

func (s *ScheduleSendUseCase) Execute(ctx context.Context) error {
	// Đăng ký công việc gửi thông báo vào lúc 8 giờ sáng hàng ngày ("0 8 * * *")
	err := s.cron.ScheduleJob("0 8 * * *", func() {

		// 1. LỖI LOGIC ĐÃ SỬA: Phải lấy thời gian và tạo ID ở TRONG hàm này
		// để mỗi lần cron job tự động kích hoạt, nó sẽ lấy đúng ngày/tháng của hôm đó.
		today := time.Now()
		month := int(today.Month())
		day := today.Day()
		id := uuid.New()

		// 2. Tốt nhất nên dùng một context mới cho background job vì ctx của lúc start
		// server có thể sẽ bị hủy (cancel) sau khi hàm Execute chạy xong.
		jobCtx := context.Background()

		var cursor string
		var limit = 100
		var hasNext = true

		for hasNext {
			settings, err := s.userNotificationRepo.GetUserNotificationSettingsByDateOfBirth(jobCtx, month, day, cursor, limit)
			if err != nil {
				// Nên dùng thư viện log để ghi nhận lỗi ở đây rồi break
				break
			}

			hasNext = settings.HasNext
			cursor = settings.NextCursor

			// Ép kiểu an toàn (nên check nil trước khi parse nếu cần)
			data := settings.Data.(*[]entity.UserNotificationSetting)

			// Nếu mảng rỗng thì thoát loop luôn
			if data == nil || len(*data) == 0 {
				break
			}
			for _, setting := range *data {
				err := s.pool.Run(jobCtx, func() {
					// Tạo payload thông báo sinh nhật
					notifType := sharedEnums.NotifBirthday
					arr := []*sharedEnums.NotificationType{&notifType}
					senderuser := []string{setting.UserID}
					payload := notificationEvent.NotificationPayload{
						BirthdayID:       id.String(),
						TypeNotification: arr,
						SendUser:         senderuser,
					}

					// Đăng sự kiện gửi thông báo sinh nhật
					errBus := s.eventbus.Publish(jobCtx, string(constants.TopicSendNotificationType), id.String(), string(constants.None), payload)
					if errBus != nil {
						// Log lỗi publish event
						return
					}
				})

				if err != nil {
					// Log lỗi khi thêm task vào worker pool
					continue // Dùng continue để bỏ qua người dùng bị lỗi, tiếp tục gửi cho người khác
				}
			}
		} // <-- ĐÃ SỬA LỖI SYNTAX: Thêm dấu đóng ngoặc cho vòng lặp for hasNext
	})

	if err != nil {
		return err
	}

	s.cron.Start() // Bắt đầu bộ lập lịch
	return nil
}
