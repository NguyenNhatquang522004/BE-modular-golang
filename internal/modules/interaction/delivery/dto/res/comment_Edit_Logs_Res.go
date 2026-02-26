package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type CommentEditLogRes struct {
	ID               string                        `json:"id"`
	TargetCollection *sharedEnums.TargetCollection `json:"target_collection"`
	TargetID         string                        `json:"target_id"`
	Version          int                           `json:"version"`
	EditedAt         time.Time                     `json:"edited_at"`
	EditorID         string                        `json:"editor_id"`
	Diff             LogDiffRes                    `json:"diff"`
	IPAddress        string                        `json:"ip_address"`
	UserAgent        string                        `json:"user_agent"`
}

type LogDiffRes struct {
	OldContent string            `json:"old_content"`
	NewContent string            `json:"new_content"`
	OldMedia   *MediaSnapshotRes `json:"old_media,omitempty"`
	NewMedia   *MediaSnapshotRes `json:"new_media,omitempty"`
}

type MediaSnapshotRes struct {
	Type string `json:"type" binding:"required,oneof=image video gif"`
	URL  string `json:"url" binding:"required,url"`
}
