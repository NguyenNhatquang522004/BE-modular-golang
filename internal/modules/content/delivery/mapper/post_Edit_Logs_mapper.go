package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- 1. REQ TO ENTITY (CREATE) ---

func ToPostEntityEditLogEntity(r req.CreatePostEntityEditLogReq) (*entity.PostEntityEditLog, error) {
	// Convert TargetID string -> ObjectID
	targetObjID, err := primitive.ObjectIDFromHex(r.TargetID)
	if err != nil {
		return nil, err
	}

	ent := &entity.PostEntityEditLog{
		ID:               primitive.NewObjectID(),
		TargetCollection: r.TargetCollection,
		TargetID:         targetObjID,
		Version:          r.Version,
		EditedAt:         time.Now(), // Thời điểm tạo log là thời điểm hiện tại
		EditorID:         r.EditorID,
		Diff: entity.LogDiff{
			OldContent:    r.Diff.OldContent,
			NewContent:    r.Diff.NewContent,
			ChangedFields: r.Diff.ChangedFields,
		},
		IPAddress: r.IPAddress,
		UserAgent: r.UserAgent,
	}

	return ent, nil
}

// --- 2. REQ TO ENTITY (UPDATE) ---

func UpdatePostEntityEditLogEntity(existingEnt *entity.PostEntityEditLog, r req.UpdatePostEntityEditLogReq) {
	// Chỉ update các trường được gửi lên (khác nil hoặc khác rỗng)

	if r.Version != nil {
		existingEnt.Version = *r.Version
	}

	if r.Diff != nil {
		existingEnt.Diff = entity.LogDiff{
			OldContent:    r.Diff.OldContent,
			NewContent:    r.Diff.NewContent,
			ChangedFields: r.Diff.ChangedFields,
		}
	}

	if r.IPAddress != "" {
		existingEnt.IPAddress = r.IPAddress
	}

	if r.UserAgent != "" {
		existingEnt.UserAgent = r.UserAgent
	}

	// Lưu ý: TargetCollection, TargetID, EditorID, EditedAt thường không cho phép sửa
	// trong logic Audit Log chuẩn, nên mình không đưa vào hàm update này.
}

// --- 3. ENTITY TO RES ---

func ToPostEntityEditLogRes(ent *entity.PostEntityEditLog) *res.PostEntityEditLogRes {
	if ent == nil {
		return nil
	}

	return &res.PostEntityEditLogRes{
		ID:               ent.ID.Hex(),
		TargetCollection: ent.TargetCollection,
		TargetID:         ent.TargetID.Hex(),
		Version:          ent.Version,
		EditedAt:         ent.EditedAt,
		EditorID:         ent.EditorID,
		Diff: res.LogDiffRes{
			OldContent:    ent.Diff.OldContent,
			NewContent:    ent.Diff.NewContent,
			ChangedFields: ent.Diff.ChangedFields,
		},
		IPAddress: ent.IPAddress,
		UserAgent: ent.UserAgent,
	}
}
