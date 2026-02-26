package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

// --- SUB DTO ---
type LogDiffRes struct {
	OldContent    string   `json:"old_content"`
	NewContent    string   `json:"new_content"`
	ChangedFields []string `json:"changed_fields"`
}

// --- MAIN DTO ---
type PostEntityEditLogRes struct {
	ID               string                       `json:"id"`
	TargetCollection sharedEnums.TargetCollection `json:"target_collection"`
	TargetID         string                       `json:"target_id"`
	Version          int                          `json:"version"`
	EditedAt         time.Time                    `json:"edited_at"`
	EditorID         string                       `json:"editor_id"`
	Diff             LogDiffRes                   `json:"diff"`
	IPAddress        string                       `json:"ip_address"`
	UserAgent        string                       `json:"user_agent"`
	CreatedAt        time.Time                    `json:"created_at"`
	UpdatedAt        time.Time                    `json:"updated_at"`
	DeletedAt        *time.Time                   `json:"deleted_at,omitempty"`
}
