package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/businessEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MapCreatePayloadToEntity(p *businessEvent.CreatePagePayload) *entity.Page {
	now := time.Now()

	// Khởi tạo ObjectID (nếu PageID được gửi lên thì parse, không thì tạo mới)
	var objID primitive.ObjectID
	if p.PageID != "" {
		id, err := primitive.ObjectIDFromHex(p.PageID)
		if err != nil {
			return nil // Hoặc handle lỗi theo cách bạn muốn
		}
		objID = id
	} else {
		return nil // Hoặc handle lỗi theo cách bạn muốn
	}

	page := &entity.Page{
		ID:            objID,
		CreatorUserID: p.UserID,
		Name:          p.Name,
		Slug:          p.Slug,
		CategoryID:    p.CategoryID,
		IsVerified:    false, // Default khi mới tạo
		// Status:        enum.PageStatusPending, // Bạn có thể set default status ở đây
		Bio:           p.Bio,
		Website:       p.Website,
		Email:         p.Email,
		PhoneNumber:   p.PhoneNumber,
		Avatar:        mapAvatar(p.Avatar),
		Cover:         mapCover(p.Cover),
		Address:       mapAddress(p.Address),
		BusinessHours: mapBusinessHours(p.BusinessHours),
		CTAButton:     mapCTA(p.CTAButton),
		Settings:      mapSettings(p.Settings),
		Stats: &entity.PageStats{ // Khởi tạo Stats mặc định là 0
			FollowersCount: 0,
			LikesCount:     0,
			RatingScore:    0.0,
			ReviewCount:    0,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	return page
}
func ApplyUpdatePayloadToEntity(page *entity.Page, p *businessEvent.UpdatePagePayload) {
	page.UpdatedAt = time.Now()

	// Primitive fields: Check nil trước khi gán
	if p.Name != nil {
		page.Name = *p.Name
	}
	if p.Slug != nil {
		page.Slug = *p.Slug
	}
	if p.CategoryID != nil {
		page.CategoryID = *p.CategoryID
	}
	if p.Bio != nil {
		page.Bio = *p.Bio
	}
	if p.Website != nil {
		page.Website = *p.Website
	}
	if p.Email != nil {
		page.Email = *p.Email
	}
	if p.PhoneNumber != nil {
		page.PhoneNumber = *p.PhoneNumber
	}

	// Nested fields: Ghi đè object nếu client có gửi lên
	if p.Avatar != nil {
		page.Avatar = mapAvatar(p.Avatar)
	}
	if p.Cover != nil {
		page.Cover = mapCover(p.Cover)
	}
	if p.Address != nil {
		page.Address = mapAddress(p.Address)
	}
	if p.CTAButton != nil {
		page.CTAButton = mapCTA(p.CTAButton)
	}
	if p.Settings != nil {
		page.Settings = mapSettings(p.Settings)
	}

	// Array (Slice): Trong Go, slice nil có nghĩa là không gửi, khác với slice rỗng []
	// Nếu != nil, ta sẽ ghi đè toàn bộ mảng cũ bằng mảng mới.
	if p.BusinessHours != nil {
		page.BusinessHours = mapBusinessHours(p.BusinessHours)
	}
}

// ==========================================
// 3. HELPER MAPPERS (Cho cả Create & Update)
// ==========================================

func mapAvatar(p *businessEvent.AvatarPayload) *entity.PageAvatar {
	if p == nil {
		return nil
	}
	return &entity.PageAvatar{
		ID:  primitive.NewObjectID(),
		URL: p.URL,
	}
}

func mapCover(p *businessEvent.CoverPayload) *entity.PageCover {
	if p == nil {
		return nil
	}
	return &entity.PageCover{
		ID:        primitive.NewObjectID(),
		URL:       p.URL,
		PositionY: p.PositionY,
	}
}

func mapAddress(p *businessEvent.AddressPayload) *entity.PageAddress {
	if p == nil {
		return nil
	}
	return &entity.PageAddress{
		Street:      p.Street,
		City:        p.City,
		Zipcode:     p.Zipcode,
		Coordinates: p.Coordinates,
	}
}

func mapBusinessHours(p []businessEvent.BusinessHourPayload) []entity.BusinessHour {
	if p == nil {
		return nil // Hoặc trả về mảng rỗng []entity.BusinessHour{} tùy logic của bạn
	}
	hours := make([]entity.BusinessHour, 0, len(p))
	for _, h := range p {
		hours = append(hours, entity.BusinessHour{
			Day:   h.Day,
			Open:  h.Open,
			Close: h.Close,
		})
	}
	return hours
}

func mapCTA(p *businessEvent.CTAPayload) *entity.CTAButton {
	if p == nil {
		return nil
	}
	return &entity.CTAButton{
		Type:  p.Type,
		Value: p.Value,
	}
}

func mapSettings(p *businessEvent.SettingsPayload) *entity.PageSettings {
	if p == nil {
		return nil
	}

	allowVisitorPost := false
	if p.AllowVisitorPost != nil {
		allowVisitorPost = *p.AllowVisitorPost
	}

	return &entity.PageSettings{
		AllowVisitorPost: allowVisitorPost,
		ProfanityFilter:  p.ProfanityFilter,
		MessagingStatus:  p.MessagingStatus,
	}
}
