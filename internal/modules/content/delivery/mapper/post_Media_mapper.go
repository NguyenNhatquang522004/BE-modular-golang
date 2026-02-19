package mapper

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- 1. REQ TO ENTITY (CREATE) ---

func ToPostMediaEntity(r req.CreatePostMediaReq) (*entity.PostMedia, error) {
	postObjID, err := primitive.ObjectIDFromHex(r.PostID)
	if err != nil {
		return nil, err
	}

	ent := &entity.PostMedia{
		ID:     primitive.NewObjectID(),
		PostID: postObjID,
		Items:  make([]*entity.MediaItem, 0, len(r.Items)),
	}

	// Map từng Media Item
	for _, itemReq := range r.Items {
		ent.Items = append(ent.Items, mapMediaItemToEntity(itemReq))
	}

	return ent, nil
}

// --- 2. REQ TO ENTITY (UPDATE) ---

func UpdatePostMediaEntity(existingEnt *entity.PostMedia, r req.UpdatePostMediaReq) {
	// Với Media, chiến lược update thường là thay thế toàn bộ danh sách Items
	// để đảm bảo thứ tự (Order) và đồng bộ dữ liệu chính xác nhất.

	newItems := make([]*entity.MediaItem, 0, len(r.Items))

	for _, itemReq := range r.Items {
		newItems = append(newItems, mapMediaItemToEntity(itemReq))
	}

	existingEnt.Items = newItems
}

// --- Helper: Map Media Item (Dùng chung cho Create và Update) ---
func mapMediaItemToEntity(req *req.MediaItemReq) *entity.MediaItem {
	// Xử lý ID cho Item: Nếu req có ID hợp lệ thì giữ, nếu không thì tạo mới
	var itemID primitive.ObjectID
	if req.ID != "" {
		if id, err := primitive.ObjectIDFromHex(req.ID); err == nil {
			itemID = id
		} else {
			itemID = primitive.NewObjectID()
		}
	} else {
		itemID = primitive.NewObjectID()
	}

	itemEnt := &entity.MediaItem{
		ID:           itemID,
		MediaType:    req.MediaType,
		URL:          req.URL,
		ThumbnailURL: req.ThumbnailURL,
		Metadata: entity.MediaMetadata{
			Width:     req.Metadata.Width,
			Height:    req.Metadata.Height,
			Duration:  req.Metadata.Duration,
			SizeBytes: req.Metadata.SizeBytes,
			MimeType:  req.Metadata.MimeType,
		},
		Order:       req.Order,
		TaggedUsers: make([]entity.TaggedUser, 0, len(req.TaggedUsers)),
	}

	// Map Tagged Users
	for _, uReq := range req.TaggedUsers {
		itemEnt.TaggedUsers = append(itemEnt.TaggedUsers, entity.TaggedUser{
			UserID: uReq.UserID,
			Name:   uReq.Name,
			X:      uReq.X,
			Y:      uReq.Y,
		})
	}

	return itemEnt
}

// --- 3. ENTITY TO RES ---

func ToPostMediaRes(ent *entity.PostMedia) *res.PostMediaRes {
	if ent == nil {
		return nil
	}

	response := &res.PostMediaRes{
		ID:     ent.ID.Hex(),
		PostID: ent.PostID.Hex(),
		Items:  make([]*res.MediaItemRes, 0, len(ent.Items)),
	}

	for _, itemEnt := range ent.Items {
		itemRes := &res.MediaItemRes{
			ID:           itemEnt.ID.Hex(),
			MediaType:    itemEnt.MediaType,
			URL:          itemEnt.URL,
			ThumbnailURL: itemEnt.ThumbnailURL,
			Metadata: res.MediaMetadataRes{
				Width:     itemEnt.Metadata.Width,
				Height:    itemEnt.Metadata.Height,
				Duration:  itemEnt.Metadata.Duration,
				SizeBytes: itemEnt.Metadata.SizeBytes,
				MimeType:  itemEnt.Metadata.MimeType,
			},
			Order:       itemEnt.Order,
			TaggedUsers: make([]res.TaggedUserRes, 0, len(itemEnt.TaggedUsers)),
		}

		// Map Tagged Users Res
		for _, tagEnt := range itemEnt.TaggedUsers {
			itemRes.TaggedUsers = append(itemRes.TaggedUsers, res.TaggedUserRes{
				UserID: tagEnt.UserID,
				Name:   tagEnt.Name,
				X:      tagEnt.X,
				Y:      tagEnt.Y,
			})
		}

		response.Items = append(response.Items, itemRes)
	}

	return response
}
