package mapper

import (
	"errors"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/interactionEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MapCreatedCommentPayloadToEntity(data *interactionEvent.CreatedCommentPayload) (*entity.Comment, error) {
	var entitya = &entity.Comment{}
	ID := primitive.NewObjectID() // Tạo ID mới cho comment
	if data.TargetID != "" {
		return nil, errors.New("TargetID is required and must be a valid ObjectID string")
	}
	if data.UserID == "" {
		return nil, errors.New("UserID is required")
	}
	if data.ID != "" {
		convertID, err := primitive.ObjectIDFromHex(data.ID)
		if err != nil {
			return nil, errors.New("Invalid ID format, must be a valid ObjectID string")
		}
		ID = convertID
	}
	entitya.ID = ID
	targetObjectID, err := primitive.ObjectIDFromHex(data.TargetID)
	if err != nil {
		return nil, errors.New("Invalid TargetID format, must be a valid ObjectID string")
	}
	entitya.TargetID = targetObjectID
	entitya.UserID = data.UserID
	entitya.Content = data.Content
	if data.AssetID != nil {
		assetObjectID, err := primitive.ObjectIDFromHex(*data.AssetID)
		if err != nil {
			return nil, errors.New("Invalid AssetID format, must be a valid ObjectID string")
		}
		entitya.AssetID = &assetObjectID
	}
	if data.ParentCommentID != nil {
		parentCommentObjectID, err := primitive.ObjectIDFromHex(*data.ParentCommentID)
		if err != nil {
			return nil, errors.New("Invalid ParentCommentID format, must be a valid ObjectID string")
		}
		entitya.ParentCommentID = &parentCommentObjectID
	}
	if data.RootCommentID != nil {
		rootCommentObjectID, err := primitive.ObjectIDFromHex(*data.RootCommentID)
		if err != nil {
			return nil, errors.New("Invalid RootCommentID format, must be a valid ObjectID string")
		}
		entitya.RootCommentID = &rootCommentObjectID
	}
	if data.Mentions != nil && len(*data.Mentions) > 0 {
		entitya.Mentions = *data.Mentions
	}
	if data.Media != nil {
		if data.AssetID == nil {
			return nil, errors.New("AssetID is required when Media is provided")
		}
		assetid, err := primitive.ObjectIDFromHex(*data.AssetID)
		if err != nil {
			return nil, errors.New("Invalid AssetID format, must be a valid ObjectID string")
		}
		entitya.AssetID = &assetid
		entitya.Media = &entity.CommentMedia{
			Type: data.Media.Type,
			URL:  data.Media.URL,
			DisplayMeta: entity.DisplayMeta{
				Width:     data.Media.DisplayMeta.Width,
				Height:    data.Media.DisplayMeta.Height,
				Duration:  data.Media.DisplayMeta.Duration,
				SizeBytes: data.Media.DisplayMeta.SizeBytes,
				MimeType:  data.Media.DisplayMeta.MimeType,
			},
		}
	}
	entitya.CreatedAt = data.CreatedAt
	return entitya, nil
}

func MapUpdatedCommentPayloadToEntity(data *interactionEvent.UpdatedCommentPayload, entitya *entity.Comment) (*entity.Comment, error) {
	entitya.Content = data.Content
	entitya.IsEdited = true
	entitya.LastEditedAt = &data.EditedAt

	if data.Mentions != nil {
		entitya.Mentions = *data.Mentions
	}

	if data.Media != nil {
		entitya.Media = &entity.CommentMedia{
			Type: data.Media.Type,
			URL:  data.Media.URL,
			DisplayMeta: entity.DisplayMeta{
				Width:     data.Media.DisplayMeta.Width,
				Height:    data.Media.DisplayMeta.Height,
				Duration:  data.Media.DisplayMeta.Duration,
				SizeBytes: data.Media.DisplayMeta.SizeBytes,
				MimeType:  data.Media.DisplayMeta.MimeType,
			},
		}
	}
	if data.HiddenMetadata != nil {
		entitya.HiddenMetadata = &entity.HiddenMetadata{
			IsHidden:       data.HiddenMetadata.IsHidden,
			HiddenAt:       data.HiddenMetadata.HiddenAt,
			HiddenByUserID: data.HiddenMetadata.HiddenByUserID,
			Reason:         data.HiddenMetadata.Reason,
			IsGhostBanned:  data.HiddenMetadata.IsGhostBanned,
		}
	}
	if data.DeletedMetadata != nil {
		entitya.DeletedMetadata = &entity.DeletedMetadata{
			DeletedAt:       data.DeletedMetadata.DeletedAt,
			DeletedByUserID: data.DeletedMetadata.DeletedByUserID,
		}
	}
	return entitya, nil
}
