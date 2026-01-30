package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ThemeMode string

const (
	ThemeSystem       ThemeMode = "system"
	ThemeLight        ThemeMode = "light"
	ThemeDark         ThemeMode = "dark"
	ThemeHighContrast ThemeMode = "high_contrast"
)

type FontSize string

const (
	FontMedium FontSize = "medium"
	FontLarge  FontSize = "large"
)

type PrivacyLevel string

const (
	PrivacyPublic           PrivacyLevel = "public"
	PrivacyFriends          PrivacyLevel = "friends"
	PrivacyCloseFriends     PrivacyLevel = "close_friends"
	PrivacyOnlyMe           PrivacyLevel = "only_me"
	PrivacyFriendsOfFriends PrivacyLevel = "friends_of_friends"
	PrivacyEveryone         PrivacyLevel = "everyone"
)

type EmailFrequency string

const (
	EmailDaily  EmailFrequency = "daily"
	EmailWeekly EmailFrequency = "weekly"
	EmailNever  EmailFrequency = "never"
)

type Lang_Code string

const (
	LangVi Lang_Code = "vi"
	LangEn Lang_Code = "en"
	LangJp Lang_Code = "jp"
)

type Timezone string

const (
	TimezoneHCM Timezone = "Asia/Ho_Chi_Minh"
	TimezoneNY  Timezone = "America/New_York"
	TimezoneLDN Timezone = "Europe/London"
)

type User_Settings struct {
	ID      primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	User_ID string             `bson:"user_id" json:"user_id"`

	// --- 1. APP CONFIGURATION ---
	Theme_Mode   ThemeMode `bson:"theme_mode" json:"theme_mode"`     // 'system', 'light', 'dark', 'high_contrast'
	Font_Size    FontSize  `bson:"font_size" json:"font_size"`       // 'medium', 'large'
	Compact_Mode bool      `bson:"compact_mode" json:"compact_mode"` //

	Lang_Code      Lang_Code `bson:"lang_code" json:"lang_code"`           // 'vi', 'en', 'jp'
	Timezone       Timezone  `bson:"timezone" json:"timezone"`             // 'Asia/Ho_Chi_Minh'
	Auto_Translate bool      `bson:"auto_translate" json:"auto_translate"` //

	// --- 2. PRIVACY DEFAULTS ---
	Default_Post_Audience  PrivacyLevel `bson:"default_post_audience" json:"default_post_audience"`   // 'public', 'friends', 'only_me'
	Default_Story_Audience PrivacyLevel `bson:"default_story_audience" json:"default_story_audience"` // 'friends', 'close_friends'

	// --- 3. ACCESS CONTROL ---
	Allow_Friend_Request_From    PrivacyLevel `bson:"allow_friend_request_from" json:"allow_friend_request_from"`       // 'everyone', 'friends_of_friends'
	Allow_Friend_List_View_From  PrivacyLevel `bson:"allow_friend_list_view_from" json:"allow_friend_list_view_from"`   // 'public', 'friends', 'only_me'
	Allow_Email_Lookup_From      PrivacyLevel `bson:"allow_email_lookup_from" json:"allow_email_lookup_from"`           // 'everyone', 'friends'
	Allow_Phone_Lookup_From      PrivacyLevel `bson:"allow_phone_lookup_from" json:"allow_phone_lookup_from"`           // 'everyone', 'friends'
	Allow_Search_Engine_Indexing bool         `bson:"allow_search_engine_indexing" json:"allow_search_engine_indexing"` //
	// --- 4. TIMELINE & TAGGING ---
	Allow_Timeline_Posting_From   PrivacyLevel `bson:"allow_timeline_posting_from" json:"allow_timeline_posting_from"`     // 'friends', 'only_me'
	Review_Tags_Enabled           bool         `bson:"review_tags_enabled" json:"review_tags_enabled"`                     //
	Review_Timeline_Posts_Enabled bool         `bson:"review_timeline_posts_enabled" json:"review_timeline_posts_enabled"` //

	// --- 5. NOTIFICATIONS ---
	Notifications *NotificationSettings `bson:"notifications" json:"notifications"`
	Updated_At    time.Time             `bson:"updated_at" json:"updated_at"`
}
type NotificationSettings struct {
	EmailFrequency   EmailFrequency `bson:"email_frequency" json:"email_frequency"`
	PushInteractions bool           `bson:"push_interactions" json:"push_interactions"`
	PushFriends      bool           `bson:"push_friends" json:"push_friends"`
	PushGroups       bool           `bson:"push_groups" json:"push_groups"`
	PushEvents       bool           `bson:"push_events" json:"push_events"`
	PushBirthdays    bool           `bson:"push_birthdays" json:"push_birthdays"`
}
