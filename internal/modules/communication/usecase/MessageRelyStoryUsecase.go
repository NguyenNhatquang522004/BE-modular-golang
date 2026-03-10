package usecase

import (
	"context"
	"errors"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communicationEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/pkg/pb/v1"
	"github.com/gocql/gocql"
)

type MessageRelyStoryUsecase struct {
	pb               pb.MediaServiceClient
	conversationRepo IRepositoryMongodb.IConversationsRepository
	events           events.EventBus
}

func NewMessageRelyStoryUsecase(pb pb.MediaServiceClient, conversationRepo IRepositoryMongodb.IConversationsRepository, events events.EventBus) *MessageRelyStoryUsecase {
	return &MessageRelyStoryUsecase{
		pb:               pb,
		conversationRepo: conversationRepo,
		events:           events,
	}
}

func (u *MessageRelyStoryUsecase) Execute(ctx context.Context, reqa *req.MessageRelyStoryRequest) (*res.FailedReplyStoryResponse, error) {
	// Implement the logic to handle the message rely story request here
	// This may involve validating the request, processing the data, and returning an appropriate response
	userIDResp, err := u.pb.GetUserIDbyStoryID(ctx, &pb.GetUserIDbyStoryIDRequest{
		StoryId: reqa.StoryID,
	})
	if err != nil {
		return &res.FailedReplyStoryResponse{
			TargetID:     reqa.StoryID,
			UserID:       reqa.UserID,
			Content:      reqa.Content,
			ErrorMessage: errors.New("failed to get user ID by story ID"),
		}, err
	}
	datacheck, err := u.conversationRepo.CheckConversationExists(ctx, reqa.UserID, userIDResp.UserId)
	if err != nil {
		return &res.FailedReplyStoryResponse{
			TargetID:     reqa.StoryID,
			UserID:       reqa.UserID,
			Content:      reqa.Content,
			ErrorMessage: errors.New("failed to check conversation existence"),
		}, err
	}
	if datacheck == nil {
		return &res.FailedReplyStoryResponse{
			TargetID:     reqa.StoryID,
			UserID:       reqa.UserID,
			Content:      reqa.Content,
			ErrorMessage: errors.New("conversation does not exist"),
		}, nil
	}
	convertstoryID, ok := gocql.ParseUUID(reqa.StoryID)
	if ok != nil {
		return &res.FailedReplyStoryResponse{
			TargetID:     reqa.StoryID,
			UserID:       reqa.UserID,
			Content:      reqa.Content,
			ErrorMessage: errors.New("invalid story ID format"),
		}, nil
	}
	convertUserID, ok := gocql.ParseUUID(reqa.UserID)
	if ok != nil {
		return &res.FailedReplyStoryResponse{
			TargetID:     reqa.StoryID,
			UserID:       reqa.UserID,
			Content:      reqa.Content,
			ErrorMessage: errors.New("invalid user ID format"),
		}, nil
	}
	MessageRequest := &communicationEvent.MessagePayload{
		ConversationID: datacheck.ID.Hex(),
		Bucket:         mapper.GenerateBucket(reqa.Bucket),
		MessageID:      gocql.TimeUUID(),
		SenderID:       convertUserID,
		Type:           sharedEnums.MediaTypeText,
		Content:        reqa.Content,
		Attachments:    nil,
		IsEdited:       false,
		StoryRefID:     &convertstoryID,
		IsRevoked:      false,
		CreatedAt:      reqa.Bucket,
		EventType:      constants.Created,
	}
	err = u.events.Publish(ctx, constants.TopicMessage.String(), datacheck.ID.Hex(), constants.Created.String(), MessageRequest)
	if err != nil {
		return &res.FailedReplyStoryResponse{
			TargetID:     reqa.StoryID,
			UserID:       reqa.UserID,
			Content:      reqa.Content,
			ErrorMessage: errors.New("failed to publish message rely story event"),
		}, err
	}
	if err != nil {
		return &res.FailedReplyStoryResponse{
			TargetID:     reqa.StoryID,
			UserID:       reqa.UserID,
			Content:      reqa.Content,
			ErrorMessage: errors.New("failed to publish reply story event"),
		}, err
	}
	return &res.FailedReplyStoryResponse{
		TargetID:     reqa.StoryID,
		UserID:       reqa.UserID,
		Content:      reqa.Content,
		ErrorMessage: nil,
	}, nil
}
