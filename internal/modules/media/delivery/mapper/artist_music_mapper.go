package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ToEntityArtist chuyển đổi từ Req DTO sang Entity để insert DB (Tạo mới)
func ToEntityArtist(request *req.ArtistReq) *entity.Artist {
	if request == nil {
		return nil
	}

	// Xử lý ID an toàn: Nếu client truyền lên HexID hợp lệ thì dùng, không thì Gen mới
	var objID primitive.ObjectID
	if request.ID != "" {
		id, err := primitive.ObjectIDFromHex(request.ID)
		if err == nil {
			objID = id
		}
	}
	if objID.IsZero() {
		objID = primitive.NewObjectID()
	}

	return &entity.Artist{
		ID:         objID,
		Name:       request.Name,
		Slug:       request.Slug,
		Bio:        request.Bio,
		AvatarURL:  request.AvatarURL,
		CoverURL:   request.CoverURL,
		IsVerified: request.IsVerified,
		UserID:     request.UserID,
		SocialLinks: entity.SocialLinks{
			Spotify:   request.SocialLinks.Spotify,
			Youtube:   request.SocialLinks.Youtube,
			Instagram: request.SocialLinks.Instagram,
			Facebook:  request.SocialLinks.Facebook,
			Website:   request.SocialLinks.Website,
		},
		FollowerCount: request.FollowerCount,
		TotalStreams:  request.TotalStreams,
		CreatedAt:     request.CreatedAt,
		UpdatedAt:     request.UpdatedAt,
		DeletedAt:     request.DeletedAt,
	}
}

// UpdateToEntityArtist cập nhật dữ liệu từ Req DTO ghi đè vào Entity hiện có (Để update DB)
func UpdateToEntityArtist(request *req.ArtistReq, existingEntity *entity.Artist) {
	if request == nil || existingEntity == nil {
		return
	}

	existingEntity.Name = request.Name
	existingEntity.Slug = request.Slug
	existingEntity.Bio = request.Bio
	existingEntity.AvatarURL = request.AvatarURL
	existingEntity.CoverURL = request.CoverURL
	existingEntity.IsVerified = request.IsVerified
	existingEntity.UserID = request.UserID

	existingEntity.SocialLinks.Spotify = request.SocialLinks.Spotify
	existingEntity.SocialLinks.Youtube = request.SocialLinks.Youtube
	existingEntity.SocialLinks.Instagram = request.SocialLinks.Instagram
	existingEntity.SocialLinks.Facebook = request.SocialLinks.Facebook
	existingEntity.SocialLinks.Website = request.SocialLinks.Website

	existingEntity.FollowerCount = request.FollowerCount
	existingEntity.TotalStreams = request.TotalStreams

	// Xử lý timestamps
	if !request.UpdatedAt.IsZero() {
		existingEntity.UpdatedAt = request.UpdatedAt
	} else {
		existingEntity.UpdatedAt = time.Now()
	}

	if request.DeletedAt != nil {
		existingEntity.DeletedAt = request.DeletedAt
	}
}

// ReqToResArtist chuyển đổi thẳng từ Request sang Response (Theo đúng yêu cầu của bạn)
func ReqToResArtist(request *req.ArtistReq) *res.ArtistRes {
	if request == nil {
		return nil
	}

	return &res.ArtistRes{
		ID:         request.ID,
		Name:       request.Name,
		Slug:       request.Slug,
		Bio:        request.Bio,
		AvatarURL:  request.AvatarURL,
		CoverURL:   request.CoverURL,
		IsVerified: request.IsVerified,
		UserID:     request.UserID,
		SocialLinks: res.SocialLinksRes{
			Spotify:   request.SocialLinks.Spotify,
			Youtube:   request.SocialLinks.Youtube,
			Instagram: request.SocialLinks.Instagram,
			Facebook:  request.SocialLinks.Facebook,
			Website:   request.SocialLinks.Website,
		},
		FollowerCount: request.FollowerCount,
		TotalStreams:  request.TotalStreams,
		CreatedAt:     request.CreatedAt,
		UpdatedAt:     request.UpdatedAt,
		DeletedAt:     request.DeletedAt,
	}
}

// EntityToResArtist chuyển đổi từ DB Entity sang Res DTO để trả về Client (Best Practice)
func EntityToResArtist(ent *entity.Artist) *res.ArtistRes {
	if ent == nil {
		return nil
	}

	return &res.ArtistRes{
		ID:         ent.ID.Hex(), // Chuyển từ ObjectID sang chuẩn Hex String
		Name:       ent.Name,
		Slug:       ent.Slug,
		Bio:        ent.Bio,
		AvatarURL:  ent.AvatarURL,
		CoverURL:   ent.CoverURL,
		IsVerified: ent.IsVerified,
		UserID:     ent.UserID,
		SocialLinks: res.SocialLinksRes{
			Spotify:   ent.SocialLinks.Spotify,
			Youtube:   ent.SocialLinks.Youtube,
			Instagram: ent.SocialLinks.Instagram,
			Facebook:  ent.SocialLinks.Facebook,
			Website:   ent.SocialLinks.Website,
		},
		FollowerCount: ent.FollowerCount,
		TotalStreams:  ent.TotalStreams,
		CreatedAt:     ent.CreatedAt,
		UpdatedAt:     ent.UpdatedAt,
		DeletedAt:     ent.DeletedAt,
	}
}
