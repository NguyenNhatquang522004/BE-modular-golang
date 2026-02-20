package mapper

import (
	"errors"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrInvalidPostID = errors.New("invalid post_id format")
	ErrInvalidID     = errors.New("invalid id format")
)

// ==========================
// REQUEST MAPPERS
// ==========================

func ToEntityComment(req *req.CreateCommentReq) (*entity.Comment, error) {
	postID, err := primitive.ObjectIDFromHex(req.PostID)
	if err != nil {
		return nil, ErrInvalidPostID
	}

	now := time.Now()
	comment := &entity.Comment{
		ID:           primitive.NewObjectID(),
		PostID:       postID,
		UserID:       req.UserID,
		Content:      req.Content,
		Mentions:     req.Mentions,
		Status:       *req.Status,
		ReportCount:  req.ReportCount,
		ReplyCount:   req.ReplyCount,
		MentionCount: req.MentionCount,
		IsEdited:     req.IsEdited,
		LastEditedAt: req.LastEditedAt,
		DeletedAt:    req.DeletedAt,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if req.CreatedAt != nil {
		comment.CreatedAt = *req.CreatedAt
	}
	if req.UpdatedAt != nil {
		comment.UpdatedAt = *req.UpdatedAt
	}

	if req.Status != nil {
		comment.Status = *req.Status
	}

	if comment.MentionCount == 0 && len(req.Mentions) > 0 {
		comment.MentionCount = len(req.Mentions)
	}

	// Map CommentReactionsReq -> Entity
	if req.Reactions != nil {
		comment.Reactions = entity.CommentReactions{
			Total: req.Reactions.Total,
			Like:  req.Reactions.Like,
			Love:  req.Reactions.Love,
			Haha:  req.Reactions.Haha,
			Wow:   req.Reactions.Wow,
			Sad:   req.Reactions.Sad,
			Angry: req.Reactions.Angry,
		}
	}

	if req.AssetID != nil && *req.AssetID != "" {
		assetID, err := primitive.ObjectIDFromHex(*req.AssetID)
		if err == nil {
			comment.AssetID = &assetID
		}
	}

	if req.ParentCommentID != nil && *req.ParentCommentID != "" {
		parentID, err := primitive.ObjectIDFromHex(*req.ParentCommentID)
		if err == nil {
			comment.ParentCommentID = &parentID
		}
	}

	if req.RootCommentID != nil && *req.RootCommentID != "" {
		rootID, err := primitive.ObjectIDFromHex(*req.RootCommentID)
		if err == nil {
			comment.RootCommentID = &rootID
		}
	}

	if req.Media != nil {
		comment.Media = &entity.CommentMedia{
			Type: *req.Media.Type,
			URL:  req.Media.URL,
			DisplayMeta: entity.DisplayMeta{
				Width:  req.Media.DisplayMeta.Width,
				Height: req.Media.DisplayMeta.Height,
			},
		}
	}

	if req.HiddenMetadata != nil {
		comment.HiddenMetadata = &entity.HiddenMetadata{
			IsHidden:       req.HiddenMetadata.IsHidden,
			HiddenAt:       req.HiddenMetadata.HiddenAt,
			HiddenByUserID: req.HiddenMetadata.HiddenByUserID,
			Reason:         req.HiddenMetadata.Reason,
			IsGhostBanned:  req.HiddenMetadata.IsGhostBanned,
		}
	}

	if req.DeletedMetadata != nil {
		comment.DeletedMetadata = &entity.DeletedMetadata{
			DeletedAt:       req.DeletedMetadata.DeletedAt,
			DeletedByUserID: req.DeletedMetadata.DeletedByUserID,
		}
	}

	return comment, nil
}

func UpdateToEntityComment(req *req.UpdateCommentReq, comment *entity.Comment) {
	isModified := false
	now := time.Now()

	if req.Content != nil && *req.Content != comment.Content {
		comment.Content = *req.Content
		isModified = true
	}

	if req.Mentions != nil {
		comment.Mentions = *req.Mentions
		comment.MentionCount = len(*req.Mentions)
		isModified = true
	}

	if req.Status != nil && *req.Status != comment.Status {
		comment.Status = *req.Status
		isModified = true
	}

	if req.ReportCount != nil && *req.ReportCount != comment.ReportCount {
		comment.ReportCount = *req.ReportCount
		isModified = true
	}

	if req.ReplyCount != nil && *req.ReplyCount != comment.ReplyCount {
		comment.ReplyCount = *req.ReplyCount
		isModified = true
	}

	if req.MentionCount != nil && *req.MentionCount != comment.MentionCount {
		comment.MentionCount = *req.MentionCount
		isModified = true
	}

	if req.IsEdited != nil && *req.IsEdited != comment.IsEdited {
		comment.IsEdited = *req.IsEdited
		isModified = true
	}

	if req.LastEditedAt != nil {
		comment.LastEditedAt = req.LastEditedAt
		isModified = true
	}

	if req.DeletedAt != nil {
		comment.DeletedAt = req.DeletedAt
		isModified = true
	}

	// Cập nhật Reactions nếu có truyền lên
	if req.Reactions != nil {
		comment.Reactions = entity.CommentReactions{
			Total: req.Reactions.Total,
			Like:  req.Reactions.Like,
			Love:  req.Reactions.Love,
			Haha:  req.Reactions.Haha,
			Wow:   req.Reactions.Wow,
			Sad:   req.Reactions.Sad,
			Angry: req.Reactions.Angry,
		}
		isModified = true
	}

	if req.Media != nil {
		comment.Media = &entity.CommentMedia{
			Type: *req.Media.Type,
			URL:  req.Media.URL,
			DisplayMeta: entity.DisplayMeta{
				Width:  req.Media.DisplayMeta.Width,
				Height: req.Media.DisplayMeta.Height,
			},
		}
		isModified = true
	}

	if req.HiddenMetadata != nil {
		comment.HiddenMetadata = &entity.HiddenMetadata{
			IsHidden:       req.HiddenMetadata.IsHidden,
			HiddenAt:       req.HiddenMetadata.HiddenAt,
			HiddenByUserID: req.HiddenMetadata.HiddenByUserID,
			Reason:         req.HiddenMetadata.Reason,
			IsGhostBanned:  req.HiddenMetadata.IsGhostBanned,
		}
		isModified = true
	}

	if req.DeletedMetadata != nil {
		comment.DeletedMetadata = &entity.DeletedMetadata{
			DeletedAt:       req.DeletedMetadata.DeletedAt,
			DeletedByUserID: req.DeletedMetadata.DeletedByUserID,
		}
		isModified = true
	}

	if isModified {
		if req.Content != nil || req.Media != nil {
			comment.IsEdited = true
			comment.LastEditedAt = &now
		}
		comment.UpdatedAt = now
	}
}

// ==========================
// RESPONSE MAPPERS
// ==========================

func ToResponse(ent *entity.Comment) *res.CommentRes {
	if ent == nil {
		return nil
	}

	resa := &res.CommentRes{
		ID:           ent.ID.Hex(),
		PostID:       ent.PostID.Hex(),
		UserID:       ent.UserID,
		Content:      ent.Content,
		Mentions:     ent.Mentions,
		Status:       &ent.Status,
		ReportCount:  ent.ReportCount,
		ReplyCount:   ent.ReplyCount,
		MentionCount: ent.MentionCount,
		IsEdited:     ent.IsEdited,
		LastEditedAt: ent.LastEditedAt,
		CreatedAt:    ent.CreatedAt,
		UpdatedAt:    ent.UpdatedAt,
		DeletedAt:    ent.DeletedAt,
		// Map Entity -> CommentReactionsRes
		Reactions: &res.CommentReactionsRes{
			Total: ent.Reactions.Total,
			Like:  ent.Reactions.Like,
			Love:  ent.Reactions.Love,
			Haha:  ent.Reactions.Haha,
			Wow:   ent.Reactions.Wow,
			Sad:   ent.Reactions.Sad,
			Angry: ent.Reactions.Angry,
		},
	}

	if ent.AssetID != nil {
		assetStr := ent.AssetID.Hex()
		resa.AssetID = &assetStr
	}
	if ent.ParentCommentID != nil {
		parentStr := ent.ParentCommentID.Hex()
		resa.ParentCommentID = &parentStr
	}
	if ent.RootCommentID != nil {
		rootStr := ent.RootCommentID.Hex()
		resa.RootCommentID = &rootStr
	}

	if ent.Media != nil {
		resa.Media = &res.CommentMediaRes{
			Type: ent.Media.Type,
			URL:  ent.Media.URL,
			DisplayMeta: res.DisplayMetaRes{
				Width:  ent.Media.DisplayMeta.Width,
				Height: ent.Media.DisplayMeta.Height,
			},
		}
	}

	if ent.HiddenMetadata != nil {
		resa.HiddenMetadata = &res.HiddenMetadataRes{
			IsHidden:       ent.HiddenMetadata.IsHidden,
			HiddenAt:       ent.HiddenMetadata.HiddenAt,
			HiddenByUserID: ent.HiddenMetadata.HiddenByUserID,
			Reason:         ent.HiddenMetadata.Reason,
			IsGhostBanned:  ent.HiddenMetadata.IsGhostBanned,
		}
	}

	if ent.DeletedMetadata != nil {
		resa.DeletedMetadata = &res.DeletedMetadataRes{
			DeletedAt:       ent.DeletedMetadata.DeletedAt,
			DeletedByUserID: ent.DeletedMetadata.DeletedByUserID,
		}
	}

	return resa
}

func ReqToResponse(req *req.CreateCommentReq, genID string) *res.CommentRes {
	now := time.Now()
	resa := &res.CommentRes{
		ID:              genID,
		PostID:          req.PostID,
		UserID:          req.UserID,
		AssetID:         req.AssetID,
		Content:         req.Content,
		Mentions:        req.Mentions,
		ParentCommentID: req.ParentCommentID,
		RootCommentID:   req.RootCommentID,
		Status:          req.Status,
		ReportCount:     req.ReportCount,
		ReplyCount:      req.ReplyCount,
		MentionCount:    req.MentionCount,
		IsEdited:        req.IsEdited,
		LastEditedAt:    req.LastEditedAt,
		CreatedAt:       now,
		UpdatedAt:       now,
		DeletedAt:       req.DeletedAt,
	}

	if req.CreatedAt != nil {
		resa.CreatedAt = *req.CreatedAt
	}
	if req.UpdatedAt != nil {
		resa.UpdatedAt = *req.UpdatedAt
	}

	if req.Status != nil {
		resa.Status = req.Status
	}

	if resa.MentionCount == 0 && len(req.Mentions) > 0 {
		resa.MentionCount = len(req.Mentions)
	}

	// Map CommentReactionsReq -> CommentReactionsRes
	if req.Reactions != nil {
		resa.Reactions = &res.CommentReactionsRes{
			Total: req.Reactions.Total,
			Like:  req.Reactions.Like,
			Love:  req.Reactions.Love,
			Haha:  req.Reactions.Haha,
			Wow:   req.Reactions.Wow,
			Sad:   req.Reactions.Sad,
			Angry: req.Reactions.Angry,
		}
	}

	if req.Media != nil {
		resa.Media = &res.CommentMediaRes{
			Type: *req.Media.Type,
			URL:  req.Media.URL,
			DisplayMeta: res.DisplayMetaRes{
				Width:  req.Media.DisplayMeta.Width,
				Height: req.Media.DisplayMeta.Height,
			},
		}
	}

	if req.HiddenMetadata != nil {
		resa.HiddenMetadata = &res.HiddenMetadataRes{
			IsHidden:       req.HiddenMetadata.IsHidden,
			HiddenAt:       req.HiddenMetadata.HiddenAt,
			HiddenByUserID: req.HiddenMetadata.HiddenByUserID,
			Reason:         req.HiddenMetadata.Reason,
			IsGhostBanned:  req.HiddenMetadata.IsGhostBanned,
		}
	}

	if req.DeletedMetadata != nil {
		resa.DeletedMetadata = &res.DeletedMetadataRes{
			DeletedAt:       req.DeletedMetadata.DeletedAt,
			DeletedByUserID: req.DeletedMetadata.DeletedByUserID,
		}
	}

	return resa
}
