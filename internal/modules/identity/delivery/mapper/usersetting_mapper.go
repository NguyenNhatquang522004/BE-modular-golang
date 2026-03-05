package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
)

func ToEntityUserSetting(req *req.UserSettingReq) *entity.UserSetting {
	e := &entity.UserSetting{

		// 1. Identity
		User_ID: req.User_ID,

		Theme_Mode:     req.Theme_Mode,
		Font_Size:      req.Font_Size,
		Compact_Mode:   req.Compact_Mode,
		Lang_Code:      req.Lang_Code,
		Timezone:       req.Timezone,
		Auto_Translate: req.Auto_Translate,

		// --- 3. PRIVACY DEFAULTS ---
		Default_Post_Audience:  req.Default_Post_Audience,
		Default_Story_Audience: req.Default_Story_Audience,

		// --- 4. ACCESS CONTROL ---
		Allow_Friend_Request_From:    req.Allow_Friend_Request_From,
		Allow_Friend_List_View_From:  req.Allow_Friend_List_View_From,
		Allow_Search_Engine_Indexing: req.Allow_Search_Engine_Indexing,

		// --- 5. TIMELINE & TAGGING ---
		Allow_Timeline_Posting_From:   req.Allow_Timeline_Posting_From,
		Review_Tags_Enabled:           req.Review_Tags_Enabled,
		Review_Timeline_Posts_Enabled: req.Review_Timeline_Posts_Enabled,

		// --- 6. METADATA ---
		Updated_At: time.Now(), // Luôn lấy giờ hiện tại của server
	}

	// --- 7. XỬ LÝ NESTED STRUCT (NOTIFICATIONS) ---
	// Kiểm tra nil để tránh panic hoặc ghi đè dữ liệu rỗng
	e.Notifications = &entity.NotificationSettings{
		EmailFrequency:   req.Notifications.EmailFrequency,
		PushInteractions: req.Notifications.PushInteractions,
		PushFriends:      req.Notifications.PushFriends,
		PushGroups:       req.Notifications.PushGroups,
		PushEvents:       req.Notifications.PushEvents,
		PushBirthdays:    req.Notifications.PushBirthdays,
	}

	return e
}
func ToCreateIniEntityUserSetting() *entity.UserSetting {
	e := &entity.UserSetting{
		// 1. Identity
		User_ID: req.User_ID,
		Theme_Mode:     req.Theme_Mode,
		Font_Size:      req.Font_Size,
		Compact_Mode:   req.Compact_Mode,
		Lang_Code:      req.Lang_Code,
		Timezone:       req.Timezone,
		Auto_Translate: req.Auto_Translate,

		// --- 3. PRIVACY DEFAULTS ---
		Default_Post_Audience:  req.Default_Post_Audience,
		Default_Story_Audience: req.Default_Story_Audience,

		// --- 4. ACCESS CONTROL ---
		Allow_Friend_Request_From:    req.Allow_Friend_Request_From,
		Allow_Friend_List_View_From:  req.Allow_Friend_List_View_From,
		Allow_Search_Engine_Indexing: req.Allow_Search_Engine_Indexing,

		// --- 5. TIMELINE & TAGGING ---
		Allow_Timeline_Posting_From:   req.Allow_Timeline_Posting_From,
		Review_Tags_Enabled:           req.Review_Tags_Enabled,
		Review_Timeline_Posts_Enabled: req.Review_Timeline_Posts_Enabled,

		// --- 6. METADATA ---
		Updated_At: time.Now(), // Luôn lấy giờ hiện tại của server
	}

	// --- 7. XỬ LÝ NESTED STRUCT (NOTIFICATIONS) ---
	// Kiểm tra nil để tránh panic hoặc ghi đè dữ liệu rỗng
	e.Notifications = &entity.NotificationSettings{
		EmailFrequency:   req.Notifications.EmailFrequency,
		PushInteractions: req.Notifications.PushInteractions,
		PushFriends:      req.Notifications.PushFriends,
		PushGroups:       req.Notifications.PushGroups,
		PushEvents:       req.Notifications.PushEvents,
		PushBirthdays:    req.Notifications.PushBirthdays,
	}

	return e
}