package mapper

import (
	"fmt"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToAlbumEntity(req *mediaEvent.CreateAlbumPayload) (*entity.Album, error) {
	now := time.Now()

	// Khởi tạo Entity với các trường bắt buộc và giá trị mặc định hệ thống
	album := &entity.Album{
		ID:          primitive.NewObjectID(), // Tự sinh Mongo ID mới
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		Privacy: entity.AlbumPrivacy{
			Level:     req.Privacy.Level,
			AllowList: req.Privacy.AllowList,
			BlockList: req.Privacy.BlockList,
		},
		AssetCount: 0, // Mặc định album mới chưa có ảnh
		CreatedAt:  now,
		UpdatedAt:  now,
		// Reactions và CommentCount mặc định bằng 0 do Zero-value của Go
	}

	// Xử lý an toàn GroupID (nếu client có truyền lên)
	if req.GroupID != nil && *req.GroupID != "" {
		groupID, err := primitive.ObjectIDFromHex(*req.GroupID)
		if err != nil {
			// Trả về lỗi rõ ràng để tầng Service/Controller handle
			return nil, fmt.Errorf("invalid group_id format: %w", err)
		}
		album.GroupID = &groupID
	}

	return album, nil
}
func ApplyAlbumUpdate(existing entity.Album, req *mediaEvent.UpdateAlbumPayload) (*entity.Album, error) {
	// Kiểm tra từng con trỏ trong DTO. Nếu khác nil nghĩa là client muốn cập nhật.
	if req.Title != nil {
		existing.Title = *req.Title
	}

	if req.Description != nil {
		existing.Description = *req.Description
	}

	if req.Type != nil {
		existing.Type = *req.Type
	}

	if req.Privacy != nil {
		existing.Privacy.Level = req.Privacy.Level
		// Nạp đè toàn bộ mảng (thực tế tuỳ nghiệp vụ có thể merge, nhưng update nguyên mảng là phổ biến nhất)
		existing.Privacy.AllowList = req.Privacy.AllowList
		existing.Privacy.BlockList = req.Privacy.BlockList
	}

	if req.CoverAssetID != nil {
		// Dùng để handle trường hợp update ObjectID
		if *req.CoverAssetID != "" {
			coverID, err := primitive.ObjectIDFromHex(*req.CoverAssetID)
			if err != nil {
				return nil, fmt.Errorf("invalid cover_asset_id format: %w", err)
			}
			existing.CoverAssetID = coverID
		} else {
			// Tuỳ nghiệp vụ: Nếu client gửi string rỗng `""`, có thể hiểu là muốn xoá Cover.
			// Ở Mongo, ta set nó về chuỗi ObjectID Zero (hoặc đổi kiểu ở DB thành pointer nếu muốn lưu null)
			existing.CoverAssetID = primitive.NilObjectID
		}
	}

	// Cập nhật lại thời gian Update (Bắt buộc trong best practice)
	existing.UpdatedAt = time.Now()

	// Trả về pointer trỏ tới bản copy đã được update thành công
	return &existing, nil
}
