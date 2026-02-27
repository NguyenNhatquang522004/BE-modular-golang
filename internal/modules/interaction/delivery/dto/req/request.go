package req

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
)

type CommentPostAndReplyRequest struct {
	Duration  float64 `json:"duration,omitempty"`
	SizeBytes int64   `json:"size_bytes,omitempty"`
	MimeType  string  `json:"mime_type,omitempty"`
	*CreateCommentReq
}
type CommentEditUseRequest struct {
	*CommentEditLogReq
}
type CommentReactionUseRequest struct {
	*EntityReactionReq
}
type ReactionPostUseRequest struct {
	*EntityReactionReq
}

type DeleteCommentRequest struct {
	CommentID string `json:"comment_id" validate:"required"`
}
type BookmarkPostRequest struct {
	*UserSavedItemReq
}

type DeleteBookmarkRequest struct {
	SaveItemID string `json:"save_item_id" validate:"required"`
}
type CommentCountRequest struct {
	CommentID    string              `json:"comment_id,omitempty"` // Nếu có comment_id thì đếm reply của comment đó, không có thì đếm comment của post
	ReplyCount   int                 `json:"reply_count,omitempty"`
	MentionCount int                 `json:"mention_count,omitempty"`
	ReportCount  int                 `json:"report_count,omitempty"`
	EventType    constants.TopicName `json:"event_type,omitempty"`
}
type GetCommentPostRequest struct {
	TargetID string `json:"target_id" validate:"required"`
	*dto.PaginationReq
}
type GetCommentDetailRequest struct {
	CommentID string `json:"comment_id" validate:"required"`
}

type GetBookMarkPostRequest struct {
	UserID string `json:"user_id" validate:"required"`
	*dto.PaginationReq
}
