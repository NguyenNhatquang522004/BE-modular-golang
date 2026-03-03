package usecase

import (
	"context"
	"errors"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeletePrivateConversationGroupUsecase struct {
	conversationRepo            IRepositoryMongodb.IConversationsRepository
	conversationParticipantRepo IRepositoryMongodb.IConversationParticipantsRepository
	calllog                     IRepositoryMongodb.ICallLogsRepository
	messageReactionRepo         IRepositoryCassandra.IMessageReactionsRepository
	messageRepo                 IRepositoryCassandra.IMessageRepository
	conversationReadRepo        IRepositoryCassandra.IConversationReadStateRepository
	pool                        IRepositoryShare.IWorkerPool
}

func NewDeletePrivateConversationGroupUsecase(conversationRepo IRepositoryMongodb.IConversationsRepository, conversationParticipantRepo IRepositoryMongodb.IConversationParticipantsRepository, calllog IRepositoryMongodb.ICallLogsRepository, messageReactionRepo IRepositoryCassandra.IMessageReactionsRepository, messageRepo IRepositoryCassandra.IMessageRepository, conversationReadRepo IRepositoryCassandra.IConversationReadStateRepository, pool IRepositoryShare.IWorkerPool) *DeletePrivateConversationGroupUsecase {
	return &DeletePrivateConversationGroupUsecase{
		conversationRepo:            conversationRepo,
		conversationParticipantRepo: conversationParticipantRepo,
		calllog:                     calllog,
		messageReactionRepo:         messageReactionRepo,
		messageRepo:                 messageRepo,
		conversationReadRepo:        conversationReadRepo,
		pool:                        pool,
	}
}
func (u *DeletePrivateConversationGroupUsecase) Execute(ctx context.Context, req *req.DeletePrivateConversationGroupRequest) (*res.FailedPrivateConversationResponse, error) {
	cursor := ""
	limit := 100
	hasnext := true
	workercout := 6
	resultChan := make(chan error, workercout)
	for hasnext == true {
		dataconversation, err := u.conversationRepo.GetConversationByGroupID(ctx, req.GroupID, cursor, limit)
		if err != nil {
			return &res.FailedPrivateConversationResponse{
				ConversationID: req.GroupID,
				ErrorMessage:   err,
			}, nil
		}
		if dataconversation.HasNext == false {
			hasnext = false
			break
		}
		var conversationIDs []primitive.ObjectID
		var conversationID2s []string
		cursor = dataconversation.NextCursor
		hasnext = true
		if dataconversation.Data == nil {
			hasnext = false
			break
		}
		if len(dataconversation.Data.([]*entity.Conversation)) == 0 {
			hasnext = false
			break
		}
		for _, item := range dataconversation.Data.([]*entity.Conversation) {
			if item.RelatedGroupID.Hex() != req.GroupID {
				hasnext = false
				break
			}
			conversationIDs = append(conversationIDs, item.ID)
			conversationID2s = append(conversationID2s, item.ID.Hex())
		}
		for i := 0; i < workercout; i++ {
			err = u.pool.Run(ctx, func() {
				switch i {
				case 0:
					_, _, err = u.conversationRepo.DeleteBulkConversations(ctx, conversationID2s)
					if err != nil {
						resultChan <- err
						return
					}
					resultChan <- nil
				case 1:
					err = u.conversationParticipantRepo.DeleteBulkConversationParticipantsByConversationIDs(ctx, conversationIDs)
					if err != nil {
						resultChan <- err
						return
					}
					resultChan <- nil
				case 2:
					err = u.calllog.DeleteBulkCallLogsByConversationIDs(ctx, conversationIDs)
					if err != nil {
						resultChan <- err
						return
					}
					resultChan <- nil
				case 3:
					err = u.messageReactionRepo.DeleteBulkReactionByConversationID(ctx, conversationID2s)
					if err != nil {
						resultChan <- err
						return
					}
					resultChan <- nil
				case 4:
					err = u.messageRepo.DeleteBulkMessagesByConversationID(ctx, conversationID2s)
					if err != nil {
						resultChan <- err
						return
					}
					resultChan <- nil
				case 5:
					err = u.conversationReadRepo.DeleteBulkConversationReadState(ctx, conversationID2s)
					if err != nil {
						resultChan <- err
						return
					}
					resultChan <- nil
				}

			})
		}
		u.pool.Wait()
		conversationID2s = nil
		conversationIDs = nil
		// Process dataconversation here
		// Update cursor and hasnext based on the response
	}
	return &res.FailedPrivateConversationResponse{
		ConversationID: req.GroupID,
		ErrorMessage:   errors.New("Failed to delete private conversation group"),
	}, nil
}
