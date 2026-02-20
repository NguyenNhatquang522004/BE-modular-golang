package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/enum"
)

type CommentEditLogReq struct {
	ID               string                 `json:"id"` // Có thể rỗng khi Create, dùng cho Update
	TargetCollection *enum.TargetCollection `json:"target_collection" binding:"required"`
	TargetID         string                 `json:"target_id" binding:"required"` // String format của ObjectID
	Version          int                    `json:"version"`
	EditedAt         time.Time              `json:"edited_at"`
	EditorID         string                 `json:"editor_id" binding:"required,uuid"` // Đảm bảo là UUID từ Postgres
	Diff             *LogDiffReq             `json:"diff" binding:"required"`
	IPAddress        string                 `json:"ip_address"`
	UserAgent        string                 `json:"user_agent"`
}

type LogDiffReq struct {
	OldContent string            `json:"old_content"`
	NewContent string            `json:"new_content"`
	OldMedia   *MediaSnapshotReq `json:"old_media,omitempty"`
	NewMedia   *MediaSnapshotReq `json:"new_media,omitempty"`
}
type MediaSnapshotReq struct {
	Type string `json:"type" binding:"required,oneof=image video gif"`
	URL  string `json:"url" binding:"required,url"`
}
