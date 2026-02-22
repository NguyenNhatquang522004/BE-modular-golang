package mapper

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ToEntityPageFollower chuyển đổi từ Request DTO sang MongoDB Entity
func ToEntityPageFollower(r *req.PageFollowerReq) (*entity.PageFollower, error) {
	if r == nil {
		return nil, nil
	}

	ent := &entity.PageFollower{
		UserID: r.UserID,
	}

	// Chỉ map ID nếu client có truyền lên (trường hợp import/migrate)
	if r.ID != "" {
		id, err := primitive.ObjectIDFromHex(r.ID)
		if err != nil {
			return nil, err
		}
		ent.ID = id
	}

	// Map PageID từ string sang ObjectID
	if r.PageID != "" {
		pageID, err := primitive.ObjectIDFromHex(r.PageID)
		if err != nil {
			return nil, err
		}
		ent.PageID = pageID
	}

	// Ánh xạ Settings nếu tồn tại
	if r.Settings != nil {
		ent.Settings = &entity.FollowerSettings{
			NotificationLevel: r.Settings.NotificationLevel,
			IsFavorite:        r.Settings.IsFavorite,
		}
	}

	return ent, nil
}

// UpdateToEntityPageFollower cập nhật dữ liệu từ Request DTO vào một Entity đã tồn tại
func UpdateToEntityPageFollower(r *req.PageFollowerReq, ent *entity.PageFollower) error {
	if r == nil || ent == nil {
		return nil
	}

	if r.PageID != "" {
		pageID, err := primitive.ObjectIDFromHex(r.PageID)
		if err != nil {
			return err
		}
		ent.PageID = pageID
	}

	if r.UserID != "" {
		ent.UserID = r.UserID
	}

	if r.Settings != nil {
		// Nếu Entity chưa có Settings thì khởi tạo mới
		if ent.Settings == nil {
			ent.Settings = &entity.FollowerSettings{}
		}
		ent.Settings.NotificationLevel = r.Settings.NotificationLevel
		ent.Settings.IsFavorite = r.Settings.IsFavorite
	}

	return nil
}
