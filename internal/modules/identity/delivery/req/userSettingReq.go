package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSettingReq struct {
	ID      primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	User_ID string             `bson:"user_id" json:"user_id"`

	// --- 1. APP CONFIGURATION ---
	Theme_Mode   enum.ThemeMode `bson:"theme_mode" json:"theme_mode"`     // 'system', 'light', 'dark', 'high_contrast'
	Font_Size    enum.FontSize  `bson:"font_size" json:"font_size"`       // 'medium', 'large'
	Compact_Mode bool           `bson:"compact_mode" json:"compact_mode"` //

	Lang_Code      enum.LangCode `bson:"lang_code" json:"lang_code"`           // 'vi', 'en', 'jp'
	Timezone       enum.Timezone `bson:"timezone" json:"timezone"`             // 'Asia/Ho_Chi_Minh'
	Auto_Translate bool          `bson:"auto_translate" json:"auto_translate"` //

	// --- 2. PRIVACY DEFAULTS ---
	Default_Post_Audience  enum.PrivacyLevel `bson:"default_post_audience" json:"default_post_audience"`   // 'public', 'friends', 'only_me'
	Default_Story_Audience enum.PrivacyLevel `bson:"default_story_audience" json:"default_story_audience"` // 'friends', 'close_friends'

	// --- 3. ACCESS CONTROL ---
	Allow_Friend_Request_From    enum.PrivacyLevel `bson:"allow_friend_request_from" json:"allow_friend_request_from"`       // 'everyone', 'friends_of_friends'
	Allow_Friend_List_View_From  enum.PrivacyLevel `bson:"allow_friend_list_view_from" json:"allow_friend_list_view_from"`   // 'public', 'friends', 'only_me'
	Allow_Email_Lookup_From      enum.PrivacyLevel `bson:"allow_email_lookup_from" json:"allow_email_lookup_from"`           // 'everyone', 'friends'
	Allow_Phone_Lookup_From      enum.PrivacyLevel `bson:"allow_phone_lookup_from" json:"allow_phone_lookup_from"`           // 'everyone', 'friends'
	Allow_Search_Engine_Indexing bool              `bson:"allow_search_engine_indexing" json:"allow_search_engine_indexing"` //
	// --- 4. TIMELINE & TAGGING ---
	Allow_Timeline_Posting_From   enum.PrivacyLevel `bson:"allow_timeline_posting_from" json:"allow_timeline_posting_from"`     // 'friends', 'only_me'
	Review_Tags_Enabled           bool              `bson:"review_tags_enabled" json:"review_tags_enabled"`                     //
	Review_Timeline_Posts_Enabled bool              `bson:"review_timeline_posts_enabled" json:"review_timeline_posts_enabled"` //

	// --- 5. NOTIFICATIONS ---
	Notifications *NotificationSettings `bson:"notifications" json:"notifications"`
	Updated_At    time.Time             `bson:"updated_at" json:"updated_at"`
}
type NotificationSettings struct {
	EmailFrequency   string `bson:"email_frequency" json:"email_frequency"`
	PushInteractions bool   `bson:"push_interactions" json:"push_interactions"`
	PushFriends      bool   `bson:"push_friends" json:"push_friends"`
	PushGroups       bool   `bson:"push_groups" json:"push_groups"`
	PushEvents       bool   `bson:"push_events" json:"push_events"`
	PushBirthdays    bool   `bson:"push_birthdays" json:"push_birthdays"`
}

func (req *UserSettingReq) ToEntity(userID string) *entity.UserSetting {
	e := &entity.UserSetting{

		// 1. Identity
		User_ID: userID,

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
		Allow_Email_Lookup_From:      req.Allow_Email_Lookup_From,
		Allow_Phone_Lookup_From:      req.Allow_Phone_Lookup_From,
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
