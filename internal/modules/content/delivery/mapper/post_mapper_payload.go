package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/contentEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToCreateEntityPostPayload(id primitive.ObjectID, r *contentEvent.CreatePostPayload) *entity.Post {
	if r == nil {
		return nil
	}
	post := &entity.Post{
		ID:      id,
		UserID:  r.UserID, // UserID sẽ được gán sau khi giải mã token, không lấy từ payload
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

func UpdateEntityPostFromPayload(post *entity.Post, r *contentEvent.UpdatePostPayload) {
	if r.Content != nil {
		post.Content = *r.Content
	}
	if r.Privacy != nil {
		post.Privacy.Scope = *r.Privacy.Scope
		post.Privacy.AllowComment = *r.Privacy.AllowComment
		post.Privacy.AllowShare = *r.Privacy.AllowShare
	}
	if r.Hashtags != nil {
		post.Hashtags = *r.Hashtags
	}
	if r.Mentions != nil {
		post.Mentions = *r.Mentions
	}
	post.IsEdited = true
	post.UpdatedAt = time.Now()
}

func UpdateEntityPostMediaFromPayload(postMedia *entity.PostMedia, r *contentEvent.UpdatePostPayload) {
	var taggedUsers []entity.TaggedUser
	for _, item := range *r.Media {
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
	for index, item := range *r.Media {
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
	postMedia.Items = items
	postMedia.UpdatedAt = time.Now()
}

func UpdateEntityPostExtensionFromPayload(postExtension *entity.PostExtension, r *contentEvent.UpdatePostPayload) {
	if r.Extension != nil {
		postExtension.BackgroundData = &entity.BackgroundData{
			ThemeID:   r.Extension.BackgroundData.ThemeID,
			TextColor: r.Extension.BackgroundData.TextColor,
		}
		postExtension.QnAData = &entity.QnAData{
			Question:   r.Extension.QnAData.Question,
			ButtonText: r.Extension.QnAData.ButtonText,
		}
		postExtension.ActivityData = &entity.ActivityData{
			Type:       r.Extension.ActivityData.Type,
			ObjectID:   r.Extension.ActivityData.ObjectID,
			ObjectName: r.Extension.ActivityData.ObjectName,
		}
		postExtension.LocationDetail = &entity.LocationDetail{
			Type:        "Point",                                                                              // Luôn là "Point" theo chuẩn GeoJSON
			Coordinates: []float64{r.Extension.LocationDetail.Longitude, r.Extension.LocationDetail.Latitude}, // Cần map chính xác nếu có dữ liệu tọa độ
			Address:     r.Extension.LocationDetail.Address,
		}
		postExtension.UpdatedAt = time.Now()
	}
}
func UpdateEntityPostSettingFromPayload(postSetting *entity.PostSetting, r *contentEvent.UpdatePostPayload) {
	if r.Setting != nil {
		if r.Setting.Schedule != nil {
			postSetting.Schedule = &entity.PostSchedule{
				IsScheduled:        true,
				PublishTime:        r.Setting.Schedule.PublishTime,
				PublisherUserID:    postSetting.Schedule.PublisherUserID,    // Giữ nguyên PublisherUserID
				AuthorRoleSnapshot: postSetting.Schedule.AuthorRoleSnapshot, // Giữ nguyên AuthorRoleSnapshot
			}
		}
		if r.Setting.Targeting != nil {
			postSetting.Targeting = &entity.PostTargeting{
				Locations: r.Setting.Targeting.Locations,
				AgeMin:    r.Setting.Targeting.AgeMin,
				AgeMax:    r.Setting.Targeting.AgeMax,
				Languages: r.Setting.Targeting.Languages,
				Interests: r.Setting.Targeting.Interests,
			}
		}
		postSetting.UpdatedAt = time.Now()
	}
}
