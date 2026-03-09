package mapper

import (
	"errors"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/interactionEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// =============================================================================
// REQUEST MAPPER
// =============================================================================

// ToEntity: Chuyển từ Request DTO sang Entity (Dùng cho Create)
func ToEntityEditLogsPayload(req *interactionEvent.EditTargetPayload) (*entity.EntityEditLog, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	// Parse TargetID từ string sang ObjectID
	targetID, err := primitive.ObjectIDFromHex(req.TargetID)
	if err != nil {
		return nil, errors.New("invalid target_id format")
	}

	ent := &entity.EntityEditLog{
		TargetCollection: req.TargetCollection,
		TargetID:         targetID,
		EditedAt:         req.EditedAt,
		EditorID:         req.EditorID,
		Diff:             mapDiffReqToEntityPayload(&req.Diff),
		IPAddress:        req.IPAddress,
		UserAgent:        req.UserAgent,
	}

	// Nếu ID có truyền lên (trường hợp custom ID hoặc copy), thì parse ID
	if req.ID != "" {
		id, err := primitive.ObjectIDFromHex(req.ID)
		if err == nil {
			ent.ID = id
		}
	}

	return ent, nil
}

// UpdateToEntity: Cập nhật các trường từ Request vào một Entity đang có sẵn (Dùng cho Update)
func UpdateToEntityEditLogsPayload(req *interactionEvent.EditTargetPayload, ent *entity.EntityEditLog) error {
	if req == nil || ent == nil {
		return errors.New("request and entity cannot be nil")
	}

	if req.TargetID != "" {
		targetID, err := primitive.ObjectIDFromHex(req.TargetID)
		if err != nil {
			return errors.New("invalid target_id format")
		}
		ent.TargetID = targetID
	}

	// Cập nhật các trường còn lại
	ent.TargetCollection = req.TargetCollection
	ent.EditedAt = req.EditedAt
	ent.EditorID = req.EditorID
	ent.Diff = mapDiffReqToEntityPayload(&req.Diff)
	ent.IPAddress = req.IPAddress
	ent.UserAgent = req.UserAgent

	return nil
}

// Helper function để map sub-struct Diff (Req -> Entity)
func mapDiffReqToEntityPayload(diffReq *interactionEvent.LogDiffPayload) entity.LogDiff {
	diff := entity.LogDiff{
		OldContent: diffReq.OldContent,
		NewContent: diffReq.NewContent,
	}

	if diffReq.OldMedia != nil {
		diff.OldMedia = &entity.MediaSnapshot{
			Type: diffReq.OldMedia.Type,
			URL:  diffReq.OldMedia.URL,
		}
	}

	if diffReq.NewMedia != nil {
		diff.NewMedia = &entity.MediaSnapshot{
			Type: diffReq.NewMedia.Type,
			URL:  diffReq.NewMedia.URL,
		}
	}

	return diff
}

// =============================================================================
// RESPONSE MAPPER
// =============================================================================

// ToRes: Chuyển từ Entity trả về Response DTO cho Client
func ToResEditLogsPayload(ent *entity.EntityEditLog) *res.CommentEditLogRes {
	if ent == nil {
		return nil
	}

	resa := &res.CommentEditLogRes{
		ID:               ent.ID.Hex(), // Chuyển ngược ObjectID thành chuỗi Hex
		TargetCollection: &ent.TargetCollection,
		TargetID:         ent.TargetID.Hex(),
		Version:          ent.Version,
		EditedAt:         ent.EditedAt,
		EditorID:         ent.EditorID,
		Diff:             mapDiffEntityToRes(ent.Diff),
		IPAddress:        ent.IPAddress,
		UserAgent:        ent.UserAgent,
	}

	return resa
}

// Helper function để map sub-struct Diff (Entity -> Res)
func mapDiffEntityToResPayload(diffEnt entity.LogDiff) res.LogDiffRes {
	diffRes := res.LogDiffRes{
		OldContent: diffEnt.OldContent,
		NewContent: diffEnt.NewContent,
	}

	if diffEnt.OldMedia != nil {
		diffRes.OldMedia = &res.MediaSnapshotRes{
			Type: diffEnt.OldMedia.Type,
			URL:  diffEnt.OldMedia.URL,
		}
	}

	if diffEnt.NewMedia != nil {
		diffRes.NewMedia = &res.MediaSnapshotRes{
			Type: diffEnt.NewMedia.Type,
			URL:  diffEnt.NewMedia.URL,
		}
	}

	return diffRes
}
