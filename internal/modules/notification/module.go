package notification

import (
	"context"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/entity"
	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ModuleNotification struct {
}

func (m *ModuleNotification) InitMongo(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Collection: notification_templates
	templateCol := db.Collection(entity.CollectionNotificationTemplates)
	_, err := templateCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		// Đảm bảo mỗi loại thông báo (VD: POST_LIKE) chỉ có duy nhất 1 template
		Keys:    bson.D{{Key: "type", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	// 2. Collection: user_notification_settings
	settingsCol := db.Collection(entity.CollectionUserNotificationSettings)
	_, err = settingsCol.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			// Mỗi user chỉ có một bản ghi cấu hình thông báo
			Keys:    bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			// Multikey Index: Tối ưu cho việc tìm kiếm/xóa một token FCM cụ thể trong mảng FCMTokens
			Keys: bson.D{{Key: "fcm_tokens.token", Value: 1}},
		},
	})

	return err
}
func (m *ModuleNotification) InitCassandra(session *gocql.Session) error {
	// Khởi tạo bảng notifications (Cassandra)
	// PRIMARY KEY ((user_id), created_at, notification_id)
	// - Partition Key: user_id (Gom dữ liệu theo người dùng)
	// - Clustering Key 1: created_at (Sắp xếp theo thời gian)
	// - Clustering Key 2: notification_id (Đảm bảo duy nhất nếu tạo noti cùng lúc)
	query := `
	CREATE TABLE IF NOT EXISTS notifications (
		user_id uuid,
		created_at timestamp,
		notification_id uuid,
		type text,
		actor_id uuid,
		actor_name text,
		actor_avatar text,
		target_id text,
		target_preview text,
		is_read boolean,
		is_clicked boolean,
		group_key text,
		PRIMARY KEY ((user_id), created_at, notification_id)
	) WITH CLUSTERING ORDER BY (created_at DESC, notification_id DESC);`

	return session.Query(query).Exec()
}
