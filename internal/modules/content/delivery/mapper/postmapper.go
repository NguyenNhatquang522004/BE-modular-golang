package mapper

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- MAPPER YÊU CẦU 1: Request -> Entity (Dùng cho Create) ---
func ToEntityPost(r *req.PostReq) *entity.Post {
	if r == nil {
		return nil
	}

	// Xử lý ObjectID
	objectID := primitive.NewObjectID()
	if r.ID != "" {
		if oid, err := primitive.ObjectIDFromHex(r.ID); err == nil {
			objectID = oid
		}
	}

	post := &entity.Post{
		ID:      objectID,
		UserID:  r.UserID,
		Type:    r.Type,
		Content: r.Content,
		Slug:    r.Slug,
		// Map explicitly to avoid conversion errors
		Privacy: entity.PostPrivacy{
			Scope:        r.Privacy.Scope,
			AllowComment: r.Privacy.AllowComment,
			AllowShare:   r.Privacy.AllowShare,
		},
		Status:   r.Status,
		IsPinned: r.IsPinned,
		IsEdited: r.IsEdited,
		// Map explicitly to avoid conversion errors
		Stats: entity.PostStats{
			TotalReactions:   r.Stats.TotalReactions,
			Comments:         r.Stats.Comments,
			Shares:           r.Stats.Shares,
			Views:            r.Stats.Views,
			TopReactionTypes: r.Stats.TopReactionTypes,
			Like:             r.Stats.Like,
			Love:             r.Stats.Love,
			Haha:             r.Stats.Haha,
			Wow:              r.Stats.Wow,
			Sad:              r.Stats.Sad,
			Angry:            r.Stats.Angry,
		},
		Hashtags:    r.Hashtags,
		Mentions:    r.Mentions,
		PublishedAt: r.PublishedAt,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		DeletedAt:   r.DeletedAt,
	}

	if r.Context != nil {
		post.Context = &entity.PostContext{
			Type:     r.Context.Type,
			TargetID: r.Context.TargetID,
		}
	}

	if r.Summary != nil {
		post.Summary = &entity.PostSummary{
			FeelingIcon:       r.Summary.FeelingIcon,
			FeelingName:       r.Summary.FeelingName,
			LocationName:      r.Summary.LocationName,
			HasMedia:          r.Summary.HasMedia,
			MediaCount:        r.Summary.MediaCount,
			ThumbnailURL:      r.Summary.ThumbnailURL,
			BackgroundThemeID: r.Summary.BackgroundThemeID,
		}
	}

	return post
}

// --- MAPPER YÊU CẦU 2: Update Request -> Existing Entity ---
func UpdateToEntityPost(r *req.PostReq, post *entity.Post) {
	if r == nil || post == nil {
		return
	}

	// Best Practice: Update không can thiệp vào ID, UserID, và CreatedAt
	post.Type = r.Type
	post.Content = r.Content
	post.Slug = r.Slug

	// Map Privacy explicitly
	post.Privacy.Scope = r.Privacy.Scope
	post.Privacy.AllowComment = r.Privacy.AllowComment
	post.Privacy.AllowShare = r.Privacy.AllowShare

	post.Status = r.Status
	post.IsPinned = r.IsPinned
	post.IsEdited = true // Đã update thì thường cờ IsEdited sẽ là true

	// Map Stats explicitly
	post.Stats.TotalReactions = r.Stats.TotalReactions
	post.Stats.Comments = r.Stats.Comments
	post.Stats.Shares = r.Stats.Shares
	post.Stats.Views = r.Stats.Views
	post.Stats.TopReactionTypes = r.Stats.TopReactionTypes
	post.Stats.Like = r.Stats.Like
	post.Stats.Love = r.Stats.Love
	post.Stats.Haha = r.Stats.Haha
	post.Stats.Wow = r.Stats.Wow
	post.Stats.Sad = r.Stats.Sad
	post.Stats.Angry = r.Stats.Angry

	post.Hashtags = r.Hashtags
	post.Mentions = r.Mentions
	post.PublishedAt = r.PublishedAt
	post.UpdatedAt = r.UpdatedAt
	post.DeletedAt = r.DeletedAt

	if r.Context != nil {
		if post.Context == nil {
			post.Context = &entity.PostContext{}
		}
		post.Context.Type = r.Context.Type
		post.Context.TargetID = r.Context.TargetID
	}

	if r.Summary != nil {
		if post.Summary == nil {
			post.Summary = &entity.PostSummary{}
		}
		post.Summary.FeelingIcon = r.Summary.FeelingIcon
		post.Summary.FeelingName = r.Summary.FeelingName
		post.Summary.LocationName = r.Summary.LocationName
		post.Summary.HasMedia = r.Summary.HasMedia
		post.Summary.MediaCount = r.Summary.MediaCount
		post.Summary.ThumbnailURL = r.Summary.ThumbnailURL
		post.Summary.BackgroundThemeID = r.Summary.BackgroundThemeID
	}
}

// --- MAPPER RES: Entity -> Response (Chuẩn hóa trả về Client) ---
func ToResPost(post *entity.Post) *res.PostRes {
	if post == nil {
		return nil
	}

	result := &res.PostRes{
		ID:      post.ID.Hex(),
		UserID:  post.UserID,
		Type:    post.Type,
		Content: post.Content,
		Slug:    post.Slug,
		// Map explicitly
		Privacy: res.PostPrivacyRes{
			Scope:        post.Privacy.Scope,
			AllowComment: post.Privacy.AllowComment,
			AllowShare:   post.Privacy.AllowShare,
		},
		Status:   post.Status,
		IsPinned: post.IsPinned,
		IsEdited: post.IsEdited,
		// Map explicitly
		Stats: res.PostStatsRes{
			TotalReactions:   post.Stats.TotalReactions,
			Comments:         post.Stats.Comments,
			Shares:           post.Stats.Shares,
			Views:            post.Stats.Views,
			TopReactionTypes: post.Stats.TopReactionTypes,
			Like:             post.Stats.Like,
			Love:             post.Stats.Love,
			Haha:             post.Stats.Haha,
			Wow:              post.Stats.Wow,
			Sad:              post.Stats.Sad,
			Angry:            post.Stats.Angry,
		},
		Hashtags:    post.Hashtags,
		Mentions:    post.Mentions,
		PublishedAt: post.PublishedAt,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		DeletedAt:   post.DeletedAt,
	}

	if post.Context != nil {
		result.Context = &res.PostContextRes{
			Type:     post.Context.Type,
			TargetID: post.Context.TargetID,
		}
	}

	if post.Summary != nil {
		result.Summary = &res.PostSummaryRes{
			FeelingIcon:       post.Summary.FeelingIcon,
			FeelingName:       post.Summary.FeelingName,
			LocationName:      post.Summary.LocationName,
			HasMedia:          post.Summary.HasMedia,
			MediaCount:        post.Summary.MediaCount,
			ThumbnailURL:      post.Summary.ThumbnailURL,
			BackgroundThemeID: post.Summary.BackgroundThemeID,
		}
	}

	return result
}

// --- MAPPER YÊU CẦU 3: Chuyển trực tiếp từ Req -> Res (Theo đúng ý bạn yêu cầu) ---
func ReqToResPost(r *req.PostReq) *res.PostRes {
	if r == nil {
		return nil
	}

	result := &res.PostRes{
		ID:          r.ID,
		UserID:      r.UserID,
		Type:        r.Type,
		Content:     r.Content,
		Slug:        r.Slug,
		Privacy:     res.PostPrivacyRes(r.Privacy),
		Status:      r.Status,
		IsPinned:    r.IsPinned,
		IsEdited:    r.IsEdited,
		Stats:       res.PostStatsRes(r.Stats),
		Hashtags:    r.Hashtags,
		Mentions:    r.Mentions,
		PublishedAt: r.PublishedAt,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		DeletedAt:   r.DeletedAt,
	}

	if r.Context != nil {
		result.Context = &res.PostContextRes{
			Type:     r.Context.Type,
			TargetID: r.Context.TargetID,
		}
	}

	if r.Summary != nil {
		result.Summary = &res.PostSummaryRes{
			FeelingIcon:       r.Summary.FeelingIcon,
			FeelingName:       r.Summary.FeelingName,
			LocationName:      r.Summary.LocationName,
			HasMedia:          r.Summary.HasMedia,
			MediaCount:        r.Summary.MediaCount,
			ThumbnailURL:      r.Summary.ThumbnailURL,
			BackgroundThemeID: r.Summary.BackgroundThemeID,
		}
	}

	return result
}
