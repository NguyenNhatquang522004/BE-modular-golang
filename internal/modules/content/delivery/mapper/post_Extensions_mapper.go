package mapper

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- 1. REQ TO ENTITY (CREATE) ---

func ToPostExtensionEntity(r req.CreatePostExtensionReq) (*entity.PostExtension, error) {
	postObjID, err := primitive.ObjectIDFromHex(r.PostID)
	if err != nil {
		return nil, err
	}

	ent := &entity.PostExtension{
		ID:     primitive.NewObjectID(),
		PostID: postObjID,
	}

	// Mapping optional parts
	if r.ShareData != nil {
		shareEnt, err := mapShareDataToEntity(r.ShareData)
		if err == nil {
			ent.ShareData = shareEnt
		}
	}
	if r.BackgroundData != nil {
		ent.BackgroundData = &entity.BackgroundData{
			ThemeID:   r.BackgroundData.ThemeID,
			TextColor: r.BackgroundData.TextColor,
		}
	}
	if r.QnAData != nil {
		ent.QnAData = &entity.QnAData{
			Question:   r.QnAData.Question,
			ButtonText: r.QnAData.ButtonText,
		}
	}
	if r.ActivityData != nil {
		ent.ActivityData = &entity.ActivityData{
			Type:       r.ActivityData.Type,
			ObjectID:   r.ActivityData.ObjectID,
			ObjectName: r.ActivityData.ObjectName,
		}
	}
	if r.LocationDetail != nil {
		ent.LocationDetail = &entity.LocationDetail{
			Type:        "Point", // Bắt buộc theo chuẩn GeoJSON
			Coordinates: r.LocationDetail.Coordinates,
			Address:     r.LocationDetail.Address,
			MapURL:      r.LocationDetail.MapURL,
		}
	}

	return ent, nil
}

// --- 2. REQ TO ENTITY (UPDATE) ---

func UpdatePostExtensionEntity(existingEnt *entity.PostExtension, r req.UpdatePostExtensionReq) {
	// Logic update: Chỉ thay thế các block dữ liệu được gửi lên (khác nil)

	if r.ShareData != nil {
		if shareEnt, err := mapShareDataToEntity(r.ShareData); err == nil {
			existingEnt.ShareData = shareEnt
		}
	}

	if r.BackgroundData != nil {
		existingEnt.BackgroundData = &entity.BackgroundData{
			ThemeID:   r.BackgroundData.ThemeID,
			TextColor: r.BackgroundData.TextColor,
		}
	}

	if r.QnAData != nil {
		existingEnt.QnAData = &entity.QnAData{
			Question:   r.QnAData.Question,
			ButtonText: r.QnAData.ButtonText,
		}
	}

	if r.ActivityData != nil {
		existingEnt.ActivityData = &entity.ActivityData{
			Type:       r.ActivityData.Type,
			ObjectID:   r.ActivityData.ObjectID,
			ObjectName: r.ActivityData.ObjectName,
		}
	}

	if r.LocationDetail != nil {
		existingEnt.LocationDetail = &entity.LocationDetail{
			Type:        "Point",
			Coordinates: r.LocationDetail.Coordinates,
			Address:     r.LocationDetail.Address,
			MapURL:      r.LocationDetail.MapURL,
		}
	}
}

// --- Helper: Map Share Data (Vì logic phức tạp có ID) ---
func mapShareDataToEntity(r *req.ShareDataReq) (*entity.ShareData, error) {
	origID, err := primitive.ObjectIDFromHex(r.OriginalPostID)
	if err != nil {
		return nil, err
	}
	parentID, err := primitive.ObjectIDFromHex(r.ParentPostID)
	if err != nil {
		return nil, err
	}

	return &entity.ShareData{
		OriginalPostID: origID,
		ParentPostID:   parentID,
		Snapshot: entity.ShareSnapshot{
			AuthorID:       r.Snapshot.AuthorID,
			AuthorName:     r.Snapshot.AuthorName,
			AuthorAvatar:   r.Snapshot.AuthorAvatar,
			ContentExcerpt: r.Snapshot.ContentExcerpt,
			MediaThumb:     r.Snapshot.MediaThumb,
			CreatedAt:      r.Snapshot.CreatedAt,
		},
	}, nil
}

// --- 3. ENTITY TO RES ---

func ToPostExtensionRes(ent *entity.PostExtension) *res.PostExtensionRes {
	if ent == nil {
		return nil
	}

	resObj := &res.PostExtensionRes{
		ID:     ent.ID.Hex(),
		PostID: ent.PostID.Hex(),
	}

	// Map ShareData
	if ent.ShareData != nil {
		resObj.ShareData = &res.ShareDataRes{
			OriginalPostID: ent.ShareData.OriginalPostID.Hex(),
			ParentPostID:   ent.ShareData.ParentPostID.Hex(),
			Snapshot: res.ShareSnapshotRes{
				AuthorID:       ent.ShareData.Snapshot.AuthorID,
				AuthorName:     ent.ShareData.Snapshot.AuthorName,
				AuthorAvatar:   ent.ShareData.Snapshot.AuthorAvatar,
				ContentExcerpt: ent.ShareData.Snapshot.ContentExcerpt,
				MediaThumb:     ent.ShareData.Snapshot.MediaThumb,
				CreatedAt:      ent.ShareData.Snapshot.CreatedAt,
			},
		}
	}

	// Map BackgroundData
	if ent.BackgroundData != nil {
		resObj.BackgroundData = &res.BackgroundDataRes{
			ThemeID:   ent.BackgroundData.ThemeID,
			TextColor: ent.BackgroundData.TextColor,
		}
	}

	// Map QnAData
	if ent.QnAData != nil {
		resObj.QnAData = &res.QnADataRes{
			Question:   ent.QnAData.Question,
			ButtonText: ent.QnAData.ButtonText,
		}
	}

	// Map ActivityData
	if ent.ActivityData != nil {
		resObj.ActivityData = &res.ActivityDataRes{
			Type:       ent.ActivityData.Type,
			ObjectID:   ent.ActivityData.ObjectID,
			ObjectName: ent.ActivityData.ObjectName,
		}
	}

	// Map LocationDetail
	if ent.LocationDetail != nil {
		resObj.LocationDetail = &res.LocationDetailRes{
			Type:        ent.LocationDetail.Type,
			Coordinates: ent.LocationDetail.Coordinates,
			Address:     ent.LocationDetail.Address,
			MapURL:      ent.LocationDetail.MapURL,
		}
	}

	return resObj
}
