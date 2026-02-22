package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// -----------------------------------------------------------------------------
// REQ MAPPER
// -----------------------------------------------------------------------------

// ToEntityPages chuyển từ PageReq sang Entity (Thường dùng cho Create)
func ToEntityPages(r *req.PageReq) *entity.Page {
	if r == nil {
		return nil
	}

	objID := primitive.NewObjectID()
	if r.ID != "" {
		if parsedID, err := primitive.ObjectIDFromHex(r.ID); err == nil {
			objID = parsedID
		}
	}

	e := &entity.Page{
		ID:            objID,
		CreatorUserID: r.CreatorUserID,
		Name:          r.Name,
		Slug:          r.Slug,
		CategoryID:    r.CategoryID,
		IsVerified:    r.IsVerified,
		Status:        r.Status,
		Avatar:        &entity.PageAvatar{URL: r.Avatar.URL},
		Cover:         &entity.PageCover{URL: r.Cover.URL, PositionY: r.Cover.PositionY},
		Bio:           r.Bio,
		Website:       r.Website,
		Email:         r.Email,
		PhoneNumber:   r.PhoneNumber,
		Address: &entity.PageAddress{
			Street:      r.Address.Street,
			City:        r.Address.City,
			Zipcode:     r.Address.Zipcode,
			Coordinates: copyFloatSlicePages(r.Address.Coordinates),
		},
		Settings: &entity.PageSettings{
			AllowVisitorPost: r.Settings.AllowVisitorPost,
			ProfanityFilter:  r.Settings.ProfanityFilter,
			MessagingStatus:  r.Settings.MessagingStatus,
		},
		Stats: &entity.PageStats{
			FollowersCount: r.Stats.FollowersCount,
			LikesCount:     r.Stats.LikesCount,
			RatingScore:    r.Stats.RatingScore,
			ReviewCount:    r.Stats.ReviewCount,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if r.CTAButton != nil {
		e.CTAButton = &entity.CTAButton{
			Type:  r.CTAButton.Type,
			Value: r.CTAButton.Value,
		}
	}

	for _, bh := range r.BusinessHours {
		e.BusinessHours = append(e.BusinessHours, entity.BusinessHour{
			Day:   bh.Day,
			Open:  bh.Open,
			Close: bh.Close,
		})
	}

	return e
}

// UpdateToEntityPages cập nhật các trường từ PageReq vào Entity có sẵn (Thường dùng cho Update)
func UpdateToEntityPages(r *req.PageReq, e *entity.Page) {
	if r == nil || e == nil {
		return
	}

	// Lưu ý: ID và CreatorUserID thường là immutable, nhưng map 100% theo yêu cầu
	if r.ID != "" {
		if parsedID, err := primitive.ObjectIDFromHex(r.ID); err == nil {
			e.ID = parsedID
		}
	}
	e.CreatorUserID = r.CreatorUserID
	e.Name = r.Name
	e.Slug = r.Slug
	e.CategoryID = r.CategoryID
	e.IsVerified = r.IsVerified
	e.Status = r.Status

	e.Avatar = &entity.PageAvatar{URL: r.Avatar.URL}
	e.Cover = &entity.PageCover{URL: r.Cover.URL, PositionY: r.Cover.PositionY}

	e.Bio = r.Bio
	e.Website = r.Website
	e.Email = r.Email
	e.PhoneNumber = r.PhoneNumber

	e.Address = &entity.PageAddress{
		Street:      r.Address.Street,
		City:        r.Address.City,
		Zipcode:     r.Address.Zipcode,
		Coordinates: copyFloatSlicePages(r.Address.Coordinates),
	}

	// Ghi đè BusinessHours
	e.BusinessHours = make([]entity.BusinessHour, 0, len(r.BusinessHours))
	for _, bh := range r.BusinessHours {
		e.BusinessHours = append(e.BusinessHours, entity.BusinessHour{
			Day:   bh.Day,
			Open:  bh.Open,
			Close: bh.Close,
		})
	}

	if r.CTAButton != nil {
		e.CTAButton = &entity.CTAButton{
			Type:  r.CTAButton.Type,
			Value: r.CTAButton.Value,
		}
	} else {
		e.CTAButton = nil
	}

	e.Settings = &entity.PageSettings{
		AllowVisitorPost: r.Settings.AllowVisitorPost,
		ProfanityFilter:  r.Settings.ProfanityFilter,
		MessagingStatus:  r.Settings.MessagingStatus,
	}

	e.Stats = &entity.PageStats{
		FollowersCount: r.Stats.FollowersCount,
		LikesCount:     r.Stats.LikesCount,
		RatingScore:    r.Stats.RatingScore,
		ReviewCount:    r.Stats.ReviewCount,
	}

	// Update timestamp
	e.UpdatedAt = time.Now()
}

// -----------------------------------------------------------------------------
// RES MAPPER
// -----------------------------------------------------------------------------

// ToResponsePages chuyển từ Entity sang PageRes
func ToResponsePages(e *entity.Page) *res.PageRes {
	if e == nil {
		return nil
	}

	resp := &res.PageRes{
		ID:            e.ID.Hex(), // Parse primitive.ObjectID -> string
		CreatorUserID: e.CreatorUserID,
		Name:          e.Name,
		Slug:          e.Slug,
		CategoryID:    e.CategoryID,
		IsVerified:    e.IsVerified,
		Status:        e.Status,
		Avatar:        &res.PageAvatarRes{URL: e.Avatar.URL},
		Cover:         &res.PageCoverRes{URL: e.Cover.URL, PositionY: e.Cover.PositionY},
		Bio:           e.Bio,
		Website:       e.Website,
		Email:         e.Email,
		PhoneNumber:   e.PhoneNumber,
		Address: &res.PageAddressRes{
			Street:      e.Address.Street,
			City:        e.Address.City,
			Zipcode:     e.Address.Zipcode,
			Coordinates: copyFloatSlicePages(e.Address.Coordinates),
		},
		Settings: &res.PageSettingsRes{
			AllowVisitorPost: e.Settings.AllowVisitorPost,
			ProfanityFilter:  e.Settings.ProfanityFilter,
			MessagingStatus:  e.Settings.MessagingStatus,
		},
		Stats: &res.PageStatsRes{
			FollowersCount: e.Stats.FollowersCount,
			LikesCount:     e.Stats.LikesCount,
			RatingScore:    e.Stats.RatingScore,
			ReviewCount:    e.Stats.ReviewCount,
		},
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}

	if e.CTAButton != nil {
		resp.CTAButton = &res.CTAButtonRes{
			Type:  e.CTAButton.Type,
			Value: e.CTAButton.Value,
		}
	}

	for _, bh := range e.BusinessHours {
		resp.BusinessHours = append(resp.BusinessHours, res.BusinessHourRes{
			Day:   bh.Day,
			Open:  bh.Open,
			Close: bh.Close,
		})
	}

	return resp
}

// ReqToResponsePages chuyển trực tiếp từ Req sang Res (Theo yêu cầu: "mapper Res có hàm chuyển từ req sang res")
func ReqToResponsePages(r *req.PageReq) *res.PageRes {
	if r == nil {
		return nil
	}

	resp := &res.PageRes{
		ID:            r.ID,
		CreatorUserID: r.CreatorUserID,
		Name:          r.Name,
		Slug:          r.Slug,
		CategoryID:    r.CategoryID,
		IsVerified:    r.IsVerified,
		Status:        r.Status,
		Avatar:        &res.PageAvatarRes{URL: r.Avatar.URL},
		Cover:         &res.PageCoverRes{URL: r.Cover.URL, PositionY: r.Cover.PositionY},
		Bio:           r.Bio,
		Website:       r.Website,
		Email:         r.Email,
		PhoneNumber:   r.PhoneNumber,
		Address: &res.PageAddressRes{
			Street:      r.Address.Street,
			City:        r.Address.City,
			Zipcode:     r.Address.Zipcode,
			Coordinates: copyFloatSlicePages(r.Address.Coordinates),
		},
		Settings: &res.PageSettingsRes{
			AllowVisitorPost: r.Settings.AllowVisitorPost,
			ProfanityFilter:  r.Settings.ProfanityFilter,
			MessagingStatus:  r.Settings.MessagingStatus,
		},
		Stats: &res.PageStatsRes{
			FollowersCount: r.Stats.FollowersCount,
			LikesCount:     r.Stats.LikesCount,
			RatingScore:    r.Stats.RatingScore,
			ReviewCount:    r.Stats.ReviewCount,
		},
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}

	if r.CTAButton != nil {
		resp.CTAButton = &res.CTAButtonRes{
			Type:  r.CTAButton.Type,
			Value: r.CTAButton.Value,
		}
	}

	for _, bh := range r.BusinessHours {
		resp.BusinessHours = append(resp.BusinessHours, res.BusinessHourRes{
			Day:   bh.Day,
			Open:  bh.Open,
			Close: bh.Close,
		})
	}

	return resp
}

// Hàm helper để deep copy slice (Tránh memory leak / tham chiếu không mong muốn)
func copyFloatSlicePages(src []float64) []float64 {
	if src == nil {
		return nil
	}
	dst := make([]float64, len(src))
	copy(dst, src)
	return dst
}
