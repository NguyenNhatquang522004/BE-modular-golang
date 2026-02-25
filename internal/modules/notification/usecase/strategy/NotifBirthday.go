package strategy

import (
	"context"
	"errors"
	"log"
	"slices"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/notificationEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/socket"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/entity"
	"github.com/gocql/gocql"
)

type NotifBirthday struct {
	notificationRepo             IRepositoryCassandra.INotificationsRepository
	notificationTemplateRepo     IRepositoryMongodb.INotificationTemplatesRepository
	userNotificationSettingsRepo IRepositoryMongodb.IUserNotificationSettingsRepository
	pool                         IRepositoryShare.IWorkerPool
	socket                       socket.Manager
}

func NewNotifBirthday(notificationRepo IRepositoryCassandra.INotificationsRepository,
	notificationTemplateRepo IRepositoryMongodb.INotificationTemplatesRepository,
	userNotificationSettingsRepo IRepositoryMongodb.IUserNotificationSettingsRepository,
	pool IRepositoryShare.IWorkerPool, socketManager socket.Manager) *NotifBirthday {
	return &NotifBirthday{
		notificationRepo:             notificationRepo,
		notificationTemplateRepo:     notificationTemplateRepo,
		userNotificationSettingsRepo: userNotificationSettingsRepo,
		pool:                         pool,
		socket:                       socketManager,
	}
}
func (r *NotifBirthday) Execute(ctx context.Context, req *notificationEvent.NotificationPayload) error {
	// Implement the logic for executing the notification strategy
	thisEnum := r.GetType()
	isExist := slices.Contains(req.TypeNotification, &thisEnum)
	if !isExist {
		return errors.New(" notification type is not valid for this strategy")
	}
	dataActor, err := r.userNotificationSettingsRepo.GetUserNotificationSettingsByUserID(ctx, req.UserID)
	if err != nil {
		return err
	}
	dataTemplate, err := r.notificationTemplateRepo.GetOneTemplateByType(ctx, thisEnum)
	if err != nil {
		return err
	}

	for _, userID := range req.SendUser {
		err := r.pool.Run(ctx, func() {
			datauser, err := r.userNotificationSettingsRepo.GetUserNotificationSettingsByUserID(ctx, userID)
			if err != nil {
				log.Printf("Error fetching notification settings for user %s: %v", userID, err)
				// Handle the error appropriately, e.g., log it or return an error
				return
			}
			if datauser.Settings.PushEnabled == false || datauser.Settings.PushBirthdays == false {
				// User has muted all notifications or birthday notifications, skip sending
				return
			}
			safectx := context.WithoutCancel(ctx)
			UserIDfinal, ok := gocql.ParseUUID(userID)
			if ok != nil {
				// Handle the error appropriately, e.g., log it or return an error
				return
			}
			ActorID, ok := gocql.ParseUUID(req.BirthdayID)
			if ok != nil {
				// Handle the error appropriately, e.g., log it or return an error
				return
			}
			var entityNotification = &entity.Notification{
				UserID:         UserIDfinal,
				CreatedAt:      time.Now(),
				NotificationID: gocql.TimeUUID(),
				Type:           sharedEnums.NotifBirthday,
				ActorID:        ActorID,
				ActorName:      dataActor.Name,
				ActorAvatar:    dataActor.Avatar,
				TargetID:       req.BirthdayID,
				TargetPreview:  "",
				IsRead:         false,
				IsClicked:      false,
				GroupKey:       ""}
			err = r.notificationRepo.CreateNotification(safectx, entityNotification)
			if err != nil {
				log.Printf("Error creating notification for user %s: %v", userID, err)
				// Handle the error appropriately, e.g., log it or return an error
				return
			}
			transmit := map[string]string{
				"actorname":   dataActor.Name,
				"actoravatar": dataActor.Avatar,
			}
			dataTemplate.Template = transmit
			r.socket.SendToUser(userID, socket.Message{
				Payload: dataTemplate,
			})
		})
		if err != nil {
			log.Printf("Error running notification task for user %s: %v", userID, err)
			// Handle the error appropriately, e.g., log it or return an error
			return err
		}

	}
	r.pool.Wait() // Wait for all goroutines to finish
	return nil
}
func (r *NotifBirthday) GetType() sharedEnums.NotificationType {
	return sharedEnums.NotifBirthday
}
