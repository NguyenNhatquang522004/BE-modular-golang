package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToArtistEntity(req *mediaEvent.CreateArtistPayload) *entity.Artist {
	now := time.Now().UTC() // Luôn dùng UTC để lưu trữ vào Database
	if req.ArtistID == "" {
		return nil // Trường ArtistID là bắt buộc, nếu thiếu thì trả về nil để đánh dấu lỗi
	}
	convertedArtistID, err := primitive.ObjectIDFromHex(req.ArtistID)
	if err != nil {
		return nil // Nếu ArtistID không phải là ObjectID hợp lệ, trả về nil để đánh dấu lỗi
	}
	artist := &entity.Artist{
		ID:            convertedArtistID, // Auto-generate MongoDB ID
		Name:          req.Name,
		Slug:          req.Slug,
		Bio:           req.Bio,
		AvatarURL:     req.AvatarURL,
		CoverURL:      req.CoverURL,
		UserID:        req.UserID,
		IsVerified:    false, // Default: Mới tạo thì chưa được verify
		FollowerCount: 0,     // Default: Mới tạo chưa có follower
		TotalStreams:  0,     // Default: Mới tạo chưa có lượt nghe
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Xử lý SocialLinks nếu client có gửi lên
	if req.SocialLinks != nil {
		artist.SocialLinks = entity.SocialLinks{
			Spotify:   req.SocialLinks.Spotify,
			Youtube:   req.SocialLinks.Youtube,
			Instagram: req.SocialLinks.Instagram,
			Facebook:  req.SocialLinks.Facebook,
			Website:   req.SocialLinks.Website,
		}
	}

	return artist
}

func ApplyUpdateToArtist(existingArtist *entity.Artist, req mediaEvent.UpdateArtistPayload) *entity.Artist {
	// Luôn cập nhật thời gian Update
	existingArtist.UpdatedAt = time.Now().UTC()

	// Partial Update: Kiểm tra con trỏ khác nil trước khi gán dữ liệu (dereference)
	if req.Name != nil {
		existingArtist.Name = *req.Name
	}

	if req.Slug != nil {
		existingArtist.Slug = *req.Slug
	}

	if req.Bio != nil {
		existingArtist.Bio = *req.Bio
	}

	if req.AvatarURL != nil {
		existingArtist.AvatarURL = *req.AvatarURL
	}

	if req.CoverURL != nil {
		existingArtist.CoverURL = *req.CoverURL
	}

	// Xử lý ghi đè object SocialLinks
	if req.SocialLinks != nil {
		existingArtist.SocialLinks = entity.SocialLinks{
			Spotify:   req.SocialLinks.Spotify,
			Youtube:   req.SocialLinks.Youtube,
			Instagram: req.SocialLinks.Instagram,
			Facebook:  req.SocialLinks.Facebook,
			Website:   req.SocialLinks.Website,
		}
	}

	return existingArtist
}
