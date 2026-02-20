package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 1. ToEntity: Chuyển từ Req -> Entity để Insert
func ToEntityAlbum(req *req.AlbumReq) (*entity.Album, error) {
	if req == nil {
		return nil, nil
	}

	albuma := &entity.Album{
		UserID:           req.UserID,
		Title:            req.Title,
		Description:      req.Description,
		Type:             req.Type,
		AssetCount:       req.AssetCount,
		LastAssetAddedAt: req.LastAssetAddedAt,
		CommentCount:     req.CommentCount,
		Reactions: entity.AlbumReactionStats{
			Total: req.Reactions.Total,
			Like:  req.Reactions.Like,
			Love:  req.Reactions.Love,
		},
		Privacy: entity.AlbumPrivacy{
			Level:     req.Privacy.Level,
			AllowList: req.Privacy.AllowList,
			BlockList: req.Privacy.BlockList,
		},
	}

	// Xử lý ID
	if req.ID != "" {
		id, err := primitive.ObjectIDFromHex(req.ID)
		if err == nil {
			albuma.ID = id
		}
	} else {
		albuma.ID = primitive.NewObjectID()
	}

	// Xử lý GroupID
	if req.GroupID != "" {
		groupID, err := primitive.ObjectIDFromHex(req.GroupID)
		if err == nil {
			albuma.GroupID = &groupID
		}
	}

	// Xử lý CoverAssetID
	if req.CoverAssetID != "" {
		coverID, err := primitive.ObjectIDFromHex(req.CoverAssetID)
		if err == nil {
			albuma.CoverAssetID = coverID
		}
	}

	// Xử lý Timestamps
	now := time.Now()
	if req.CreatedAt != nil {
		albuma.CreatedAt = *req.CreatedAt
	} else {
		albuma.CreatedAt = now
	}

	if req.UpdatedAt != nil {
		albuma.UpdatedAt = *req.UpdatedAt
	} else {
		albuma.UpdatedAt = now
	}
	albuma.DeletedAt = req.DeletedAt

	return albuma, nil
}

// 2. UpdateToEntity: Đổ dữ liệu từ UpdateReq vào Entity có sẵn (Partial Update)
func UpdateToEntityAlbum(req *req.UpdateAlbumReq, album *entity.Album) {
	if req == nil || album == nil {
		return
	}

	if req.Title != nil {
		album.Title = *req.Title
	}
	if req.Description != nil {
		album.Description = *req.Description
	}
	if req.Type != nil {
		album.Type = *req.Type
	}
	if req.AssetCount != nil {
		album.AssetCount = *req.AssetCount
	}
	if req.CommentCount != nil {
		album.CommentCount = *req.CommentCount
	}
	if req.LastAssetAddedAt != nil {
		album.LastAssetAddedAt = req.LastAssetAddedAt
	}

	if req.CoverAssetID != nil && *req.CoverAssetID != "" {
		coverID, err := primitive.ObjectIDFromHex(*req.CoverAssetID)
		if err == nil {
			album.CoverAssetID = coverID
		}
	}

	if req.Reactions != nil {
		album.Reactions.Total = req.Reactions.Total
		album.Reactions.Like = req.Reactions.Like
		album.Reactions.Love = req.Reactions.Love
	}

	if req.Privacy != nil {
		album.Privacy.Level = req.Privacy.Level
		album.Privacy.AllowList = req.Privacy.AllowList
		album.Privacy.BlockList = req.Privacy.BlockList
	}

	album.UpdatedAt = time.Now()
}

// 3. EntityToResponse: Chuyển từ Entity -> Res (Chuẩn hóa trả về Client)
func EntityToResponseAlbum(album *entity.Album) *res.AlbumRes {
	if album == nil {
		return nil
	}

	resa := &res.AlbumRes{
		ID:               album.ID.Hex(),
		UserID:           album.UserID,
		Title:            album.Title,
		Description:      album.Description,
		Type:             album.Type,
		AssetCount:       album.AssetCount,
		LastAssetAddedAt: album.LastAssetAddedAt,
		CommentCount:     album.CommentCount,
		Reactions: &res.AlbumReactionStatsRes{
			Total: album.Reactions.Total,
			Like:  album.Reactions.Like,
			Love:  album.Reactions.Love,
		},
		Privacy: &res.AlbumPrivacyRes{
			Level:     album.Privacy.Level,
			AllowList: album.Privacy.AllowList,
			BlockList: album.Privacy.BlockList,
		},
		CreatedAt: album.CreatedAt,
		UpdatedAt: album.UpdatedAt,
		DeletedAt: album.DeletedAt,
	}

	if album.GroupID != nil {
		resa.GroupID = album.GroupID.Hex()
	}

	if !album.CoverAssetID.IsZero() {
		resa.CoverAssetID = album.CoverAssetID.Hex()
	}

	return resa
}

// 4. ReqToResponseAlbum: Hàm mapping trực tiếp từ Req -> Res theo đúng yêu cầu
func ReqToResponseAlbum(req *req.AlbumReq) *res.AlbumRes {
	if req == nil {
		return nil
	}

	resa := &res.AlbumRes{
		ID:               req.ID,
		UserID:           req.UserID,
		GroupID:          req.GroupID,
		Title:            req.Title,
		Description:      req.Description,
		Type:             req.Type,
		CoverAssetID:     req.CoverAssetID,
		AssetCount:       req.AssetCount,
		LastAssetAddedAt: req.LastAssetAddedAt,
		CommentCount:     req.CommentCount,
		Reactions: &res.AlbumReactionStatsRes{
			Total: req.Reactions.Total,
			Like:  req.Reactions.Like,
			Love:  req.Reactions.Love,
		},
		Privacy: &res.AlbumPrivacyRes{
			Level:     req.Privacy.Level,
			AllowList: req.Privacy.AllowList,
			BlockList: req.Privacy.BlockList,
		},
		DeletedAt: req.DeletedAt,
	}

	now := time.Now()
	if req.CreatedAt != nil {
		resa.CreatedAt = *req.CreatedAt
	} else {
		resa.CreatedAt = now
	}

	if req.UpdatedAt != nil {
		resa.UpdatedAt = *req.UpdatedAt
	} else {
		resa.UpdatedAt = now
	}

	return resa
}
