package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/contentEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToCreateEntityPostPayload(r *contentEvent.CreatePostPayload) *entity.Post {
	if r == nil {
		return nil
	}
	post := &entity.Post{
		UserID:  "", // UserID sẽ được gán sau khi giải mã token, không lấy từ payload
		Type:    r.Type,
		Content: r.Content,
		// Map explicitly to avoid conversion errors
		Privacy: entity.PostPrivacy{
			Scope:        r.Privacy.Scope,
			AllowComment: *r.Privacy.AllowComment,
			AllowShare:   *r.Privacy.AllowShare,
		},
		Context: &entity.PostContext{
			Type:     r.Context.Type,
			TargetID: r.Context.TargetID,
		},
		Hashtags:  r.Hashtags,
		Mentions:  r.Mentions,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return post
}
func ToCreateEntityPostMediaPayload(postID string, mediaItems []contentEvent.MediaItemPayload) *entity.PostMedia {
	convertedPostID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		// Xử lý lỗi nếu postID không hợp lệ
		return nil
	}
	var taggedUsers []entity.TaggedUser
	for _, item := range mediaItems {
		for _, tagged := range item.TaggedUsers {
			taggedUser := entity.TaggedUser{
				UserID: tagged.UserID,
				Name:   tagged.Name,
				X:      tagged.X,
				Y:      tagged.Y,
			}
			taggedUsers = append(taggedUsers, taggedUser)
		}
	}
	var items []*entity.MediaItem
	for index, item := range mediaItems {
		mediaItem := &entity.MediaItem{
			ID:           primitive.NewObjectID(),
			MediaType:    item.MediaType,
			URL:          item.URL,
			ThumbnailURL: item.ThumbnailURL,
			Metadata: entity.MediaMetadata{
				Width:     item.Width,
				Height:    item.Height,
				Duration:  item.Duration,
				SizeBytes: item.SizeBytes,
				MimeType:  item.MimeType,
			},
			Order:       index,
			TaggedUsers: taggedUsers,
		}
		items = append(items, mediaItem)
	}
	entity := &entity.PostMedia{
		ID:     primitive.NewObjectID(),
		PostID: convertedPostID,
		Items:  items,
	}
	return entity
}

func ToCreateEntityPostExtensionPayload(postID string, extension *contentEvent.ExtensionPayload) *entity.PostExtension {
	if extension == nil {
		return nil
	}
	convertedPostID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		// Xử lý lỗi nếu postID không hợp lệ
		return nil
	}
	convertoriginalPostID, err := primitive.ObjectIDFromHex(extension.ShareData.OriginalPostID)
	entity := &entity.PostExtension{
		ID:     primitive.NewObjectID(),
		PostID: convertedPostID,
		ShareData: &entity.ShareData{
			ParentPostID:   convertedPostID,
			OriginalPostID: convertoriginalPostID,
		}, // Cần map chi tiết nếu ShareData có cấu trúc phức tạp
		BackgroundData: &entity.BackgroundData{
			ThemeID:   extension.BackgroundData.ThemeID,
			TextColor: extension.BackgroundData.TextColor,
		},
		QnAData: &entity.QnAData{
			Question:   extension.QnAData.Question,
			ButtonText: extension.QnAData.ButtonText,
		},
		ActivityData: &entity.ActivityData{
			Type:       extension.ActivityData.Type,
			ObjectID:   extension.ActivityData.ObjectID,
			ObjectName: extension.ActivityData.ObjectName,
		},
		LocationDetail: &entity.LocationDetail{
			Type:        "Point",                                                                          // Luôn là "Point" theo chuẩn GeoJSON
			Coordinates: []float64{extension.LocationDetail.Longitude, extension.LocationDetail.Latitude}, // Cần map chính xác nếu có dữ liệu tọa độ
			Address:     extension.LocationDetail.Address,
		},
	}
	return entity
}

func ToCreateEntityPostSettingPayload(postID string, userID string, role sharedEnums.RoleType, setting *contentEvent.SettingPayload) *entity.PostSetting {
	if setting == nil {
		return nil
	}
	convertedPostID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		// Xử lý lỗi nếu postID không hợp lệ
		return nil
	}

	postSetting := &entity.PostSetting{
		ID:        primitive.NewObjectID(),
		PostID:    convertedPostID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if setting.Schedule != nil {
		postSetting.Schedule = &entity.PostSchedule{
			IsScheduled:        true,
			PublishTime:        setting.Schedule.PublishTime,
			PublisherUserID:    userID, // UserID sẽ được gán sau khi giải mã token, không lấy từ payload
			AuthorRoleSnapshot: role,   // Vai trò sẽ được gán sau khi giải mã token, không lấy từ payload
		}
	}

	if setting.Targeting != nil {
		postSetting.Targeting = &entity.PostTargeting{
			Locations: setting.Targeting.Locations,
			AgeMin:    setting.Targeting.AgeMin,
			AgeMax:    setting.Targeting.AgeMax,
			Genders:   setting.Targeting.Genders,
			Languages: setting.Targeting.Languages,
			Interests: setting.Targeting.Interests,
		}
	}

	return postSetting
}
