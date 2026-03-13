package mediaEvent

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/enum"
)

type ReplyStoryPayload struct {
	StoryID string `json:"story_id"`
	ReplyID int    `json:"reply_id"`
}
type StartStopVideoLiveStreamPayload struct {
	LiveSessionID string              `json:"live_session_id"`
	SegmentLen    int                 `json:"segment_len"`
	OwnerID       string              `json:"owner_id"`
	Name          string              `json:"name"`
	EventType     constants.EventType `json:"event_type"`
}

type ReactLiveStreamPayload struct {
	LiveSessionID string                     `json:"live_session_id"`
	UserID        string                     `json:"user_id"`
	Total         int                        `json:"total"`
	TargetType    sharedEnums.ReactionTarget `json:"target_type" validate:"required"`
	ReactionCode  sharedEnums.ReactionCode   `json:"reaction_code"` // "❤️", "😂" hoặc ID sticker
	CreatedAt     time.Time                  `json:"created_at"`
	EventType     constants.EventType        `json:"event_type"` // "view" hoặc "unview"
}
type ReactCounterReelPayload struct {
	ReelID    string              `json:"reel_id"`
	UserID    string              `json:"user_id"`
	Comments  int                 `json:"comments"`
	Saves     int                 `json:"saves"`
	Shares    int                 `json:"shares"`
	EventType constants.EventType `json:"event_type"` // "view" hoặc "unview"
}

type ReactReelPayload struct {
	ReelID       string                     `json:"reel_id"`
	UserID       string                     `json:"user_id"`
	Total        int                        `json:"total"`
	TargetType   sharedEnums.ReactionTarget `json:"target_type" validate:"required"`
	ReactionCode sharedEnums.ReactionCode   `json:"reaction_code"` // "❤️", "😂" hoặc ID sticker
	CreatedAt    time.Time                  `json:"created_at"`
	EventType    constants.EventType        `json:"event_type"` // "view" hoặc "unview"
}

type DeleteMediaRelationTargetPayload struct {
	TargetID string `json:"target_id"`
}

// /
type CreateMediaAssetsPayload struct {
	// Chuyển toàn bộ primitive.ObjectID của MongoDB thành string
	// // Nếu client gửi lên có nghĩa là update, nếu không có nghĩa là create mới
	UserID string             `json:"user_id"` // ID người upload (UUID từ Postgres)
	Items  []MediaItemPayload `json:"items"`
}

// --- 3. SUB-STRUCT: MEDIA ITEM ---
type MediaItemPayload struct {
	MediaID      string                `json:"media_id,omitempty"` // Nếu client gửi lên có nghĩa là update, nếu không có nghĩa là create mới
	PostID       string                `json:"post_id"`
	AlbumID      string                `json:"album_id,omitempty"`
	GroupID      string                `json:"group_id,omitempty"`
	CommentID    string                `json:"comment_id,omitempty"`
	PageID       string                `json:"page_id,omitempty"`
	StoryID      string                `json:"story_id,omitempty"`
	ReelID       string                `json:"reel_id,omitempty"`
	MessageID    string                `json:"message_id,omitempty"`
	MediaType    sharedEnums.MediaType `json:"media_type"` // Sử dụng Enum đã định nghĩa
	URL          string                `json:"url"`
	ThumbnailURL string                `json:"thumbnail_url"`
	Metadata     MetadataPayload       `json:"metadata"`
	Order        int                   `json:"order"`
	Hashtags     []string              `json:"hashtags,omitempty"`
	// Lưu ý: Nếu module Media KHÔNG quan tâm đến TaggedUser (chỉ quan tâm xử lý file/ảnh),
	// bạn hoàn toàn có thể lược bỏ TaggedUsers ở đây để giữ payload nhẹ (Thin Payload).
	// Dưới đây vẫn giữ lại để đảm bảo đủ 100% data như entity của bạn.
	TaggedUsers []TaggedUserPayload `json:"tagged_users,omitempty"`
}

// --- 4. SUB-STRUCT: METADATA ---
type MetadataPayload struct {
	Width     int     `json:"width,omitempty"`
	Height    int     `json:"height,omitempty"`
	Duration  float64 `json:"duration,omitempty"`
	SizeBytes int64   `json:"size_bytes"`
	MimeType  string  `json:"mime_type"`
}

// --- 5. SUB-STRUCT: TAGGED USER ---
type TaggedUserPayload struct {
	UserID string  `json:"user_id"` // Đã là string từ Postgres UUID
	Name   string  `json:"name"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
}

type UpdateMediaAssetsPayload struct {
	MediaID string `json:"media_id"`
	// Các trường có thể cập nhật
	AlbumID      string               `json:"album_id,omitempty"`
	URL          *string              `json:"url,omitempty"`
	ThumbnailURL *string              `json:"thumbnail_url,omitempty"`
	Order        *int                 `json:"order,omitempty"`
	Hashtags     *[]string            `json:"hashtags,omitempty"`
	Metadata     *MetadataPayload     `json:"metadata,omitempty"`
	TaggedUsers  *[]TaggedUserPayload `json:"tagged_users,omitempty"`
	Privacy      *MediaPrivacyPayload `json:"privacy,omitempty"`
}

type DeleteMediaAssetsPayload struct {
	MediaID   string `json:"media_id"`
	MessageID string `json:"message_id"`
	GroupID   string `json:"group_id"`
	PageID    string `json:"page_id"`
}
type MediaPrivacyPayload struct {
	// Level string hoặc dùng Enum PrivacyScope tái sử dụng
	Level            sharedEnums.PrivacyScope `bson:"level" json:"level"`
	InheritFromAlbum bool                     `bson:"inherit_from_album" json:"inherit_from_album"`
}
type AlbumStatsPayload struct {
	UserID     string              `json:"user_id"`
	AlbumID    string              `json:"album_id"`
	AssetCount int                 `json:"asset_count"`
	Like       int                 `json:"like"`
	Love       int                 `json:"love"`
	Haha       int                 `json:"haha"`
	Wow        int                 `json:"wow"`
	Sad        int                 `json:"sad"`
	Angry      int                 `json:"angry"`
	EventType  constants.EventType `json:"event_type"` // "increment" hoặc "decrement"
}

type StoryStatsPayload struct {
	UserID          string                           `json:"user_id"`
	StoryID         string                           `json:"story_id"`
	Views           int                              `json:"views"`
	Like            int                              `json:"like"`
	Love            int                              `json:"love"`
	Haha            int                              `json:"haha"`
	Wow             int                              `json:"wow"`
	Sad             int                              `json:"sad"`
	Angry           int                              `json:"angry"`
	ReplyCount      int                              `json:"reply_count"`
	ViewsCount      int                              `json:"views_count"`
	InteractionType sharedEnums.StoryInteractionType `json:"interaction_type"`
	PollOptionIndex *int                             `json:"poll_option_index,omitempty"` // Chỉ có khi InteractionType là
	Content         string                           `json:"content,omitempty"`           // Chỉ có khi InteractionType là reaction hoặc comment
	EventType       constants.EventType              `json:"event_type"`                  // "increment" hoặc "decrement"
}

type ReelStatsPayload struct {
	UserID    string              `json:"user_id"`
	ReelID    string              `json:"reel_id"`
	Views     int                 `json:"views"`
	Like      int                 `json:"like"`
	Love      int                 `json:"love"`
	Haha      int                 `json:"haha"`
	Wow       int                 `json:"wow"`
	Sad       int                 `json:"sad"`
	Angry     int                 `json:"angry"`
	Shares    int                 `json:"shares"`
	Saves     int                 `json:"saves"`
	Comments  int                 `json:"comments"`
	EventType constants.EventType `json:"event_type"` // "increment" hoặc "decrement"
}

type LiveSessionStatsPayload struct {
	LiveSessionID string              `json:"live_session_id"`
	UserID        string              `json:"user_id"`
	PeakViewers   int                 `json:"peak_viewers"`
	TotalViews    int                 `json:"total_views"`
	TotalComments int                 `json:"total_comments"`
	Like          int                 `json:"like"`
	Love          int                 `json:"love"`
	Haha          int                 `json:"haha"`
	Wow           int                 `json:"wow"`
	Sad           int                 `json:"sad"`
	Angry         int                 `json:"angry"`
	EventType     constants.EventType `json:"event_type"` // "view" hoặc "unview"
}

type LiveCommentPayload struct {
	UserID        string                   `json:"user_id"`
	StreamID      string                   `json:"stream_id"`
	CreatedAt     time.Time                `json:"created_at"`
	CommentID     string                   `json:"comment_id"`
	UserBadges    []*sharedEnums.UserBadge `json:"user_badges"`
	Content       string                   `json:"content"`
	IsPinned      bool                     `json:"is_pinned"`
	UserNickname  string                   `json:"user_nickname"`
	UserAvatarURL string                   `json:"user_avatar_url"`
	EventType     constants.EventType      `json:"event_type"`
}

// CreateStoryPayload đại diện cho payload Client gửi lên khi tạo Story mới
type CreateStoryPayload struct {
	UserID   string                `json:"user_id" binding:"required"` // ID người tạo Story (UUID từ Postgres)
	Media    StoryMediaPayload     `json:"media" binding:"required"`
	Overlays []StoryOverlayPayload `json:"overlays,omitempty" binding:"dive"` // dive: validate từng phần tử trong mảng
	Privacy  StoryPrivacyPayload   `json:"privacy" binding:"required"`
	Settings StorySettingsPayload  `json:"settings" binding:"required"`
}

type StoryMediaPayload struct {
	URL          string                `json:"url" binding:"required,url"` // Bắt buộc phải là định dạng URL
	Type         sharedEnums.MediaType `json:"type" binding:"required"`
	Duration     float64               `json:"duration" binding:"gte=0"` // Lớn hơn hoặc bằng 0
	ThumbnailURL string                `json:"thumbnail_url" binding:"omitempty,url"`
	SizeBytes    int64                 `json:"size_bytes" binding:"gte=0"`
	Width        int                   `json:"width,omitempty" binding:"gte=0"`
	Height       int                   `json:"height,omitempty" binding:"gte=0"`
	MimeType     string                `json:"mime_type,omitempty"`
}

type StoryPrivacyPayload struct {
	Type sharedEnums.PrivacyScope `json:"type" binding:"required"`
	// Nếu dùng UUID cho Postgres, validate uuid ở đây
	AllowList []string `json:"allow_list,omitempty" binding:"omitempty,dive,uuid"`
	BlockList []string `json:"block_list,omitempty" binding:"omitempty,dive,uuid"`
}

type StorySettingsPayload struct {
	// Không dùng pointer ở Create vì ta cần force client gửi các config này
	AllowReply bool `json:"allow_reply"`
	AllowShare bool `json:"allow_share"`
}

type StoryOverlayPayload struct {
	Type     sharedEnums.OverlayType `json:"type" binding:"required"`
	Position OverlayPositionPayload  `json:"position" binding:"required"`
	Data     map[string]interface{}  `json:"data" binding:"required"`
}

type OverlayPositionPayload struct {
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Rotation float64 `json:"rotation"`
	Scale    float64 `json:"scale" binding:"gt=0"` // Scale phải lớn hơn 0
}
type DeleteStoryPayload struct {
	StoryID string `json:"story_id" binding:"required"`
}
type UpdateStoryPayload struct {
	StoryID  string                      `json:"story_id" binding:"required"`
	Privacy  *UpdateStoryPrivacyPayload  `json:"privacy,omitempty"`
	Settings *UpdateStorySettingsPayload `json:"settings,omitempty"`
}

type UpdateStoryPrivacyPayload struct {
	Type      *sharedEnums.PrivacyScope `json:"type,omitempty"`
	AllowList []string                  `json:"allow_list,omitempty" binding:"omitempty,dive,uuid"`
	BlockList []string                  `json:"block_list,omitempty" binding:"omitempty,dive,uuid"`
}

type UpdateStorySettingsPayload struct {
	// Bắt buộc phải dùng con trỏ (pointer) cho boolean trong Update DTO.
	// Nếu dùng bool thường, khi client không truyền `allow_reply`, Go sẽ tự hiểu là `false`.
	// Dùng `*bool` giúp ta check: nếu nó là nil -> client không muốn update trường này.
	AllowReply *bool `json:"allow_reply,omitempty"`
	AllowShare *bool `json:"allow_share,omitempty"`
}
type DeleteReelPayload struct {
	UserID string `json:"user_id" binding:"required"` // ID người tạo Reel (UUID từ Postgres)
	PostID string `json:"post_id" binding:"required"`
	ReelID string `json:"reel_id" binding:"required"`
}
type UpdateReelPayload struct {
	UserID       string  `json:"user_id" binding:"required"` // ID người tạo Reel (UUID từ Postgres)
	ReelID       string  `json:"reel_id" binding:"required"`
	Caption      *string `json:"caption,omitempty" binding:"omitempty,max=2200"`
	AllowComment *bool   `json:"allow_comment,omitempty"`
	AllowShare   *bool   `json:"allow_share,omitempty"`

	// Slice bản chất đã có thể check nil, nhưng để phân biệt "xóa hết hashtag" (gửi mảng rỗng [])
	// và "không update hashtag" (không gửi field), ta có thể cân nhắc dùng Pointer cho mảng hoặc tự xử lý logic ở Service.
	Hashtags []string `json:"hashtags,omitempty" binding:"omitempty,dive,alphanum"`
	Mentions []string `json:"mentions,omitempty" binding:"omitempty,dive,uuid"`

	Privacy *sharedEnums.PrivacyScope `json:"privacy,omitempty" binding:"omitempty"`

	// LƯU Ý: Thông thường các nền tảng (như TikTok/IG) KHÔNG cho phép update Video file,
	// Audio file hay Remix Info sau khi đã publish. Nếu hệ thống của bạn cho phép,
	// bạn có thể thêm các trường UpdateVideoDTO, UpdateAudioMetaDTO tương tự vào đây bằng Pointer.
}
type CreateReelPayload struct {
	UserID string `json:"user_id" binding:"required"` // ID người tạo Reel (UUID từ Postgres)
	// Video bắt buộc phải có khi tạo Reel
	Video ReelVideoPayload `json:"video" binding:"required"`

	Caption      string                   `json:"caption" binding:"max=2200"`                 // Giới hạn độ dài giống Instagram
	Hashtags     []string                 `json:"hashtags,omitempty" binding:"dive,alphanum"` // dive: validate từng phần tử trong slice
	Mentions     []string                 `json:"mentions,omitempty" binding:"dive,uuid"`     // Theo entity, UserID là UUID (Postgres)
	AllowComment *bool                    `json:"allow_comment,omitempty"`
	AllowShare   *bool                    `json:"allow_share,omitempty"`
	Privacy      sharedEnums.PrivacyScope `json:"privacy" binding:"required"`

	AudioMeta AudioMetaPayload `json:"audio_meta" binding:"required"`

	// Con trỏ vì không phải Reel nào cũng là Remix
	RemixInfo *RemixInfoPayload `json:"remix_info,omitempty"`
}

type ReelVideoPayload struct {
	URL           string  `json:"url" binding:"required,url"`
	ThumbnailURL  string  `json:"thumbnail_url" binding:"required,url"`
	PreviewGifURL string  `json:"preview_gif_url,omitempty" binding:"omitempty,url"`
	Width         int     `json:"width" binding:"required,min=1"`
	Height        int     `json:"height" binding:"required,min=1"`
	Duration      float64 `json:"duration" binding:"required,min=0.1"` // Tính bằng giây
	SizeBytes     int64   `json:"size_bytes" binding:"required,min=1"`
	MimeType      string  `json:"mime_type,omitempty"`
}

type AudioMetaPayload struct {
	// Nhận string từ client, Service layer sẽ convert sang primitive.ObjectID
	TrackID         *string `json:"track_id,omitempty" binding:"omitempty,mongodb"`
	IsOriginalAudio bool    `json:"is_original_audio"`
	VolumeAdjust    float64 `json:"volumn_adjust" binding:"min=0,max=1"` // 0.0 -> 1.0 (Giữ nguyên tên field volumn)
	AudioStartTime  float64 `json:"audio_start_time" binding:"min=0"`
}

type RemixInfoPayload struct {
	// Nhận string từ client, Service layer sẽ convert sang primitive.ObjectID
	ParentReelID string         `json:"parent_reel_id" binding:"required,mongodb"`
	Type         enum.RemixType `json:"type" binding:"required"`
	IsRemixable  bool           `json:"is_remixable"`
}

type AlbumPrivacyPayload struct {
	// Yêu cầu phải có, có thể thêm tag 'oneof' nếu enum là string (ví dụ: oneof=public friends only_me custom)
	Level sharedEnums.PrivacyScope `json:"level" binding:"required"`

	// dive,uuid: Đảm bảo từng phần tử trong mảng phải là UUID hợp lệ
	AllowList []string `json:"allow_list" binding:"omitempty,dive,uuid"`
	BlockList []string `json:"block_list" binding:"omitempty,dive,uuid"`
}
type CreateAlbumPayload struct {
	UserID string
	// GroupID nhận vào dạng string để validate trước khi parse sang ObjectID ở Controller/Service.
	// Tránh lỗi panic của Unmarshal nếu client gửi sai format ObjectID.
	GroupID *string `json:"group_id" binding:"omitempty,mongodb"`

	Title       string         `json:"title" binding:"required,min=1,max=255"`
	Description string         `json:"description" binding:"omitempty,max=2000"`
	Type        enum.AlbumType `json:"type" binding:"required"`

	Privacy      AlbumPrivacyPayload `json:"privacy" binding:"required"`
	ItemMediaIDs []string            `json:"item_media_ids,omitempty" binding:"omitempty,dive,mongodb"` // Danh sách media_id (string) để thêm vào album khi tạo
}
type DeleteAlbumPayload struct {
	AlbumID      string   `json:"album_id" binding:"required"`
	ItemMediaIDs []string `json:"item_media_ids,omitempty" binding:"omitempty,dive,mongodb"` // Danh sách media_id (string) để xóa khỏi album khi xóa album
}
type UpdateAlbumPayload struct {
	ID                 string          `json:"id" binding:"required"` // ID của album cần update, bắt buộc phải có để xác định target
	Title              *string         `json:"title" binding:"omitempty,min=1,max=255"`
	Description        *string         `json:"description" binding:"omitempty,max=2000"`
	Type               *enum.AlbumType `json:"type" binding:"omitempty"`
	ItemMediaDeleteIDs []string        `json:"item_media_delete_ids,omitempty" binding:"omitempty,dive,mongodb"` // Danh sách media_id (string) để xóa khỏi album khi cập nhật
	ItemMediaAddIDs    []string        `json:"item_media_add_ids,omitempty" binding:"omitempty,dive,mongodb"`    // Danh sách media_id (string) để thêm vào album khi cập nhật
	// CoverAssetID có thể được update sau khi người dùng upload ảnh mới
	CoverAssetID *string `json:"cover_asset_id" binding:"omitempty,mongodb"`

	// Pointer tới struct để biết client có muốn update privacy hay không
	Privacy *AlbumPrivacyPayload `json:"privacy" binding:"omitempty"`
}
type CreateMusicLibraryPayload struct {
	// Dùng string cho ArtistID ở DTO để dễ validate, sau đó convert sang ObjectID ở Service
	ArtistID string `json:"artist_id" binding:"required,mongodb"`
	Title    string `json:"title" binding:"required,max=255"`
	Album    string `json:"album,omitempty" binding:"omitempty,max=255"`

	CoverURL  string `json:"cover_url" binding:"required,url"`
	StreamURL string `json:"stream_url" binding:"required,url"`
	Duration  int    `json:"duration" binding:"required,gt=0"` // Phải lớn hơn 0

	LyricsSnippet string                   `json:"lyrics_snippet,omitempty"`
	Genres        []sharedEnums.MusicGenre `json:"genre" binding:"required,min=1"` // Ít nhất 1 thể loại

	CopyrightInfo CreateCopyrightInfoPayload `json:"copyright_info" binding:"required"`
}

type CreateCopyrightInfoPayload struct {
	Provider       string   `json:"provider" binding:"required"`
	AllowedRegions []string `json:"allowed_regions,omitempty" binding:"omitempty,dive,iso3166_1_alpha2"` // Validate chuẩn mã quốc gia 2 ký tự (VD: VN, US)
}
type DeleteMusicLibraryPayload struct {
	MusicID string `json:"music_id" binding:"required"`
	ArtisID string `json:"artist_id" binding:"required,mongodb"`
}
type UpdateMusicLibraryPayload struct {
	// Dùng con trỏ (*) cho TẤT CẢ các trường để phân biệt nil (không gửi) và zero-value (0, "", false)
	// ArtistID thường là immutable (không cho phép đổi sau khi tạo), nếu hệ thống cho phép đổi thì bạn mới thêm vào đây.
	UserID  string  `json:"user_id" binding:"required"`  // ID người tạo (UUID từ Postgres)
	MusicID string  `json:"music_id" binding:"required"` // ID của bản nhạc cần update, bắt buộc phải có để xác định target
	Title   *string `json:"title,omitempty" binding:"omitempty,max=255"`
	Album   *string `json:"album,omitempty" binding:"omitempty,max=255"`

	CoverURL  *string `json:"cover_url,omitempty" binding:"omitempty,url"`
	StreamURL *string `json:"stream_url,omitempty" binding:"omitempty,url"`
	Duration  *int    `json:"duration,omitempty" binding:"omitempty,gt=0"`

	LyricsSnippet *string                   `json:"lyrics_snippet,omitempty"`
	Genres        *[]sharedEnums.MusicGenre `json:"genre,omitempty" binding:"omitempty,min=1"`

	CopyrightInfo *UpdateCopyrightInfoPayload `json:"copyright_info,omitempty"`
}

type UpdateCopyrightInfoPayload struct {
	Provider       *string   `json:"provider,omitempty"`
	AllowedRegions *[]string `json:"allowed_regions,omitempty" binding:"omitempty,dive,iso3166_1_alpha2"`
}
type SocialLinksPayload struct {
	Spotify   string `json:"spotify,omitempty" validate:"omitempty,url"`
	Youtube   string `json:"youtube,omitempty" validate:"omitempty,url"`
	Instagram string `json:"instagram,omitempty" validate:"omitempty,url"`
	Facebook  string `json:"facebook,omitempty" validate:"omitempty,url"`
	Website   string `json:"website,omitempty" validate:"omitempty,url"`
}
type CreateArtistPayload struct {
	ArtistID string `json:"artist_id,omitempty"` // Nếu client gửi lên có nghĩa là update, nếu không có nghĩa là create mới
	// Name là bắt buộc khi tạo mới
	Name string `json:"name" validate:"required,min=2,max=100"`

	// Slug thường được hệ thống tự động sinh ra từ Name (vd: son-tung-m-tp)
	// Nhưng nếu bạn cho phép user tự custom slug thì mở field này.
	Slug string `json:"slug,omitempty" validate:"omitempty,min=2,max=100"`

	Bio       string `json:"bio,omitempty" validate:"omitempty,max=1000"`
	AvatarURL string `json:"avatar_url,omitempty" validate:"omitempty,url"`
	CoverURL  string `json:"cover_url,omitempty" validate:"omitempty,url"`

	// UserID thường được trích xuất từ Token JWT (middleware) gán vào context.
	// Chỉ đưa vào payload nếu đây là API dành cho Admin tạo hộ Artist.
	UserID string `json:"user_id,omitempty" validate:"omitempty,mongodb"`

	SocialLinks *SocialLinksPayload `json:"social_links,omitempty"`
}
type DeleteArtistPayload struct {
	ArtisID string `json:"artist_id" binding:"required"` // ID của artist cần xóa, bắt buộc phải có để xác định target
}
type UpdateArtistPayload struct {
	// Con trỏ *string giúp phân biệt:
	// - Client không gửi field "name" lên API -> Name = nil -> Không cập nhật
	// - Client gửi "name": "" -> Name != nil -> Báo lỗi validation do min=2
	ArtisID string  `json:"artist_id" binding:"required"` // ID của artist cần update, bắt buộc phải có để xác định target
	Name    *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`

	Slug      *string `json:"slug,omitempty" validate:"omitempty,min=2,max=100"`
	Bio       *string `json:"bio,omitempty" validate:"omitempty,max=1000"`
	AvatarURL *string `json:"avatar_url,omitempty" validate:"omitempty,url"`
	CoverURL  *string `json:"cover_url,omitempty" validate:"omitempty,url"`

	// Update cả object SocialLinks hoặc không update.
	SocialLinks *SocialLinksPayload `json:"social_links,omitempty"`
}
type MusicStatsPayload struct {
	ID         string `json:"id"` // Có thể là MusicID hoặc ArtistID tùy context
	ArtistID   string `json:"artist_id"`
	UsageCount int    `json:"usage_count"`
}
type ArtistStatsPayload struct {
	ArtistID      string `json:"artist_id"`
	FollowerCount int    `json:"follower_count"`
	TotalStreams  int    `json:"total_streams"`
}

type ProcessMediaPayload struct {
	// Bắt buộc: ID của MediaAsset trong MongoDB (dạng Hex String)
	MediaID string `json:"media_id"`

	// Tùy chọn: Gắn thêm UserID để dễ dàng trace log trên Kibana/Grafana mà không cần query DB
	UserID string `json:"user_id,omitempty"`
}
