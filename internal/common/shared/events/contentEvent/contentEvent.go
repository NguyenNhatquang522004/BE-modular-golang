package contentEvent

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"

type SharePostPayload struct {
	UserID  string  `json:"user_id"`
	PostID  string  `json:"post_id"`
	GroupID *string `json:"group_id,omitempty"` // Nếu share vào group, có thể có group_id
	PageID  *string `json:"page_id,omitempty"`  // Nếu share vào page, có thể có page_id

}
type PostStatsPayload struct {
	UserID         string              `json:"user_id"`
	PostID         string              `json:"post_id"`
	TotalReactions int                 `json:"total_reactions"`
	Comments       int                 `json:"comments"`
	Shares         int                 `json:"shares"`
	Views          int                 `json:"views"`
	Like           int                 `json:"like"`
	Love           int                 `json:"love"`
	Haha           int                 `json:"haha"`
	Wow            int                 `json:"wow"`
	Sad            int                 `json:"sad"`
	Angry          int                 `json:"angry"`
	Type           constants.EventType `json:"type"`
}

type DeleteContentRelationTargetPayload struct {
	TargetID string `json:"target_id"`
}
