package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

// --- SUB DTO ---
type LogDiffReq struct {
	OldContent    string   `json:"old_content"`
	NewContent    string   `json:"new_content"`
	ChangedFields []string `json:"changed_fields"`
}

// --- MAIN DTO ---
type PostEntityEditLogReq struct {
	ID               string                       `json:"id,omitempty"` // Trình bày dưới dạng string
	TargetCollection sharedEnums.TargetCollection `json:"target_collection" validate:"required"`
	TargetID         string                       `json:"target_id" validate:"required"` // String của ObjectID
	Version          int                          `json:"version" validate:"required"`
	EditedAt         time.Time                    `json:"edited_at" validate:"required"`
	EditorID         string                       `json:"editor_id" validate:"required"`
	Diff             LogDiffReq                   `json:"diff" validate:"required"` // Value struct vì bắt buộc có
	IPAddress        string                       `json:"ip_address,omitempty"`
	UserAgent        string                       `json:"user_agent,omitempty"`
	CreatedAt        time.Time                    `json:"created_at"`
	UpdatedAt        time.Time                    `json:"updated_at"`
	DeletedAt        *time.Time                   `json:"deleted_at,omitempty"`
}
