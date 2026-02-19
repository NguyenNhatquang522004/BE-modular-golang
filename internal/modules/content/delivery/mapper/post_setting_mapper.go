package mapper

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- 1. REQ TO ENTITY (CREATE) ---

func ToPostSettingEntity(r req.CreatePostSettingReq) (*entity.PostSetting, error) {
	// Convert PostID from string to ObjectID
	postObjID, err := primitive.ObjectIDFromHex(r.PostID)
	if err != nil {
		return nil, err
	}

	ent := &entity.PostSetting{
		ID:     primitive.NewObjectID(), // Tự sinh ID mới
		PostID: postObjID,
	}

	// Map Schedule
	if r.Schedule != nil {
		ent.Schedule = &entity.PostSchedule{
			IsScheduled:        r.Schedule.IsScheduled,
			PublishTime:        r.Schedule.PublishTime,
			PublisherUserID:    r.Schedule.PublisherUserID,
			AuthorRoleSnapshot: r.Schedule.AuthorRoleSnapshot,
		}
	}

	// Map AdsInfo
	if r.AdsInfo != nil {
		campaignObjID, err := primitive.ObjectIDFromHex(r.AdsInfo.CampaignID)
		if err == nil { // Nếu ID lỗi thì có thể bỏ qua hoặc return err tùy logic strict
			ent.AdsInfo = &entity.AdsInfo{
				CampaignID:  campaignObjID,
				IsSponsored: r.AdsInfo.IsSponsored,
				CTALink:     r.AdsInfo.CTALink,
			}
		}
	}

	// Map Targeting
	if r.Targeting != nil {
		ent.Targeting = &entity.PostTargeting{
			Locations: r.Targeting.Locations,
			AgeMin:    r.Targeting.AgeMin,
			AgeMax:    r.Targeting.AgeMax,
			Genders:   r.Targeting.Genders,
			Languages: r.Targeting.Languages,
			Interests: r.Targeting.Interests,
		}
	}

	return ent, nil
}

// --- 2. REQ TO ENTITY (UPDATE) ---

func UpdatePostSettingEntity(existingEnt *entity.PostSetting, r req.UpdatePostSettingReq) {
	// Chỉ update nếu request gửi dữ liệu (khác nil)

	// 1. Update Schedule
	if r.Schedule != nil {
		// Nếu trong DB chưa có (nil) thì khởi tạo mới
		if existingEnt.Schedule == nil {
			existingEnt.Schedule = &entity.PostSchedule{}
		}
		// Gán giá trị
		existingEnt.Schedule.IsScheduled = r.Schedule.IsScheduled
		existingEnt.Schedule.PublishTime = r.Schedule.PublishTime
		existingEnt.Schedule.PublisherUserID = r.Schedule.PublisherUserID
		existingEnt.Schedule.AuthorRoleSnapshot = r.Schedule.AuthorRoleSnapshot
	}

	// 2. Update AdsInfo
	if r.AdsInfo != nil {
		campaignObjID, err := primitive.ObjectIDFromHex(r.AdsInfo.CampaignID)
		if err == nil { // Chỉ update nếu ID hợp lệ
			if existingEnt.AdsInfo == nil {
				existingEnt.AdsInfo = &entity.AdsInfo{}
			}
			existingEnt.AdsInfo.CampaignID = campaignObjID
			existingEnt.AdsInfo.IsSponsored = r.AdsInfo.IsSponsored
			existingEnt.AdsInfo.CTALink = r.AdsInfo.CTALink
		}
	}

	// 3. Update Targeting
	if r.Targeting != nil {
		if existingEnt.Targeting == nil {
			existingEnt.Targeting = &entity.PostTargeting{}
		}
		existingEnt.Targeting.Locations = r.Targeting.Locations
		existingEnt.Targeting.AgeMin = r.Targeting.AgeMin
		existingEnt.Targeting.AgeMax = r.Targeting.AgeMax
		existingEnt.Targeting.Genders = r.Targeting.Genders
		existingEnt.Targeting.Languages = r.Targeting.Languages
		existingEnt.Targeting.Interests = r.Targeting.Interests
	}
}

// --- 3. ENTITY TO RES ---

func ToPostSettingRes(ent *entity.PostSetting) *res.PostSettingRes {
	if ent == nil {
		return nil
	}

	response := &res.PostSettingRes{
		ID:     ent.ID.Hex(),
		PostID: ent.PostID.Hex(),
	}

	// Map Schedule Res
	if ent.Schedule != nil {
		response.Schedule = &res.PostScheduleRes{
			IsScheduled:        ent.Schedule.IsScheduled,
			PublishTime:        ent.Schedule.PublishTime,
			PublisherUserID:    ent.Schedule.PublisherUserID,
			AuthorRoleSnapshot: ent.Schedule.AuthorRoleSnapshot,
		}
	}

	// Map AdsInfo Res
	if ent.AdsInfo != nil {
		response.AdsInfo = &res.AdsInfoRes{
			CampaignID:  ent.AdsInfo.CampaignID.Hex(),
			IsSponsored: ent.AdsInfo.IsSponsored,
			CTALink:     ent.AdsInfo.CTALink,
		}
	}

	// Map Targeting Res
	if ent.Targeting != nil {
		response.Targeting = &res.PostTargetingRes{
			Locations: ent.Targeting.Locations,
			AgeMin:    ent.Targeting.AgeMin,
			AgeMax:    ent.Targeting.AgeMax,
			Genders:   ent.Targeting.Genders,
			Languages: ent.Targeting.Languages,
			Interests: ent.Targeting.Interests,
		}
	}

	return response
}
