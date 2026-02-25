package mapper

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ToEntityUserNotificationSetting: Ánh xạ từ Create Req sang Entity mới
func ToEntityUserNotificationSetting(r *req.CreateUserNotificationSettingReq) *entity.UserNotificationSetting {
	if r == nil {
		return nil
	}

	return &entity.UserNotificationSetting{
		ID:          primitive.NewObjectID(), // Tự động generate ObjectID cho bản ghi mới
		Name:        r.Name,
		UserID:      r.UserID,
		Avatar:      r.Avatar,
		DateOfBirth: r.DateOfBirth,
		Settings:    mapGeneralSettingsReqToEntity(r.Settings),
		FCMTokens:   mapFCMTokensReqToEntity(r.FCMTokens),
	}
}

// UpdateToEntityUserNotificationSetting: Cập nhật từ Update Req vào Entity hiện có (Partial Update)
func UpdateToEntityUserNotificationSetting(r *req.UpdateUserNotificationSettingReq, ent *entity.UserNotificationSetting) {
	if r == nil || ent == nil {
		return
	}

	if r.Settings != nil {
		ent.Settings = mapGeneralSettingsReqToEntity(r.Settings)
	}

	if r.FCMTokens != nil {
		ent.FCMTokens = mapFCMTokensReqToEntity(r.FCMTokens)
	}

	if r.DateOfBirth != "" {
		ent.DateOfBirth = r.DateOfBirth
	}
	if r.Name != "" {
		ent.Name = r.Name
	}
	if r.Avatar != "" {
		ent.Avatar = r.Avatar
	}
}

// EntityToResUserNotificationSetting: Ánh xạ từ Entity sang Response DTO
func EntityToResUserNotificationSetting(ent *entity.UserNotificationSetting) *res.UserNotificationSettingRes {
	if ent == nil {
		return nil
	}

	return &res.UserNotificationSettingRes{
		ID:          ent.ID.Hex(), // Convert ObjectID sang dạng chuỗi Hex
		UserID:      ent.UserID,
		Name:        ent.Name,
		Avatar:      ent.Avatar,
		DateOfBirth: ent.DateOfBirth,
		Settings:    mapGeneralSettingsEntityToRes(ent.Settings),
		FCMTokens:   mapFCMTokensEntityToRes(ent.FCMTokens),
	}
}

// ================= INTERNAL HELPERS =================

func mapGeneralSettingsReqToEntity(s *req.GeneralSettingsReq) *entity.GeneralSettings {
	if s == nil {
		return nil
	}
	return &entity.GeneralSettings{
		PushEnabled:      true,
		EmailFrequency:   enum.EmailFreqInstant,
		PushInteractions: true,
		PushFriends:      true,
		PushGroups:       true,
		PushEvents:       true,
		PushBirthdays:    true,
	}
}

func mapGeneralSettingsEntityToRes(s *entity.GeneralSettings) *res.GeneralSettingsRes {
	if s == nil {
		return nil
	}
	return &res.GeneralSettingsRes{
		PushEnabled:      s.PushEnabled,
		EmailFrequency:   &s.EmailFrequency,
		PushInteractions: s.PushInteractions,
		PushFriends:      s.PushFriends,
		PushGroups:       s.PushGroups,
		PushEvents:       s.PushEvents,
		PushBirthdays:    s.PushBirthdays,
	}
}

func mapFCMTokensReqToEntity(tokens []*req.FCMTokenReq) []*entity.FCMToken {
	if tokens == nil {
		return nil
	}
	// Best practice: Cấp phát trước bộ nhớ (pre-allocate) cho slice
	resTokens := make([]*entity.FCMToken, len(tokens))
	for i, t := range tokens {
		resTokens[i] = &entity.FCMToken{
			Token:     t.Token,
			DeviceID:  t.DeviceID,
			UpdatedAt: t.UpdatedAt,
		}
	}
	return resTokens
}

func mapFCMTokensEntityToRes(tokens []*entity.FCMToken) []*res.FCMTokenRes {
	if tokens == nil {
		return nil
	}
	resTokens := make([]*res.FCMTokenRes, len(tokens))
	for i, t := range tokens {
		resTokens[i] = &res.FCMTokenRes{
			Token:     t.Token,
			DeviceID:  t.DeviceID,
			UpdatedAt: t.UpdatedAt,
		}
	}
	return resTokens
}
