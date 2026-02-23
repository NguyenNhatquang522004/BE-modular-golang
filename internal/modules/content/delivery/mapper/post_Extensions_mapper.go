package mapper

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- 1. Request -> Entity (Tạo mới) ---
func ToPostExtensionEntityPostExtension(r *req.PostExtensionReq) *entity.PostExtension {
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

	ext := &entity.PostExtension{
		ID:        objectID,
		PostID:    postID,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		DeletedAt: r.DeletedAt,
	}

	if r.ShareData != nil {
		originalID, _ := primitive.ObjectIDFromHex(r.ShareData.OriginalPostID)
		parentID, _ := primitive.ObjectIDFromHex(r.ShareData.ParentPostID)

		ext.ShareData = &entity.ShareData{
			OriginalPostID: originalID,
			ParentPostID:   parentID,
			Snapshot: entity.ShareSnapshot{
				AuthorID:       r.ShareData.Snapshot.AuthorID,
				AuthorName:     r.ShareData.Snapshot.AuthorName,
				AuthorAvatar:   r.ShareData.Snapshot.AuthorAvatar,
				ContentExcerpt: r.ShareData.Snapshot.ContentExcerpt,
				MediaThumb:     r.ShareData.Snapshot.MediaThumb,
				CreatedAt:      r.ShareData.Snapshot.CreatedAt,
			},
		}
	}

	if r.BackgroundData != nil {
		ext.BackgroundData = &entity.BackgroundData{
			ThemeID:   r.BackgroundData.ThemeID,
			TextColor: r.BackgroundData.TextColor,
		}
	}

	if r.QnAData != nil {
		ext.QnAData = &entity.QnAData{
			Question:   r.QnAData.Question,
			ButtonText: r.QnAData.ButtonText,
		}
	}

	if r.ActivityData != nil {
		ext.ActivityData = &entity.ActivityData{
			Type:       r.ActivityData.Type,
			ObjectID:   r.ActivityData.ObjectID,
			ObjectName: r.ActivityData.ObjectName,
		}
	}

	if r.LocationDetail != nil {
		ext.LocationDetail = &entity.LocationDetail{
			Type:        r.LocationDetail.Type,
			Coordinates: r.LocationDetail.Coordinates,
			Address:     r.LocationDetail.Address,
			MapURL:      r.LocationDetail.MapURL,
		}
	}

	return ext
}

// --- 2. Update Request -> Existing Entity ---
func UpdatePostExtensionEntityPostExtension(r *req.PostExtensionReq, ext *entity.PostExtension) {
	if r == nil || ext == nil {
		return
	}

	// Bỏ qua Update ID, PostID, CreatedAt
	ext.UpdatedAt = r.UpdatedAt
	ext.DeletedAt = r.DeletedAt

	if r.ShareData != nil {
		originalID, _ := primitive.ObjectIDFromHex(r.ShareData.OriginalPostID)
		parentID, _ := primitive.ObjectIDFromHex(r.ShareData.ParentPostID)

		ext.ShareData = &entity.ShareData{
			OriginalPostID: originalID,
			ParentPostID:   parentID,
			Snapshot: entity.ShareSnapshot{
				AuthorID:       r.ShareData.Snapshot.AuthorID,
				AuthorName:     r.ShareData.Snapshot.AuthorName,
				AuthorAvatar:   r.ShareData.Snapshot.AuthorAvatar,
				ContentExcerpt: r.ShareData.Snapshot.ContentExcerpt,
				MediaThumb:     r.ShareData.Snapshot.MediaThumb,
				CreatedAt:      r.ShareData.Snapshot.CreatedAt,
			},
		}
	}

	if r.BackgroundData != nil {
		ext.BackgroundData = &entity.BackgroundData{
			ThemeID:   r.BackgroundData.ThemeID,
			TextColor: r.BackgroundData.TextColor,
		}
	}

	if r.QnAData != nil {
		ext.QnAData = &entity.QnAData{
			Question:   r.QnAData.Question,
			ButtonText: r.QnAData.ButtonText,
		}
	}

	if r.ActivityData != nil {
		ext.ActivityData = &entity.ActivityData{
			Type:       r.ActivityData.Type,
			ObjectID:   r.ActivityData.ObjectID,
			ObjectName: r.ActivityData.ObjectName,
		}
	}

	if r.LocationDetail != nil {
		ext.LocationDetail = &entity.LocationDetail{
			Type:        r.LocationDetail.Type,
			Coordinates: r.LocationDetail.Coordinates,
			Address:     r.LocationDetail.Address,
			MapURL:      r.LocationDetail.MapURL,
		}
	}
}

// --- 3. Entity -> Response ---
func ToPostExtensionResPostExtension(ext *entity.PostExtension) *res.PostExtensionRes {
	if ext == nil {
		return nil
	}

	result := &res.PostExtensionRes{
		ID:        ext.ID.Hex(),
		PostID:    ext.PostID.Hex(),
		CreatedAt: ext.CreatedAt,
		UpdatedAt: ext.UpdatedAt,
		DeletedAt: ext.DeletedAt,
	}

	if ext.ShareData != nil {
		result.ShareData = &res.ShareDataRes{
			OriginalPostID: ext.ShareData.OriginalPostID.Hex(),
			ParentPostID:   ext.ShareData.ParentPostID.Hex(),
			Snapshot: res.ShareSnapshotRes{
				AuthorID:       ext.ShareData.Snapshot.AuthorID,
				AuthorName:     ext.ShareData.Snapshot.AuthorName,
				AuthorAvatar:   ext.ShareData.Snapshot.AuthorAvatar,
				ContentExcerpt: ext.ShareData.Snapshot.ContentExcerpt,
				MediaThumb:     ext.ShareData.Snapshot.MediaThumb,
				CreatedAt:      ext.ShareData.Snapshot.CreatedAt,
			},
		}
	}

	if ext.BackgroundData != nil {
		result.BackgroundData = &res.BackgroundDataRes{
			ThemeID:   ext.BackgroundData.ThemeID,
			TextColor: ext.BackgroundData.TextColor,
		}
	}

	if ext.QnAData != nil {
		result.QnAData = &res.QnADataRes{
			Question:   ext.QnAData.Question,
			ButtonText: ext.QnAData.ButtonText,
		}
	}

	if ext.ActivityData != nil {
		result.ActivityData = &res.ActivityDataRes{
			Type:       ext.ActivityData.Type,
			ObjectID:   ext.ActivityData.ObjectID,
			ObjectName: ext.ActivityData.ObjectName,
		}
	}

	if ext.LocationDetail != nil {
		result.LocationDetail = &res.LocationDetailRes{
			Type:        ext.LocationDetail.Type,
			Coordinates: ext.LocationDetail.Coordinates,
			Address:     ext.LocationDetail.Address,
			MapURL:      ext.LocationDetail.MapURL,
		}
	}

	return result
}

// --- 4. Request -> Response ---
func ReqToPostExtensionResPostExtension(r *req.PostExtensionReq) *res.PostExtensionRes {
	if r == nil {
		return nil
	}

	result := &res.PostExtensionRes{
		ID:        r.ID,
		PostID:    r.PostID,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		DeletedAt: r.DeletedAt,
	}

	if r.ShareData != nil {
		result.ShareData = &res.ShareDataRes{
			OriginalPostID: r.ShareData.OriginalPostID,
			ParentPostID:   r.ShareData.ParentPostID,
			Snapshot: res.ShareSnapshotRes{
				AuthorID:       r.ShareData.Snapshot.AuthorID,
				AuthorName:     r.ShareData.Snapshot.AuthorName,
				AuthorAvatar:   r.ShareData.Snapshot.AuthorAvatar,
				ContentExcerpt: r.ShareData.Snapshot.ContentExcerpt,
				MediaThumb:     r.ShareData.Snapshot.MediaThumb,
				CreatedAt:      r.ShareData.Snapshot.CreatedAt,
			},
		}
	}

	if r.BackgroundData != nil {
		result.BackgroundData = &res.BackgroundDataRes{
			ThemeID:   r.BackgroundData.ThemeID,
			TextColor: r.BackgroundData.TextColor,
		}
	}

	if r.QnAData != nil {
		result.QnAData = &res.QnADataRes{
			Question:   r.QnAData.Question,
			ButtonText: r.QnAData.ButtonText,
		}
	}

	if r.ActivityData != nil {
		result.ActivityData = &res.ActivityDataRes{
			Type:       r.ActivityData.Type,
			ObjectID:   r.ActivityData.ObjectID,
			ObjectName: r.ActivityData.ObjectName,
		}
	}

	if r.LocationDetail != nil {
		result.LocationDetail = &res.LocationDetailRes{
			Type:        r.LocationDetail.Type,
			Coordinates: r.LocationDetail.Coordinates,
			Address:     r.LocationDetail.Address,
			MapURL:      r.LocationDetail.MapURL,
		}
	}

	return result
}
