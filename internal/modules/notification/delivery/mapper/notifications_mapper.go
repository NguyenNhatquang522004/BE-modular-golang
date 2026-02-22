package mapper

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/entity"
	"github.com/gocql/gocql"
)

type ReqMapper struct{}

func NewReqMapper() *ReqMapper {
	return &ReqMapper{}
}

// ToEntity chuyển từ DTO Request sang Entity (Tạo mới)
func (m *ReqMapper) ToEntity(request *req.NotificationReq) (*entity.Notification, error) {
	userID, err := gocql.ParseUUID(request.UserID)
	if err != nil {
		return nil, err
	}

	actorID, err := gocql.ParseUUID(request.ActorID)
	if err != nil {
		return nil, err
	}

	// Nếu request không truyền NotificationID, ta tự gen TimeUUID theo chuẩn Cassandra
	notificationID := gocql.TimeUUID()
	if request.NotificationID != "" {
		notificationID, err = gocql.ParseUUID(request.NotificationID)
		if err != nil {
			return nil, err
		}
	}

	return &entity.Notification{
		UserID:         userID,
		CreatedAt:      request.CreatedAt,
		NotificationID: notificationID,
		Type:           request.Type,
		ActorID:        actorID,
		ActorName:      request.ActorName,
		ActorAvatar:    request.ActorAvatar,
		TargetID:       request.TargetID,
		TargetPreview:  request.TargetPreview,
		IsRead:         request.IsRead,
		IsClicked:      request.IsClicked,
		GroupKey:       request.GroupKey,
	}, nil
}

// UpdateToEntity cập nhật dữ liệu từ DTO Request ghi đè lên Entity có sẵn
func (m *ReqMapper) UpdateToEntity(request *req.NotificationReq, e *entity.Notification) error {
	if request.UserID != "" {
		userID, err := gocql.ParseUUID(request.UserID)
		if err != nil {
			return err
		}
		e.UserID = userID
	}

	if request.ActorID != "" {
		actorID, err := gocql.ParseUUID(request.ActorID)
		if err != nil {
			return err
		}
		e.ActorID = actorID
	}

	if request.NotificationID != "" {
		notifID, err := gocql.ParseUUID(request.NotificationID)
		if err != nil {
			return err
		}
		e.NotificationID = notifID
	}

	if !request.CreatedAt.IsZero() {
		e.CreatedAt = request.CreatedAt
	}

	// Cập nhật các trường còn lại
	e.Type = request.Type
	e.ActorName = request.ActorName
	e.ActorAvatar = request.ActorAvatar
	e.TargetID = request.TargetID
	e.TargetPreview = request.TargetPreview
	e.IsRead = request.IsRead
	e.IsClicked = request.IsClicked
	e.GroupKey = request.GroupKey

	return nil
}
