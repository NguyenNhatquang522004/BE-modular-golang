package req

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"

type CommentPostAndReplyRequest struct {
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

type GetCommentPostRequest struct {
	PostID string `json:"post_id" validate:"required"`
	*dto.PaginationReq
}
type GetCommentDetailRequest struct {
	CommentID string `json:"comment_id" validate:"required"`
}
