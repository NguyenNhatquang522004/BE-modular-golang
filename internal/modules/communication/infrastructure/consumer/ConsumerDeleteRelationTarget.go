package consumer

import (
	"context"
	"errors"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communicationEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConsumerDeleteRelationTarget struct {
	conversationRepo            IRepositoryMongodb.IConversationsRepository
	conversationParticipantRepo IRepositoryMongodb.IConversationParticipantsRepository
	calllog                     IRepositoryMongodb.ICallLogsRepository
	messageReactionRepo         IRepositoryCassandra.IMessageReactionsRepository
	messageRepo                 IRepositoryCassandra.IMessageRepository
	conversationReadRepo        IRepositoryCassandra.IConversationReadStateRepository
	pool                        IRepositoryShare.IWorkerPool
	event                       events.EventBus
}

func NewConsumerDeleteRelationTarget(conversationRepo IRepositoryMongodb.IConversationsRepository, conversationParticipantRepo IRepositoryMongodb.IConversationParticipantsRepository, calllog IRepositoryMongodb.ICallLogsRepository, messageReactionRepo IRepositoryCassandra.IMessageReactionsRepository, messageRepo IRepositoryCassandra.IMessageRepository, conversationReadRepo IRepositoryCassandra.IConversationReadStateRepository, pool IRepositoryShare.IWorkerPool, event events.EventBus) *ConsumerDeleteRelationTarget {
	return &ConsumerDeleteRelationTarget{
		conversationRepo:            conversationRepo,
		conversationParticipantRepo: conversationParticipantRepo,
		calllog:                     calllog,
		messageReactionRepo:         messageReactionRepo,
		messageRepo:                 messageRepo,
		conversationReadRepo:        conversationReadRepo,
		pool:                        pool,
		event:                       event,
	}
}
func (c *ConsumerDeleteRelationTarget) ConsumerDeleteRelationTarget(ctx context.Context) error {
	err := c.event.Subscribe(ctx, constants.TopicDeleteRelationTarget.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data, ok := event.Payload.(*communicationEvent.DeletePrivateConversationGroupPayload)
		if !ok {
			return errors.New("invalid event payload")
		}
		req := data
		cursor := ""
		limit := 100
		hasnext := true
		workercout := 6
		resultChan := make(chan error, workercout)
		for hasnext == true {
			dataconversation, err := c.conversationRepo.GetConversationByGroupID(ctx, req.TargetID, cursor, limit)
			if err != nil {
				return err
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
				if item.RelatedGroupID.Hex() != req.TargetID {
					hasnext = false
					break
				}
				conversationIDs = append(conversationIDs, item.ID)
				conversationID2s = append(conversationID2s, item.ID.Hex())
			}
			for i := 0; i < workercout; i++ {
				err = c.pool.Run(ctx, func() {
					switch i {
					case 0:
						_, _, err = c.conversationRepo.DeleteBulkConversations(ctx, conversationID2s)
						if err != nil {
							resultChan <- err
							return
						}
						resultChan <- nil
					case 1:
						err = c.conversationParticipantRepo.DeleteBulkConversationParticipantsByConversationIDs(ctx, conversationIDs)
						if err != nil {
							resultChan <- err
							return
						}
						resultChan <- nil
					case 2:
						err = c.calllog.DeleteBulkCallLogsByConversationIDs(ctx, conversationIDs)
						if err != nil {
							resultChan <- err
							return
						}
						resultChan <- nil
					case 3:
						err = c.messageReactionRepo.DeleteBulkReactionByConversationID(ctx, conversationID2s)
						if err != nil {
							resultChan <- err
							return
						}
						resultChan <- nil
					case 4:
						err = c.messageRepo.DeleteBulkMessagesByConversationID(ctx, conversationID2s)
						if err != nil {
							resultChan <- err
							return
						}
						resultChan <- nil
					case 5:
						err = c.conversationReadRepo.DeleteBulkConversationReadState(ctx, conversationID2s)
						if err != nil {
							resultChan <- err
							return
						}
						resultChan <- nil
					}

				})
			}
			c.pool.Wait()
			conversationID2s = nil
			conversationIDs = nil
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *ConsumerDeleteRelationTarget) FailedDeleteRelationTarget(ctx context.Context) error
