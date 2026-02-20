package mapper

import (
	"errors"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToEntityUserSavedItem(reqPayload *req.UserSavedItemReq) (*entity.UserSavedItem, error) {
	if reqPayload == nil {
		return nil, errors.New("request cannot be nil")
	}

	targetID, err := primitive.ObjectIDFromHex(reqPayload.TargetID)
	if err != nil {
		return nil, errors.New("invalid target_id format")
	}

	// Xử lý an toàn Enum con trỏ
	var targetType enum.SavedTargetType
	if reqPayload.TargetType != nil {
		targetType = *reqPayload.TargetType
	}

	ent := &entity.UserSavedItem{
		UserID:         reqPayload.UserID,
		TargetID:       targetID,
		TargetType:     targetType,
		Snapshot:       mapSnapshotReqToEntity(reqPayload.Snapshot),
		CollectionName: reqPayload.CollectionName,
		CreatedAt:      reqPayload.CreatedAt,
	}

	// Nếu có ID truyền lên, parse nó
	if reqPayload.ID != "" {
		id, err := primitive.ObjectIDFromHex(reqPayload.ID)
		if err == nil {
			ent.ID = id
		}
	}

	return ent, nil
}

// UpdateToEntity: Cập nhật thông tin từ Request vào Entity có sẵn (Dùng cho Update)
func UpdateToEntityUserSavedItem(reqPayload *req.UserSavedItemReq, ent *entity.UserSavedItem) error {
	if reqPayload == nil || ent == nil {
		return errors.New("request and entity cannot be nil")
	}

	if reqPayload.TargetID != "" {
		targetID, err := primitive.ObjectIDFromHex(reqPayload.TargetID)
		if err != nil {
			return errors.New("invalid target_id format")
		}
		ent.TargetID = targetID
	}

	if reqPayload.TargetType != nil {
		ent.TargetType = *reqPayload.TargetType
	}

	ent.UserID = reqPayload.UserID
	ent.Snapshot = mapSnapshotReqToEntity(reqPayload.Snapshot)
	ent.CollectionName = reqPayload.CollectionName
	ent.CreatedAt = reqPayload.CreatedAt

	return nil
}

// Helper: Chuyển đổi Snapshot từ Request sang Entity
func mapSnapshotReqToEntity(snapReq *req.SavedItemSnapshotReq) entity.SavedItemSnapshot {
	if snapReq == nil {
		return entity.SavedItemSnapshot{}
	}

	return entity.SavedItemSnapshot{
		AuthorName:     snapReq.AuthorName,
		ContentPreview: snapReq.ContentPreview,
		ThumbnailURL:   snapReq.ThumbnailURL,
	}
}

// =============================================================================
// RESPONSE MAPPER
// =============================================================================

// ToRes: Chuyển từ Entity trả về Response DTO
func ToResUserSavedItem(ent *entity.UserSavedItem) *res.UserSavedItemRes {
	if ent == nil {
		return nil
	}

	// Lấy giá trị enum để lấy địa chỉ con trỏ trả về cho an toàn
	targetType := ent.TargetType

	return &res.UserSavedItemRes{
		ID:             ent.ID.Hex(),
		UserID:         ent.UserID,
		TargetID:       ent.TargetID.Hex(),
		TargetType:     &targetType,
		Snapshot:       mapSnapshotEntityToRes(ent.Snapshot),
		CollectionName: ent.CollectionName,
		CreatedAt:      ent.CreatedAt,
	}
}

// Helper: Chuyển đổi Snapshot từ Entity sang Response
func mapSnapshotEntityToRes(snapEnt entity.SavedItemSnapshot) *res.SavedItemSnapshotRes {
	return &res.SavedItemSnapshotRes{
		AuthorName:     snapEnt.AuthorName,
		ContentPreview: snapEnt.ContentPreview,
		ThumbnailURL:   snapEnt.ThumbnailURL,
	}
}
