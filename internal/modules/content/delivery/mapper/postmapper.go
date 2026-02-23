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
		ID:          objectID,
		UserID:      r.UserID,
		Type:        r.Type,
		Content:     r.Content,
		Slug:        r.Slug,
		Privacy:     entity.PostPrivacy(r.Privacy),
		Status:      r.Status,
		IsPinned:    r.IsPinned,
		IsEdited:    r.IsEdited,
		Stats:       entity.PostStats(r.Stats),
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
	post.Privacy = entity.PostPrivacy(r.Privacy)
	post.Status = r.Status
	post.IsPinned = r.IsPinned
	post.IsEdited = true // Đã update thì thường cờ IsEdited sẽ là true
	post.Stats = entity.PostStats(r.Stats)
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
		ID:          post.ID.Hex(),
		UserID:      post.UserID,
		Type:        post.Type,
		Content:     post.Content,
		Slug:        post.Slug,
		Privacy:     res.PostPrivacyRes(post.Privacy),
		Status:      post.Status,
		IsPinned:    post.IsPinned,
		IsEdited:    post.IsEdited,
		Stats:       res.PostStatsRes(post.Stats),
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
