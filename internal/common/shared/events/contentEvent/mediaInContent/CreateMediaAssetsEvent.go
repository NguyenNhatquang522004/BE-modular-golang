package mediaInContent

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	reqcontent "github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- 2. DỮ LIỆU CỐT LÕI (Payload) ---
type CreateMediaAssetsPayload struct {
	// Chuyển toàn bộ primitive.ObjectID của MongoDB thành string
	UserID string             `json:"user_id"` // ID người upload (UUID từ Postgres)
	PostID string             `json:"post_id"`
	Items  []MediaItemPayload `json:"items"`
}

// --- 3. SUB-STRUCT: MEDIA ITEM ---
type MediaItemPayload struct {
	MediaID      string                `json:"media_id"`
	MediaType    sharedEnums.MediaType `json:"media_type"` // Sử dụng Enum đã định nghĩa
	URL          string                `json:"url"`
	ThumbnailURL string                `json:"thumbnail_url"`
	Metadata     MetadataPayload       `json:"metadata"`
	Order        int                   `json:"order"`

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

func CreateMediaAssetsPayloadMediaMetadataReqtoPayloads(request *reqcontent.PostMediaReq, userid string) *CreateMediaAssetsPayload {
	// 1. KIỂM TRA NIL NGAY TỪ ĐẦU
	if request == nil {
		return nil // Hoặc trả về payload rỗng tùy theo logic xử lý lỗi ở tầng trên
	}

	// 2. TỐI ƯU HIỆU NĂNG: Cấp phát sẵn bộ nhớ cho slice bằng len(request.Items)
	itemsPayload := make([]MediaItemPayload, 0, len(request.Items))

	for _, item := range request.Items {
		// Bỏ qua nếu có phần tử nil bất thường trong mảng
		if item == nil {
			continue
		}

		// Xử lý TaggedUsers tách biệt cho dễ đọc
		var taggedUsers []TaggedUserPayload
		if len(item.TaggedUsers) > 0 {
			taggedUsers = make([]TaggedUserPayload, 0, len(item.TaggedUsers))
			for _, tu := range item.TaggedUsers {
				taggedUsers = append(taggedUsers, TaggedUserPayload{
					UserID: tu.UserID,
					Name:   tu.Name,
					X:      tu.X,
					Y:      tu.Y,
				})
			}
		}

		// Gộp vào mảng itemsPayload
		itemsPayload = append(itemsPayload, MediaItemPayload{
			MediaID:      item.ID,
			MediaType:    item.MediaType,
			URL:          item.URL,
			ThumbnailURL: item.ThumbnailURL,
			Order:        item.Order,
			Metadata: MetadataPayload{
				Width:     item.Metadata.Width,
				Height:    item.Metadata.Height,
				Duration:  item.Metadata.Duration,
				SizeBytes: item.Metadata.SizeBytes,
				MimeType:  item.Metadata.MimeType,
			},
			TaggedUsers: taggedUsers,
		})
	}

	return &CreateMediaAssetsPayload{
		PostID: request.PostID,
		UserID: userid,
		Items:  itemsPayload,
	}
}
func CreateMediaAssetsPayloadtoEntityMediaAssets(payload *CreateMediaAssetsPayload) ([]*entity.MediaAsset, error) {
	if payload == nil {
		return nil, nil
	}

	// 1. Convert PostID một lần duy nhất
	postID, err := primitive.ObjectIDFromHex(payload.PostID)
	if err != nil {
		// Trong thực tế, nếu PostID sai format thì event này có thể bị coi là invalid
		return nil, err
	}

	assets := make([]*entity.MediaAsset, 0, len(payload.Items))

	for _, item := range payload.Items {
		// 2. Convert MediaID (nếu từ Content gửi sang đã có sẵn ID)
		mediaID, err := primitive.ObjectIDFromHex(item.MediaID)
		if err != nil {
			mediaID = primitive.NewObjectID() // Nếu lỗi thì tạo mới để đảm bảo data toàn vẹn
		}

		// 3. Khởi tạo Entity cho từng item
		asset := &entity.MediaAsset{
			ID:            mediaID,
			UserID:        payload.UserID, // Lấy từ payload
			PostID:        postID,
			OriginalURL:   item.URL,
			ThumbnailURL:  item.ThumbnailURL,
			StorageFileID: item.URL,       // Mapping theo yêu cầu của bạn
			AssetType:     item.MediaType, // Cast Enum
			Order:         item.Order,

			// Metadata
			Metadata: entity.MediaMetadata{
				Width:     item.Metadata.Width,
				Height:    item.Metadata.Height,
				Duration:  item.Metadata.Duration,
				SizeBytes: item.Metadata.SizeBytes,
				MimeType:  item.Metadata.MimeType,
			},

			// Tagged Users
			TaggedUsers: mapTaggedUsers(item.TaggedUsers),

			// Default Stats cho ảnh mới
			ReactionsCount: entity.MediaReactionStats{
				Total: 0,
				Like:  0,
				Love:  0,
			},
			CommentCount: 0,

			// Timestamps
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assets = append(assets, asset)
	}

	return assets, nil
}

// Hàm helper để map Tagged Users (giúp code chính sạch hơn)
func mapTaggedUsers(payloadTags []TaggedUserPayload) []entity.MediaTag {
	if len(payloadTags) == 0 {
		return nil
	}

	tags := make([]entity.MediaTag, len(payloadTags))
	for i, t := range payloadTags {
		tags[i] = entity.MediaTag{
			UserID: t.UserID,
			Name:   t.Name,
			Position: entity.TagPosition{
				X: t.X,
				Y: t.Y,
			},
			Status: sharedEnums.ProcessingActive, // Best practice: Mặc định là Active
		}
	}
	return tags
}
