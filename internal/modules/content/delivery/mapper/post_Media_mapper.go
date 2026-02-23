package mapper

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- 1. Request -> Entity (Tạo mới) ---
func ToPostMediaEntityPostMedia(r *req.PostMediaReq) *entity.PostMedia {
	if r == nil {
		return nil
	}

	objectID := primitive.NewObjectID()
	if r.ID != "" {
		if oid, err := primitive.ObjectIDFromHex(r.ID); err == nil {
			objectID = oid
		}
	}

	postID, _ := primitive.ObjectIDFromHex(r.PostID)

	media := &entity.PostMedia{
		ID:        objectID,
		PostID:    postID,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		DeletedAt: r.DeletedAt,
	}

	if len(r.Items) > 0 {
		media.Items = make([]*entity.MediaItem, 0, len(r.Items))
		for _, itemReq := range r.Items {
			if itemReq == nil {
				continue
			}

			itemID := primitive.NewObjectID()
			if itemReq.ID != "" {
				if oid, err := primitive.ObjectIDFromHex(itemReq.ID); err == nil {
					itemID = oid
				}
			}

			itemEntity := &entity.MediaItem{
				ID:           itemID,
				MediaType:    itemReq.MediaType,
				URL:          itemReq.URL,
				ThumbnailURL: itemReq.ThumbnailURL,
				Order:        itemReq.Order,
				Metadata: entity.MediaMetadata{
					Width:     itemReq.Metadata.Width,
					Height:    itemReq.Metadata.Height,
					Duration:  itemReq.Metadata.Duration,
					SizeBytes: itemReq.Metadata.SizeBytes,
					MimeType:  itemReq.Metadata.MimeType,
				},
			}

			if len(itemReq.TaggedUsers) > 0 {
				itemEntity.TaggedUsers = make([]entity.TaggedUser, 0, len(itemReq.TaggedUsers))
				for _, tuReq := range itemReq.TaggedUsers {
					itemEntity.TaggedUsers = append(itemEntity.TaggedUsers, entity.TaggedUser{
						UserID: tuReq.UserID,
						Name:   tuReq.Name,
						X:      tuReq.X,
						Y:      tuReq.Y,
					})
				}
			}

			media.Items = append(media.Items, itemEntity)
		}
	}

	return media
}

// --- 2. Update Request -> Existing Entity ---
func UpdatePostMediaEntityPostMedia(r *req.PostMediaReq, media *entity.PostMedia) {
	if r == nil || media == nil {
		return
	}

	// Bỏ qua Update ID, PostID, CreatedAt
	media.UpdatedAt = r.UpdatedAt
	media.DeletedAt = r.DeletedAt

	// Cập nhật lại toàn bộ danh sách Items (Replace behavior)
	if r.Items != nil {
		media.Items = make([]*entity.MediaItem, 0, len(r.Items))
		for _, itemReq := range r.Items {
			if itemReq == nil {
				continue
			}

			itemID := primitive.NewObjectID()
			if itemReq.ID != "" {
				if oid, err := primitive.ObjectIDFromHex(itemReq.ID); err == nil {
					itemID = oid
				}
			}

			itemEntity := &entity.MediaItem{
				ID:           itemID,
				MediaType:    itemReq.MediaType,
				URL:          itemReq.URL,
				ThumbnailURL: itemReq.ThumbnailURL,
				Order:        itemReq.Order,
				Metadata: entity.MediaMetadata{
					Width:     itemReq.Metadata.Width,
					Height:    itemReq.Metadata.Height,
					Duration:  itemReq.Metadata.Duration,
					SizeBytes: itemReq.Metadata.SizeBytes,
					MimeType:  itemReq.Metadata.MimeType,
				},
			}

			if len(itemReq.TaggedUsers) > 0 {
				itemEntity.TaggedUsers = make([]entity.TaggedUser, 0, len(itemReq.TaggedUsers))
				for _, tuReq := range itemReq.TaggedUsers {
					itemEntity.TaggedUsers = append(itemEntity.TaggedUsers, entity.TaggedUser{
						UserID: tuReq.UserID,
						Name:   tuReq.Name,
						X:      tuReq.X,
						Y:      tuReq.Y,
					})
				}
			}

			media.Items = append(media.Items, itemEntity)
		}
	}
}

// --- 3. Entity -> Response ---
func ToPostMediaResPostMedia(media *entity.PostMedia) *res.PostMediaRes {
	if media == nil {
		return nil
	}

	result := &res.PostMediaRes{
		ID:        media.ID.Hex(),
		PostID:    media.PostID.Hex(),
		CreatedAt: media.CreatedAt,
		UpdatedAt: media.UpdatedAt,
		DeletedAt: media.DeletedAt,
	}

	if len(media.Items) > 0 {
		result.Items = make([]*res.MediaItemRes, 0, len(media.Items))
		for _, itemEntity := range media.Items {
			if itemEntity == nil {
				continue
			}

			itemRes := &res.MediaItemRes{
				ID:           itemEntity.ID.Hex(),
				MediaType:    itemEntity.MediaType,
				URL:          itemEntity.URL,
				ThumbnailURL: itemEntity.ThumbnailURL,
				Order:        itemEntity.Order,
				Metadata: res.MediaMetadataRes{
					Width:     itemEntity.Metadata.Width,
					Height:    itemEntity.Metadata.Height,
					Duration:  itemEntity.Metadata.Duration,
					SizeBytes: itemEntity.Metadata.SizeBytes,
					MimeType:  itemEntity.Metadata.MimeType,
				},
			}

			if len(itemEntity.TaggedUsers) > 0 {
				itemRes.TaggedUsers = make([]res.TaggedUserRes, 0, len(itemEntity.TaggedUsers))
				for _, tuEntity := range itemEntity.TaggedUsers {
					itemRes.TaggedUsers = append(itemRes.TaggedUsers, res.TaggedUserRes{
						UserID: tuEntity.UserID,
						Name:   tuEntity.Name,
						X:      tuEntity.X,
						Y:      tuEntity.Y,
					})
				}
			}

			result.Items = append(result.Items, itemRes)
		}
	}

	return result
}

// --- 4. Request -> Response ---
func ReqToPostMediaResPostMedia(r *req.PostMediaReq) *res.PostMediaRes {
	if r == nil {
		return nil
	}

	result := &res.PostMediaRes{
		ID:        r.ID,
		PostID:    r.PostID,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		DeletedAt: r.DeletedAt,
	}

	if len(r.Items) > 0 {
		result.Items = make([]*res.MediaItemRes, 0, len(r.Items))
		for _, itemReq := range r.Items {
			if itemReq == nil {
				continue
			}

			itemRes := &res.MediaItemRes{
				ID:           itemReq.ID,
				MediaType:    itemReq.MediaType,
				URL:          itemReq.URL,
				ThumbnailURL: itemReq.ThumbnailURL,
				Order:        itemReq.Order,
				Metadata: res.MediaMetadataRes{
					Width:     itemReq.Metadata.Width,
					Height:    itemReq.Metadata.Height,
					Duration:  itemReq.Metadata.Duration,
					SizeBytes: itemReq.Metadata.SizeBytes,
					MimeType:  itemReq.Metadata.MimeType,
				},
			}

			if len(itemReq.TaggedUsers) > 0 {
				itemRes.TaggedUsers = make([]res.TaggedUserRes, 0, len(itemReq.TaggedUsers))
				for _, tuReq := range itemReq.TaggedUsers {
					itemRes.TaggedUsers = append(itemRes.TaggedUsers, res.TaggedUserRes{
						UserID: tuReq.UserID,
						Name:   tuReq.Name,
						X:      tuReq.X,
						Y:      tuReq.Y,
					})
				}
			}

			result.Items = append(result.Items, itemRes)
		}
	}

	return result
}
