package mapper

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- 1. Request -> Entity (Dùng cho Create) ---
func ToPostSettingEntityPostSetting(r *req.PostSettingReq) *entity.PostSetting {
	if r == nil {
		return nil
	}

	// Xử lý ObjectID chính
	objectID := primitive.NewObjectID()
	if r.ID != "" {
		if oid, err := primitive.ObjectIDFromHex(r.ID); err == nil {
			objectID = oid
		}
	}

	// Xử lý PostID (Bắt buộc)
	postID, _ := primitive.ObjectIDFromHex(r.PostID)

	setting := &entity.PostSetting{
		ID:        objectID,
		PostID:    postID,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		DeletedAt: r.DeletedAt,
	}

	if r.Schedule != nil {
		setting.Schedule = &entity.PostSchedule{
			IsScheduled:        r.Schedule.IsScheduled,
			PublishTime:        r.Schedule.PublishTime,
			PublisherUserID:    r.Schedule.PublisherUserID,
			AuthorRoleSnapshot: r.Schedule.AuthorRoleSnapshot,
		}
	}

	if r.AdsInfo != nil {
		campaignID, _ := primitive.ObjectIDFromHex(r.AdsInfo.CampaignID)
		setting.AdsInfo = &entity.AdsInfo{
			CampaignID:  campaignID,
			IsSponsored: r.AdsInfo.IsSponsored,
			CTALink:     r.AdsInfo.CTALink,
		}
	}

	if r.Targeting != nil {
		setting.Targeting = &entity.PostTargeting{
			Locations: r.Targeting.Locations,
			AgeMin:    r.Targeting.AgeMin,
			AgeMax:    r.Targeting.AgeMax,
			Genders:   r.Targeting.Genders,
			Languages: r.Targeting.Languages,
			Interests: r.Targeting.Interests,
		}
	}

	return setting
}

// --- 2. Update Request -> Existing Entity ---
func UpdatePostSettingEntityPostSetting(r *req.PostSettingReq, setting *entity.PostSetting) {
	if r == nil || setting == nil {
		return
	}

	// Best Practice: Bỏ qua ID, PostID, CreatedAt khi update để đảm bảo an toàn dữ liệu
	setting.UpdatedAt = r.UpdatedAt
	setting.DeletedAt = r.DeletedAt

	if r.Schedule != nil {
		if setting.Schedule == nil {
			setting.Schedule = &entity.PostSchedule{}
		}
		setting.Schedule.IsScheduled = r.Schedule.IsScheduled
		setting.Schedule.PublishTime = r.Schedule.PublishTime
		setting.Schedule.PublisherUserID = r.Schedule.PublisherUserID
		setting.Schedule.AuthorRoleSnapshot = r.Schedule.AuthorRoleSnapshot
	}

	if r.AdsInfo != nil {
		if setting.AdsInfo == nil {
			setting.AdsInfo = &entity.AdsInfo{}
		}
		if cid, err := primitive.ObjectIDFromHex(r.AdsInfo.CampaignID); err == nil {
			setting.AdsInfo.CampaignID = cid
		}
		setting.AdsInfo.IsSponsored = r.AdsInfo.IsSponsored
		setting.AdsInfo.CTALink = r.AdsInfo.CTALink
	}

	if r.Targeting != nil {
		if setting.Targeting == nil {
			setting.Targeting = &entity.PostTargeting{}
		}
		setting.Targeting.Locations = r.Targeting.Locations
		setting.Targeting.AgeMin = r.Targeting.AgeMin
		setting.Targeting.AgeMax = r.Targeting.AgeMax
		setting.Targeting.Genders = r.Targeting.Genders
		setting.Targeting.Languages = r.Targeting.Languages
		setting.Targeting.Interests = r.Targeting.Interests
	}
}

// --- 3. Entity -> Response ---
func ToPostSettingResPostSetting(setting *entity.PostSetting) *res.PostSettingRes {
	if setting == nil {
		return nil
	}

	result := &res.PostSettingRes{
		ID:        setting.ID.Hex(),
		PostID:    setting.PostID.Hex(),
		CreatedAt: setting.CreatedAt,
		UpdatedAt: setting.UpdatedAt,
		DeletedAt: setting.DeletedAt,
	}

	if setting.Schedule != nil {
		result.Schedule = &res.PostScheduleRes{
			IsScheduled:        setting.Schedule.IsScheduled,
			PublishTime:        setting.Schedule.PublishTime,
			PublisherUserID:    setting.Schedule.PublisherUserID,
			AuthorRoleSnapshot: setting.Schedule.AuthorRoleSnapshot,
		}
	}

	if setting.AdsInfo != nil {
		result.AdsInfo = &res.AdsInfoRes{
			CampaignID:  setting.AdsInfo.CampaignID.Hex(),
			IsSponsored: setting.AdsInfo.IsSponsored,
			CTALink:     setting.AdsInfo.CTALink,
		}
	}

	if setting.Targeting != nil {
		result.Targeting = &res.PostTargetingRes{
			Locations: setting.Targeting.Locations,
			AgeMin:    setting.Targeting.AgeMin,
			AgeMax:    setting.Targeting.AgeMax,
			Genders:   setting.Targeting.Genders,
			Languages: setting.Targeting.Languages,
			Interests: setting.Targeting.Interests,
		}
	}

	return result
}

// --- 4. Request -> Response ---
func ReqToPostSettingResPostSetting(r *req.PostSettingReq) *res.PostSettingRes {
	if r == nil {
		return nil
	}

	result := &res.PostSettingRes{
		ID:        r.ID,
		PostID:    r.PostID,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		DeletedAt: r.DeletedAt,
	}

	if r.Schedule != nil {
		result.Schedule = &res.PostScheduleRes{
			IsScheduled:        r.Schedule.IsScheduled,
			PublishTime:        r.Schedule.PublishTime,
			PublisherUserID:    r.Schedule.PublisherUserID,
			AuthorRoleSnapshot: r.Schedule.AuthorRoleSnapshot,
		}
	}

	if r.AdsInfo != nil {
		result.AdsInfo = &res.AdsInfoRes{
			CampaignID:  r.AdsInfo.CampaignID,
			IsSponsored: r.AdsInfo.IsSponsored,
			CTALink:     r.AdsInfo.CTALink,
		}
	}

	if r.Targeting != nil {
		result.Targeting = &res.PostTargetingRes{
			Locations: r.Targeting.Locations,
			AgeMin:    r.Targeting.AgeMin,
			AgeMax:    r.Targeting.AgeMax,
			Genders:   r.Targeting.Genders,
			Languages: r.Targeting.Languages,
			Interests: r.Targeting.Interests,
		}
	}

	return result
}
