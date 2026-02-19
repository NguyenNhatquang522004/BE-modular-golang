package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PostMapper struct để gom nhóm các phương thức (Optional, có thể dùng function thường)
type PostMapper struct{}

// 1. ToEntity: Chuyển từ Create Request sang Entity
func (m *PostMapper) ToEntity(req *req.CreatePostRequest, userID string) *entity.Post {
	now := time.Now()

	// Khởi tạo Post Context nếu có
	var postContext *entity.PostContext
	if req.Context != nil {
		postContext = &entity.PostContext{
			Type:     req.Context.Type,
			TargetID: req.Context.TargetID,
		}
	}

	// Khởi tạo Post Summary nếu có
	var postSummary *entity.PostSummary
	if req.Summary != nil {
		postSummary = &entity.PostSummary{
			FeelingIcon:       req.Summary.FeelingIcon,
			FeelingName:       req.Summary.FeelingName,
			LocationName:      req.Summary.LocationName,
			HasMedia:          req.Summary.HasMedia,
			MediaCount:        req.Summary.MediaCount,
			ThumbnailURL:      req.Summary.ThumbnailURL,
			BackgroundThemeID: req.Summary.BackgroundThemeID,
		}
	}

	// Mặc định Status
	status := enum.StatusPublished // Giả sử mặc định là Published

	return &entity.Post{
		ID:      primitive.NewObjectID(), // Tạo ID mới ngay tại đây
		UserID:  userID,
		Type:    req.Type,
		Context: postContext,
		Content: req.Content,
		Slug:    "", // Slug sẽ được generate ở Service layer
		Summary: postSummary,
		Privacy: entity.PostPrivacy{
			Scope:        req.Privacy.Scope,
			AllowComment: req.Privacy.AllowComment,
			AllowShare:   req.Privacy.AllowShare,
		},
		Status:      status,
		IsPinned:    false,
		IsEdited:    false,
		Stats:       entity.PostStats{}, // Mặc định stats bằng 0
		Hashtags:    req.Hashtags,
		Mentions:    req.Mentions,
		PublishedAt: req.PublishedAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// 2. UpdateToEntity: Cập nhật Entity hiện tại dựa trên Update Request
func (m *PostMapper) UpdateToEntity(post *entity.Post, req *req.UpdatePostRequest) {
	// Chỉ update những trường client gửi lên
	if req.Content != "" {
		post.Content = req.Content
		post.IsEdited = true
	}

	if req.Summary != nil {
		if post.Summary == nil {
			post.Summary = &entity.PostSummary{}
		}
		// Update từng trường hoặc ghi đè tùy logic, ở đây mình ghi đè các trường metadata
		post.Summary.FeelingIcon = req.Summary.FeelingIcon
		post.Summary.FeelingName = req.Summary.FeelingName
		post.Summary.LocationName = req.Summary.LocationName
		post.Summary.HasMedia = req.Summary.HasMedia
		post.Summary.MediaCount = req.Summary.MediaCount
		post.Summary.ThumbnailURL = req.Summary.ThumbnailURL
		post.Summary.BackgroundThemeID = req.Summary.BackgroundThemeID
	}

	if req.Privacy != nil {
		post.Privacy.Scope = req.Privacy.Scope
		post.Privacy.AllowComment = req.Privacy.AllowComment
		post.Privacy.AllowShare = req.Privacy.AllowShare
	}

	if req.Status != nil {
		post.Status = *req.Status
	}

	// Xử lý field boolean pointer
	if req.IsPinned != nil {
		post.IsPinned = *req.IsPinned
	}

	if req.Hashtags != nil {
		post.Hashtags = req.Hashtags
	}

	if req.Mentions != nil {
		post.Mentions = req.Mentions
	}

	post.UpdatedAt = time.Now()
}

// 3. ToResponse: Chuyển từ Entity sang Response DTO
func (m *PostMapper) ToResponse(post *entity.Post) *res.PostResponse {
	if post == nil {
		return nil
	}

	// Map Context
	var ctxRes *res.PostContextResponse
	if post.Context != nil {
		ctxRes = &res.PostContextResponse{
			Type:     post.Context.Type,
			TargetID: post.Context.TargetID,
		}
	}

	// Map Summary
	var sumRes *res.PostSummaryResponse
	if post.Summary != nil {
		sumRes = &res.PostSummaryResponse{
			FeelingIcon:       post.Summary.FeelingIcon,
			FeelingName:       post.Summary.FeelingName,
			LocationName:      post.Summary.LocationName,
			HasMedia:          post.Summary.HasMedia,
			MediaCount:        post.Summary.MediaCount,
			ThumbnailURL:      post.Summary.ThumbnailURL,
			BackgroundThemeID: post.Summary.BackgroundThemeID,
		}
	}

	return &res.PostResponse{
		ID:      post.ID.Hex(), // Convert ObjectID to Hex String
		UserID:  post.UserID,
		Type:    post.Type,
		Context: ctxRes,
		Content: post.Content,
		Slug:    post.Slug,
		Summary: sumRes,
		Privacy: res.PostPrivacyResponse{
			Scope:        post.Privacy.Scope,
			AllowComment: post.Privacy.AllowComment,
			AllowShare:   post.Privacy.AllowShare,
		},
		Status:   post.Status,
		IsPinned: post.IsPinned,
		IsEdited: post.IsEdited,
		Stats: res.PostStatsResponse{
			TotalReactions:   post.Stats.TotalReactions,
			Comments:         post.Stats.Comments,
			Shares:           post.Stats.Shares,
			Views:            post.Stats.Views,
			TopReactionTypes: post.Stats.TopReactionTypes,
		},
		Hashtags:    post.Hashtags,
		Mentions:    post.Mentions,
		PublishedAt: post.PublishedAt,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	}
}
