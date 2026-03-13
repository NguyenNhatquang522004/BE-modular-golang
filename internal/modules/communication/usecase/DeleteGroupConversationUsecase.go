package usecase

import (
	"context"
	"errors"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryMongodb"
)

type DeleteGroupConversationUsecase struct {
	conversationRepo            IRepositoryMongodb.IConversationsRepository
	conversationParticipantRepo IRepositoryMongodb.IConversationParticipantsRepository
	calllog                     IRepositoryMongodb.ICallLogsRepository
	messageReactionRepo         IRepositoryCassandra.IMessageReactionsRepository
	messageRepo                 IRepositoryCassandra.IMessageRepository
	conversationReadRepo        IRepositoryCassandra.IConversationReadStateRepository
	pool                        IRepositoryShare.IWorkerPool
}

func NewDeleteGroupConversationUsecase(conversationRepo IRepositoryMongodb.IConversationsRepository, conversationParticipantRepo IRepositoryMongodb.IConversationParticipantsRepository, calllog IRepositoryMongodb.ICallLogsRepository, messageReactionRepo IRepositoryCassandra.IMessageReactionsRepository, messageRepo IRepositoryCassandra.IMessageRepository, conversationReadRepo IRepositoryCassandra.IConversationReadStateRepository, pool IRepositoryShare.IWorkerPool) *DeleteGroupConversationUsecase {
	return &DeleteGroupConversationUsecase{
		conversationRepo:            conversationRepo,
		conversationParticipantRepo: conversationParticipantRepo,
		calllog:                     calllog,
		messageReactionRepo:         messageReactionRepo,
		messageRepo:                 messageRepo,
		conversationReadRepo:        conversationReadRepo,
		pool:                        pool,
	}
}

func (u *DeleteGroupConversationUsecase) Execute(ctx context.Context, req *req.DeleteGroupConversationRequest) (*res.FailedChannelGroupConversationResponse, error) {
	dataConversation, err := u.conversationRepo.GetConversationByID(ctx, req.ConversationID)
	if err != nil {
		return &res.FailedChannelGroupConversationResponse{
			ChannelID:      "unknown",
			ConversationID: req.ConversationID,
			UserID:         "unknown",
			ErrorMessage:   errors.New("Failed to get conversation"),
		}, err
	}
	if dataConversation == nil {
		return &res.FailedChannelGroupConversationResponse{
			ChannelID:      "unknown",
			ConversationID: req.ConversationID,
			UserID:         "unknown",
			ErrorMessage:   errors.New("Conversation not found"),
		}, nil
	}
	if dataConversation.Type != sharedEnums.TypeGroup {
		return &res.FailedChannelGroupConversationResponse{
			ChannelID:      "unknown",
			ConversationID: req.ConversationID,
			UserID:         "unknown",
			ErrorMessage:   errors.New("Conversation is not a group conversation"),
		}, nil
	}
	workercount := 6
	resultChan := make(chan *res.FailedChannelGroupConversationResponse, workercount)
	for i := 0; i < workercount; i++ {
		switch i {
		case 0:
			err := u.conversationParticipantRepo.DeleteConversationParticipantByConversationID(ctx, dataConversation.ID.Hex())
			if err != nil {
				resultChan <- &res.FailedChannelGroupConversationResponse{
					ChannelID:      "unknown",
					ConversationID: dataConversation.ID.Hex(),
					UserID:         "unknown",
					ErrorMessage:   errors.New("Failed to delete conversation participants"),
				}

			}
			resultChan <- &res.FailedChannelGroupConversationResponse{
				ChannelID:      "unknown",
				ConversationID: dataConversation.ID.Hex(),
				UserID:         "unknown",
				ErrorMessage:   nil,
			}
		case 1:
			err := u.calllog.DeleteCallLogByConversationID(ctx, dataConversation.ID.Hex())
			if err != nil {
				resultChan <- &res.FailedChannelGroupConversationResponse{
					ChannelID:      "unknown",
					ConversationID: dataConversation.ID.Hex(),
					UserID:         "unknown",
					ErrorMessage:   errors.New("Failed to delete call logs"),
				}
			} else {
				resultChan <- &res.FailedChannelGroupConversationResponse{
					ChannelID:      "unknown",
					ConversationID: dataConversation.ID.Hex(),
					UserID:         "unknown",
					ErrorMessage:   nil,
				}
			}
		case 2:
			err := u.messageReactionRepo.DeleteReactionByConversationID(ctx, dataConversation.ID.Hex())
			if err != nil {
				resultChan <- &res.FailedChannelGroupConversationResponse{
					ChannelID:      "unknown",
					ConversationID: dataConversation.ID.Hex(),
					UserID:         "unknown",
					ErrorMessage:   errors.New("Failed to delete message reactions"),
				}
			} else {
				resultChan <- &res.FailedChannelGroupConversationResponse{
					ChannelID:      "unknown",
					ConversationID: dataConversation.ID.Hex(),
					UserID:         "unknown",
					ErrorMessage:   nil,
				}
			}
		case 3:
			err := u.messageRepo.DeleteMessagesByConversationID(ctx, dataConversation.ID.Hex())
			if err != nil {
				resultChan <- &res.FailedChannelGroupConversationResponse{
					ChannelID:      "unknown",
					ConversationID: dataConversation.ID.Hex(),
					UserID:         "unknown",
					ErrorMessage:   errors.New("Failed to delete messages"),
				}
			} else {
				resultChan <- &res.FailedChannelGroupConversationResponse{
					ChannelID:      "unknown",
					ConversationID: dataConversation.ID.Hex(),
					UserID:         "unknown",
					ErrorMessage:   nil,
				}
			}
		case 4:
			err := u.conversationReadRepo.DeleteConversationReadStatesByConversationID(ctx, dataConversation.ID.Hex())
			if err != nil {
				resultChan <- &res.FailedChannelGroupConversationResponse{
					ChannelID:      "unknown",
					ConversationID: dataConversation.ID.Hex(),
					UserID:         "unknown",
					ErrorMessage:   errors.New("Failed to delete conversation read states"),
				}
			} else {
				resultChan <- &res.FailedChannelGroupConversationResponse{
					ChannelID:      "unknown",
					ConversationID: dataConversation.ID.Hex(),
					UserID:         "unknown",
					ErrorMessage:   nil,
				}
			}
		case 5:
			err := u.conversationRepo.DeleteConversation(ctx, dataConversation.ID.Hex())
			if err != nil {
				resultChan <- &res.FailedChannelGroupConversationResponse{
					ChannelID:      "unknown",
					ConversationID: dataConversation.ID.Hex(),
					UserID:         "unknown",
					ErrorMessage:   errors.New("Failed to delete conversation"),
				}
			} else {
				resultChan <- &res.FailedChannelGroupConversationResponse{
					ChannelID:      "unknown",
					ConversationID: dataConversation.ID.Hex(),
					UserID:         "unknown",
					ErrorMessage:   nil,
				}
			}
		}

	}

	return &res.FailedChannelGroupConversationResponse{
		ChannelID:      "unknown",
		ConversationID: req.ConversationID,
		UserID:         "unknown",
		ErrorMessage:   nil,
	}, nil
}
