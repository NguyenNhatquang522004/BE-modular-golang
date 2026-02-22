package cassandra

import (
	"context"
	"fmt"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/cassandraErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/entity"
	"github.com/gocql/gocql"
)

type NotificationsRepository struct {
	session *gocql.Session
	redis   IRepositoryShare.IRedis
	pool    IRepositoryShare.IWorkerPool
}

func NewNotificationsRepository(session *gocql.Session, redis IRepositoryShare.IRedis, pool IRepositoryShare.IWorkerPool) *NotificationsRepository {
	return &NotificationsRepository{
		session: session,
		redis:   redis,
		pool:    pool,
	}
}
func (r *NotificationsRepository) CreateNotification(ctx context.Context, notification *entity.Notification) error {
	tableName := entity.Notification{}.TableName()
	query := fmt.Sprintf(`INSERT INTO %s (user_id, created_at, notification_id, type, actor_id, actor_name, actor_avatar, target_id, target_preview) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, tableName)
	return r.session.Query(query,
		notification.UserID,
		notification.CreatedAt,
		notification.NotificationID,
		notification.Type,
		notification.ActorID,
		notification.ActorName,
		notification.ActorAvatar,
		notification.TargetID,
		notification.TargetPreview,
	).WithContext(ctx).Exec()
}
func (r *NotificationsRepository) CreateBulkNotifications(ctx context.Context, notifications []*entity.Notification) ([]*entity.Notification, []*cassandraErrors.NotificationBulkError, error) {
	tableName := entity.Notification{}.TableName()
	if len(notifications) == 0 {
		return nil, nil, nil
	}
	type taskResult struct {
		notification *entity.Notification
		err          error
	}
	query := fmt.Sprintf(`INSERT INTO %s (user_id, created_at, notification_id, type, actor_id, actor_name, actor_avatar, target_id, target_preview) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, tableName)
	resultsCh := make(chan taskResult, len(notifications))
	for _, notification := range notifications {
		if notification.UserID == (gocql.UUID{}) || notification.CreatedAt.IsZero() || notification.NotificationID == (gocql.UUID{}) {
			resultsCh <- taskResult{notification: nil, err: fmt.Errorf("invalid notification data: missing required fields")}
			continue
		}
		notification := notification // capture range variable
		err := r.pool.Run(ctx, func() {
			safeCtx := context.WithoutCancel(ctx)
			err := r.session.Query(query,
				notification.UserID,
				notification.CreatedAt,
				notification.NotificationID,
				notification.Type,
				notification.ActorID,
				notification.ActorName,
				notification.ActorAvatar,
				notification.TargetID,
				notification.TargetPreview,
			).WithContext(safeCtx).Exec()
			resultsCh <- taskResult{notification: notification, err: err}
		})
		if err != nil {
			resultsCh <- taskResult{notification: nil, err: fmt.Errorf("failed to execute query for notification with UserID %s and NotificationID %s: %w", notification.UserID, notification.NotificationID, err)}
			continue
		}
	}
	var successfulNotifications []*entity.Notification
	var bulkErrors []*cassandraErrors.NotificationBulkError
	for result := range resultsCh {
		if result.err != nil {
			var userID string
			var createdAt time.Time
			var notificationID string
			if result.notification != nil {
				userID = result.notification.UserID.String()
				createdAt = result.notification.CreatedAt
				notificationID = result.notification.NotificationID.String()
			}
			bulkErrors = append(bulkErrors, &cassandraErrors.NotificationBulkError{
				UserID:         userID,
				CreatedAt:      createdAt,
				NotificationID: notificationID,
				Error:          result.err.Error(),
			})
		} else {
			successfulNotifications = append(successfulNotifications, result.notification)
		}
	}
	return successfulNotifications, bulkErrors, nil
}
func (r *NotificationsRepository) GetNotificationsByUserID(ctx context.Context, userID string, createdAt time.Time, notificationID string) (*entity.Notification, error) {
	tableName := entity.Notification{}.TableName()
	var notifications = &entity.Notification{}
	finalUserID, err := gocql.ParseUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid userID: %w", err)
	}
	finalNotificationID, err := gocql.ParseUUID(notificationID)
	if err != nil {
		return nil, fmt.Errorf("invalid notificationID: %w", err)
	}
	query := fmt.Sprintf(`SELECT user_id, created_at, notification_id, type, actor_id, actor_name, actor_avatar, target_id, target_preview FROM %s WHERE user_id = ? AND created_at = ? AND notification_id = ? `, tableName)
	err = r.session.Query(query, finalUserID, createdAt, finalNotificationID).WithContext(ctx).Scan(
		&notifications.UserID,
		&notifications.CreatedAt,
		&notifications.NotificationID,
		&notifications.Type,
		&notifications.ActorID,
		&notifications.ActorName,
		&notifications.ActorAvatar,
		&notifications.TargetID,
		&notifications.TargetPreview,
	)
	if err != nil {
		if err == gocql.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get notifications for user %s: %w", userID, err)
	}
	return notifications, nil
}
func (r *NotificationsRepository) UpdateNotification(ctx context.Context, notification *entity.Notification) error {
	tableName := entity.Notification{}.TableName()
	query := fmt.Sprintf(`UPDATE %s SET type = ?, actor_id = ?, actor_name = ?, actor_avatar = ?, target_id = ?, target_preview = ? WHERE user_id = ? AND created_at = ? AND notification_id = ?`, tableName)
	return r.session.Query(query,
		notification.Type,
		notification.ActorID,
		notification.ActorName,
		notification.ActorAvatar,
		notification.TargetID,
		notification.TargetPreview,
		notification.UserID,
		notification.CreatedAt,
		notification.NotificationID,
	).WithContext(ctx).Exec()
}
func (r *NotificationsRepository) UpdateBulkNotifications(ctx context.Context, notifications []*entity.Notification) ([]*entity.Notification, []*cassandraErrors.NotificationBulkError, error) {
	tableName := entity.Notification{}.TableName()
	if len(notifications) == 0 {
		return nil, nil, nil
	}
	type taskResult struct {
		notification *entity.Notification
		err          error
	}
	query := fmt.Sprintf(`UPDATE %s SET type = ?, actor_id = ?, actor_name = ?, actor_avatar = ?, target_id = ?, target_preview = ? WHERE user_id = ? AND created_at = ? AND notification_id = ?`, tableName)
	resultsCh := make(chan taskResult, len(notifications))
	for _, notification := range notifications {
		if notification.UserID == (gocql.UUID{}) || notification.CreatedAt.IsZero() || notification.NotificationID == (gocql.UUID{}) {
			resultsCh <- taskResult{notification: nil, err: fmt.Errorf("invalid notification data: missing required fields")}
			continue
		}
		notification := notification // capture range variable
		err := r.pool.Run(ctx, func() {
			safeCtx := context.WithoutCancel(ctx)
			err := r.session.Query(query,
				notification.Type,
				notification.ActorID,
				notification.ActorName,
				notification.ActorAvatar,
				notification.TargetID,
				notification.TargetPreview,
				notification.UserID,
				notification.CreatedAt,
				notification.NotificationID,
			).WithContext(safeCtx).Exec()
			resultsCh <- taskResult{notification: notification, err: err}
		})
		if err != nil {
			resultsCh <- taskResult{notification: nil, err: fmt.Errorf("failed to execute query for notification with UserID %s and NotificationID %s: %w", notification.UserID, notification.NotificationID, err)}
			continue
		}
	}
	var successfulNotifications []*entity.Notification
	var bulkErrors []*cassandraErrors.NotificationBulkError
	for result := range resultsCh {
		if result.err != nil {
			var userID string
			var createdAt time.Time
			var notificationID string
			if result.notification != nil {
				userID = result.notification.UserID.String()
				createdAt = result.notification.CreatedAt
				notificationID = result.notification.NotificationID.String()
			}
			bulkErrors = append(bulkErrors, &cassandraErrors.NotificationBulkError{
				UserID:         userID,
				CreatedAt:      createdAt,
				NotificationID: notificationID,
				Error:          result.err.Error(),
			})
		} else {
			successfulNotifications = append(successfulNotifications, result.notification)
		}
	}
	return successfulNotifications, bulkErrors, nil
}
func (r *NotificationsRepository) DeleteNotification(ctx context.Context, userID string, createdAt time.Time, notificationID string) error {
	tableName := entity.Notification{}.TableName()
	finalUserID, err := gocql.ParseUUID(userID)
	if err != nil {
		return fmt.Errorf("invalid userID: %w", err)
	}
	finalNotificationID, err := gocql.ParseUUID(notificationID)
	if err != nil {
		return fmt.Errorf("invalid notificationID: %w", err)
	}
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id = ? AND created_at = ? AND notification_id = ?`, tableName)
	return r.session.Query(query, finalUserID, createdAt, finalNotificationID).WithContext(ctx).Exec()
}
func (r *NotificationsRepository) DeleteBulkNotifications(ctx context.Context, notifications []*entity.Notification) ([]*entity.Notification, []*cassandraErrors.NotificationBulkError, error) {
	tableName := entity.Notification{}.TableName()
	if len(notifications) == 0 {
		return nil, nil, nil
	}
	type taskResult struct {
		notification *entity.Notification
		err          error
	}
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id = ? AND created_at = ? AND notification_id = ?`, tableName)
	resultsCh := make(chan taskResult, len(notifications))
	for _, notification := range notifications {
		if notification.UserID == (gocql.UUID{}) || notification.CreatedAt.IsZero() || notification.NotificationID == (gocql.UUID{}) {
			resultsCh <- taskResult{notification: nil, err: fmt.Errorf("invalid notification data: missing required fields")}
			continue
		}
		notification := notification // capture range variable
		err := r.pool.Run(ctx, func() {
			safeCtx := context.WithoutCancel(ctx)
			err := r.session.Query(query,
				notification.UserID,
				notification.CreatedAt,
				notification.NotificationID,
			).WithContext(safeCtx).Exec()
			resultsCh <- taskResult{notification: notification, err: err}
		})
		if err != nil {
			resultsCh <- taskResult{notification: nil, err: fmt.Errorf("failed to execute query for notification with UserID %s and NotificationID %s: %w", notification.UserID, notification.NotificationID, err)}
			continue
		}
	}
	var successfulNotifications []*entity.Notification
	var bulkErrors []*cassandraErrors.NotificationBulkError
	for result := range resultsCh {
		if result.err != nil {
			var userID string
			var createdAt time.Time
			var notificationID string
			if result.notification != nil {
				userID = result.notification.UserID.String()
				createdAt = result.notification.CreatedAt
				notificationID = result.notification.NotificationID.String()
			}
			bulkErrors = append(bulkErrors, &cassandraErrors.NotificationBulkError{
				UserID:         userID,
				CreatedAt:      createdAt,
				NotificationID: notificationID,
				Error:          result.err.Error(),
			})
		} else {
			successfulNotifications = append(successfulNotifications, result.notification)
		}
	}
	return successfulNotifications, bulkErrors, nil
}
