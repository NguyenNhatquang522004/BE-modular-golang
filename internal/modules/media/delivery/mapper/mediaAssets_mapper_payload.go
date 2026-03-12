package mapper

import (
	"errors"
	"fmt"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToEntityCreateMediaAssetsPayload(data *mediaEvent.CreateMediaAssetsPayload) ([]*entity.MediaAsset, error) {
	var mediaAsset []*entity.MediaAsset
	if data == nil {
		return nil, nil
	}
	if data.UserID == "" {
		return nil, errors.New("userID is required")
	}
	for _, item := range data.Items {
		mediaAssetItem := &entity.MediaAsset{}
		ID := primitive.NewObjectID()
		if item.MediaID != "" {
			converid, err := primitive.ObjectIDFromHex(item.MediaID)
			if err != nil {
				return nil, errors.New("invalid mediaID: " + err.Error())
			}
			ID = converid
		}
		mediaAssetItem.ID = ID
		mediaAssetItem.UserID = data.UserID
		if item.PostID != "" {
			postID, err := primitive.ObjectIDFromHex(item.PostID)
			if err != nil {
				return nil, errors.New("invalid postID: " + err.Error())
			}
			mediaAssetItem.PostID = postID
		}
		if item.StoryID != "" {
			storyID, err := primitive.ObjectIDFromHex(item.StoryID)
			if err != nil {
				return nil, errors.New("invalid storyID: " + err.Error())
			}
			mediaAssetItem.StoryID = storyID
		}
		if item.AlbumID != "" {
			albumID, err := primitive.ObjectIDFromHex(item.AlbumID)
			if err != nil {
				return nil, errors.New("invalid albumID: " + err.Error())
			}
			mediaAssetItem.AlbumID = albumID
		}
		if item.CommentID != "" {
			commentID, err := primitive.ObjectIDFromHex(item.CommentID)
			if err != nil {
				return nil, errors.New("invalid commentID: " + err.Error())
			}
			mediaAssetItem.CommentID = commentID
		}
		if item.ReelID != "" {
			reelID, err := primitive.ObjectIDFromHex(item.ReelID)
			if err != nil {
				return nil, errors.New("invalid reelID: " + err.Error())
			}
			mediaAssetItem.ReelID = reelID
		}
		if item.GroupID != "" {
			groupID, err := primitive.ObjectIDFromHex(item.GroupID)
			if err != nil {
				return nil, errors.New("invalid groupID: " + err.Error())
			}
			mediaAssetItem.GroupID = groupID
		}
		if item.PageID != "" {
			pageID, err := primitive.ObjectIDFromHex(item.PageID)
			if err != nil {
				return nil, errors.New("invalid pageID: " + err.Error())
			}
			mediaAssetItem.PageID = pageID
		}
		mediaAssetItem.AssetType = item.MediaType
		mediaAssetItem.StorageFileID = item.URL
		mediaAssetItem.ThumbnailURL = item.URL // Tạm thời set thumbnail = original, sau này có thể xử lý riêng
		mediaAssetItem.Order = item.Order
		mediaAssetItem.Hashtags = item.Hashtags
		mediaAssetItem.Privacy = &entity.MediaPrivacy{
			Level:            sharedEnums.ScopeFriends, // Mặc định public, sau này có thể mở rộng
			InheritFromAlbum: false,
		}
		mediaAssetItem.Metadata = entity.MediaMetadata{
			Width:     item.Metadata.Width,
			Height:    item.Metadata.Height,
			Duration:  item.Metadata.Duration,
			SizeBytes: item.Metadata.SizeBytes,
			MimeType:  item.Metadata.MimeType,
		}
		for _, taggedUser := range item.TaggedUsers {
			mediaAssetItem.TaggedUsers = append(mediaAssetItem.TaggedUsers, entity.MediaTag{
				UserID: taggedUser.UserID,
				Name:   taggedUser.Name,
				Position: entity.TagPosition{
					X: taggedUser.X,
					Y: taggedUser.Y,
				},
			})
		}
		mediaAsset = append(mediaAsset, mediaAssetItem)
	}
	return mediaAsset, nil
}

func UpdateEntityMediaAssetsFromPayload(mediaAsset *entity.MediaAsset, data *mediaEvent.UpdateMediaAssetsPayload) error {
	if data.AlbumID != "" {
		albumID, err := primitive.ObjectIDFromHex(data.AlbumID)
		if err != nil {
			return fmt.Errorf("invalid albumID: %w", err)
		}
		mediaAsset.AlbumID = albumID
	}
	if data.URL != nil {
		mediaAsset.StorageFileID = *data.URL
	}
	if data.ThumbnailURL != nil {
		mediaAsset.ThumbnailURL = *data.ThumbnailURL
	}
	if data.Order != nil {
		mediaAsset.Order = *data.Order
	}
	if data.Hashtags != nil {
		mediaAsset.Hashtags = *data.Hashtags
	}
	if data.Metadata != nil {

		mediaAsset.Metadata = entity.MediaMetadata{
			Width:     data.Metadata.Width,
			Height:    data.Metadata.Height,
			Duration:  data.Metadata.Duration,
			SizeBytes: data.Metadata.SizeBytes,
			MimeType:  data.Metadata.MimeType,
		}
	}
	if data.TaggedUsers != nil {
		var taggedUsers []entity.MediaTag
		for _, taggedUser := range *data.TaggedUsers {
			taggedUsers = append(taggedUsers, entity.MediaTag{
				UserID: taggedUser.UserID,
				Name:   taggedUser.Name,
				Position: entity.TagPosition{
					X: taggedUser.X,
					Y: taggedUser.Y,
				},
			})
		}
		mediaAsset.TaggedUsers = taggedUsers
	}
	if data.Privacy != nil {
		mediaAsset.Privacy = &entity.MediaPrivacy{
			Level:            data.Privacy.Level,
			InheritFromAlbum: data.Privacy.InheritFromAlbum,
		}
	}
	return nil
}
