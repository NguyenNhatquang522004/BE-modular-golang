package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 1. ToEntityMediaAsset: Chuyển từ Req -> Entity để Insert
func ToEntityMediaAsset(r *req.MediaAssetReq) (*entity.MediaAsset, error) {
	if r == nil {
		return nil, nil
	}

	e := &entity.MediaAsset{
		UserID:        r.UserID,
		StorageFileID: r.StorageFileID,
		OriginalURL:   r.OriginalURL,
		ThumbnailURL:  r.ThumbnailURL,
		AssetType:     r.AssetType,
		Caption:       r.Caption,
		Hashtags:      r.Hashtags,
		CommentCount:  r.CommentCount,
		Order:         r.Order,
		DeletedAt:     r.DeletedAt,
		Metadata: entity.MediaMetadata{
			Width:     r.Metadata.Width,
			Height:    r.Metadata.Height,
			Duration:  r.Metadata.Duration,
			SizeBytes: r.Metadata.SizeBytes,
			MimeType:  r.Metadata.MimeType,
		},
		ReactionsCount: entity.MediaReactionStats{
			Total: r.ReactionsCount.Total,
			Like:  r.ReactionsCount.Like,
			Love:  r.ReactionsCount.Love,
		},
	}

	// Xử lý ID
	if r.ID != "" {
		id, err := primitive.ObjectIDFromHex(r.ID)
		if err == nil {
			e.ID = id
		}
	} else {
		e.ID = primitive.NewObjectID()
	}

	// Xử lý các ID liên kết (Album, Post, Group)
	if r.AlbumID != "" {
		if id, err := primitive.ObjectIDFromHex(r.AlbumID); err == nil {
			e.AlbumID = id
		}
	}
	if r.PostID != "" {
		if id, err := primitive.ObjectIDFromHex(r.PostID); err == nil {
			e.PostID = id
		}
	}
	if r.GroupID != "" {
		if id, err := primitive.ObjectIDFromHex(r.GroupID); err == nil {
			e.GroupID = id
		}
	}

	// Xử lý mảng TaggedUsers
	if len(r.TaggedUsers) > 0 {
		e.TaggedUsers = make([]entity.MediaTag, len(r.TaggedUsers))
		for i, tag := range r.TaggedUsers {
			e.TaggedUsers[i] = entity.MediaTag{
				UserID: tag.UserID,
				Name:   tag.Name,
				Status: tag.Status,
				Position: entity.TagPosition{
					X: tag.Position.X,
					Y: tag.Position.Y,
				},
			}
		}
	}

	// Xử lý Privacy (Pointer)
	if r.Privacy != nil {
		e.Privacy = &entity.MediaPrivacy{
			Level:            r.Privacy.Level,
			InheritFromAlbum: r.Privacy.InheritFromAlbum,
		}
	}

	// Xử lý Timestamps
	now := time.Now()
	if r.CreatedAt != nil {
		e.CreatedAt = *r.CreatedAt
	} else {
		e.CreatedAt = now
	}
	if r.UpdatedAt != nil {
		e.UpdatedAt = *r.UpdatedAt
	} else {
		e.UpdatedAt = now
	}

	return e, nil
}

// 2. UpdateToEntityMediaAsset: Cập nhật Partial Update
func UpdateToEntityMediaAsset(r *req.UpdateMediaAssetReq, e *entity.MediaAsset) {
	if r == nil || e == nil {
		return
	}

	if r.AlbumID != nil && *r.AlbumID != "" {
		if id, err := primitive.ObjectIDFromHex(*r.AlbumID); err == nil {
			e.AlbumID = id
		}
	}
	if r.PostID != nil && *r.PostID != "" {
		if id, err := primitive.ObjectIDFromHex(*r.PostID); err == nil {
			e.PostID = id
		}
	}
	if r.GroupID != nil && *r.GroupID != "" {
		if id, err := primitive.ObjectIDFromHex(*r.GroupID); err == nil {
			e.GroupID = id
		}
	}
	if r.StorageFileID != nil {
		e.StorageFileID = *r.StorageFileID
	}
	if r.OriginalURL != nil {
		e.OriginalURL = *r.OriginalURL
	}
	if r.ThumbnailURL != nil {
		e.ThumbnailURL = *r.ThumbnailURL
	}
	if r.AssetType != nil {
		e.AssetType = *r.AssetType
	}
	if r.Metadata != nil {
		e.Metadata.Width = r.Metadata.Width
		e.Metadata.Height = r.Metadata.Height
		e.Metadata.Duration = r.Metadata.Duration
		e.Metadata.SizeBytes = r.Metadata.SizeBytes
		e.Metadata.MimeType = r.Metadata.MimeType
	}
	if r.Caption != nil {
		e.Caption = *r.Caption
	}
	if r.Hashtags != nil {
		e.Hashtags = r.Hashtags
	}
	if r.TaggedUsers != nil {
		e.TaggedUsers = make([]entity.MediaTag, len(r.TaggedUsers))
		for i, tag := range r.TaggedUsers {
			e.TaggedUsers[i] = entity.MediaTag{
				UserID:   tag.UserID,
				Name:     tag.Name,
				Status:   tag.Status,
				Position: entity.TagPosition{X: tag.Position.X, Y: tag.Position.Y},
			}
		}
	}
	if r.ReactionsCount != nil {
		e.ReactionsCount.Total = r.ReactionsCount.Total
		e.ReactionsCount.Like = r.ReactionsCount.Like
		e.ReactionsCount.Love = r.ReactionsCount.Love
	}
	if r.CommentCount != nil {
		e.CommentCount = *r.CommentCount
	}
	if r.Order != nil {
		e.Order = *r.Order
	}
	if r.Privacy != nil {
		e.Privacy = &entity.MediaPrivacy{
			Level:            r.Privacy.Level,
			InheritFromAlbum: r.Privacy.InheritFromAlbum,
		}
	}
	if r.DeletedAt != nil {
		e.DeletedAt = r.DeletedAt
	}

	e.UpdatedAt = time.Now()
}

// 3. ReqToResMediaAsset: Chuyển thẳng từ Req -> Res
func ReqToResMediaAsset(r *req.MediaAssetReq) *res.MediaAssetRes {
	if r == nil {
		return nil
	}

	response := &res.MediaAssetRes{
		ID:            r.ID,
		UserID:        r.UserID,
		AlbumID:       r.AlbumID,
		PostID:        r.PostID,
		GroupID:       r.GroupID,
		StorageFileID: r.StorageFileID,
		OriginalURL:   r.OriginalURL,
		ThumbnailURL:  r.ThumbnailURL,
		AssetType:     r.AssetType,
		Caption:       r.Caption,
		Hashtags:      r.Hashtags,
		CommentCount:  r.CommentCount,
		Order:         r.Order,
		DeletedAt:     r.DeletedAt,
		Metadata: res.MediaMetadataRes{
			Width:     r.Metadata.Width,
			Height:    r.Metadata.Height,
			Duration:  r.Metadata.Duration,
			SizeBytes: r.Metadata.SizeBytes,
			MimeType:  r.Metadata.MimeType,
		},
		ReactionsCount: res.MediaReactionStatsRes{
			Total: r.ReactionsCount.Total,
			Like:  r.ReactionsCount.Like,
			Love:  r.ReactionsCount.Love,
		},
	}

	if len(r.TaggedUsers) > 0 {
		response.TaggedUsers = make([]res.MediaTagRes, len(r.TaggedUsers))
		for i, tag := range r.TaggedUsers {
			response.TaggedUsers[i] = res.MediaTagRes{
				UserID:   tag.UserID,
				Name:     tag.Name,
				Status:   tag.Status,
				Position: res.TagPositionRes{X: tag.Position.X, Y: tag.Position.Y},
			}
		}
	}

	if r.Privacy != nil {
		response.Privacy = &res.MediaPrivacyRes{
			Level:            r.Privacy.Level,
			InheritFromAlbum: r.Privacy.InheritFromAlbum,
		}
	}

	now := time.Now()
	if r.CreatedAt != nil {
		response.CreatedAt = *r.CreatedAt
	} else {
		response.CreatedAt = now
	}
	if r.UpdatedAt != nil {
		response.UpdatedAt = *r.UpdatedAt
	} else {
		response.UpdatedAt = now
	}

	return response
}

// 4. EntityToResMediaAsset: Best Practice, dùng khi Query từ DB ra
func EntityToResMediaAsset(e *entity.MediaAsset) *res.MediaAssetRes {
	if e == nil {
		return nil
	}

	response := &res.MediaAssetRes{
		ID:            e.ID.Hex(),
		UserID:        e.UserID,
		StorageFileID: e.StorageFileID,
		OriginalURL:   e.OriginalURL,
		ThumbnailURL:  e.ThumbnailURL,
		AssetType:     e.AssetType,
		Caption:       e.Caption,
		Hashtags:      e.Hashtags,
		CommentCount:  e.CommentCount,
		Order:         e.Order,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
		DeletedAt:     e.DeletedAt,
		Metadata: res.MediaMetadataRes{
			Width:     e.Metadata.Width,
			Height:    e.Metadata.Height,
			Duration:  e.Metadata.Duration,
			SizeBytes: e.Metadata.SizeBytes,
			MimeType:  e.Metadata.MimeType,
		},
		ReactionsCount: res.MediaReactionStatsRes{
			Total: e.ReactionsCount.Total,
			Like:  e.ReactionsCount.Like,
			Love:  e.ReactionsCount.Love,
		},
	}

	// Chỉ chuyển sang Hex string nếu ObjectID tồn tại (khác zero value)
	if !e.AlbumID.IsZero() {
		response.AlbumID = e.AlbumID.Hex()
	}
	if !e.PostID.IsZero() {
		response.PostID = e.PostID.Hex()
	}
	if !e.GroupID.IsZero() {
		response.GroupID = e.GroupID.Hex()
	}

	if len(e.TaggedUsers) > 0 {
		response.TaggedUsers = make([]res.MediaTagRes, len(e.TaggedUsers))
		for i, tag := range e.TaggedUsers {
			response.TaggedUsers[i] = res.MediaTagRes{
				UserID:   tag.UserID,
				Name:     tag.Name,
				Status:   tag.Status,
				Position: res.TagPositionRes{X: tag.Position.X, Y: tag.Position.Y},
			}
		}
	}

	if e.Privacy != nil {
		response.Privacy = &res.MediaPrivacyRes{
			Level:            e.Privacy.Level,
			InheritFromAlbum: e.Privacy.InheritFromAlbum,
		}
	}

	return response
}
