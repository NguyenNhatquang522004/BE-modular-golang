package mapper

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- 1. Request -> Entity (Tạo mới) ---
func ToPostEntityEditLogEntity(r *req.PostEntityEditLogReq) *entity.PostEntityEditLog {
	if r == nil {
		return nil
	}

	// Xử lý ID chính
	objectID := primitive.NewObjectID()
	if r.ID != "" {
		if oid, err := primitive.ObjectIDFromHex(r.ID); err == nil {
			objectID = oid
		}
	}

	// Xử lý TargetID
	targetID, _ := primitive.ObjectIDFromHex(r.TargetID)

	logEntity := &entity.PostEntityEditLog{
		ID:               objectID,
		TargetCollection: r.TargetCollection,
		TargetID:         targetID,
		Version:          r.Version,
		EditedAt:         r.EditedAt,
		EditorID:         r.EditorID,
		Diff: entity.LogDiff{
			OldContent:    r.Diff.OldContent,
			NewContent:    r.Diff.NewContent,
			ChangedFields: r.Diff.ChangedFields,
		},
		IPAddress: r.IPAddress,
		UserAgent: r.UserAgent,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		DeletedAt: r.DeletedAt,
	}

	return logEntity
}

// --- 2. Update Request -> Existing Entity ---
func UpdatePostEntityEditLogEntity(r *req.PostEntityEditLogReq, logEntity *entity.PostEntityEditLog) {
	if r == nil || logEntity == nil {
		return
	}

	// Lưu ý: Đối với bảng Log, thường ít khi Update.
	// Dưới đây chỉ map những trường có thể thay đổi hợp lý, bảo vệ ID, TargetID, CreatedAt.
	logEntity.TargetCollection = r.TargetCollection
	logEntity.Version = r.Version
	logEntity.EditedAt = r.EditedAt
	logEntity.EditorID = r.EditorID
	logEntity.IPAddress = r.IPAddress
	logEntity.UserAgent = r.UserAgent
	logEntity.UpdatedAt = r.UpdatedAt
	logEntity.DeletedAt = r.DeletedAt

	logEntity.Diff = entity.LogDiff{
		OldContent:    r.Diff.OldContent,
		NewContent:    r.Diff.NewContent,
		ChangedFields: r.Diff.ChangedFields,
	}
}

// --- 3. Entity -> Response ---
func ToPostEntityEditLogRes(logEntity *entity.PostEntityEditLog) *res.PostEntityEditLogRes {
	if logEntity == nil {
		return nil
	}

	return &res.PostEntityEditLogRes{
		ID:               logEntity.ID.Hex(),
		TargetCollection: logEntity.TargetCollection,
		TargetID:         logEntity.TargetID.Hex(),
		Version:          logEntity.Version,
		EditedAt:         logEntity.EditedAt,
		EditorID:         logEntity.EditorID,
		Diff: res.LogDiffRes{
			OldContent:    logEntity.Diff.OldContent,
			NewContent:    logEntity.Diff.NewContent,
			ChangedFields: logEntity.Diff.ChangedFields,
		},
		IPAddress: logEntity.IPAddress,
		UserAgent: logEntity.UserAgent,
		CreatedAt: logEntity.CreatedAt,
		UpdatedAt: logEntity.UpdatedAt,
		DeletedAt: logEntity.DeletedAt,
	}
}

// --- 4. Request -> Response ---
func ReqToPostEntityEditLogRes(r *req.PostEntityEditLogReq) *res.PostEntityEditLogRes {
	if r == nil {
		return nil
	}

	return &res.PostEntityEditLogRes{
		ID:               r.ID,
		TargetCollection: r.TargetCollection,
		TargetID:         r.TargetID,
		Version:          r.Version,
		EditedAt:         r.EditedAt,
		EditorID:         r.EditorID,
		Diff: res.LogDiffRes{
			OldContent:    r.Diff.OldContent,
			NewContent:    r.Diff.NewContent,
			ChangedFields: r.Diff.ChangedFields,
		},
		IPAddress: r.IPAddress,
		UserAgent: r.UserAgent,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		DeletedAt: r.DeletedAt,
	}
}
