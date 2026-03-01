package consumer

import (
	"context"
	"errors"
	"log"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/entity"
	"github.com/gocql/gocql"
)

type ConsumerMessage struct {
	conversationRepo        IRepositoryMongodb.IConversationsRepository
	conversationParticipant IRepositoryMongodb.IConversationParticipantsRepository
	messageStateRepo        IRepositoryCassandra.IConversationReadStateRepository
	messageReact            IRepositoryCassandra.IMessageReactionsRepository
	messageRepo             IRepositoryCassandra.IMessageRepository
	events                  events.EventBus
	pool                    IRepositoryShare.IWorkerPool
}

func NewConsumerMessage(conversationRepo IRepositoryMongodb.IConversationsRepository, conversationParticipant IRepositoryMongodb.IConversationParticipantsRepository, messageStateRepo IRepositoryCassandra.IConversationReadStateRepository, messageReact IRepositoryCassandra.IMessageReactionsRepository, messageRepo IRepositoryCassandra.IMessageRepository, events events.EventBus, pool IRepositoryShare.IWorkerPool) *ConsumerMessage {
	return &ConsumerMessage{
		conversationRepo:        conversationRepo,
		conversationParticipant: conversationParticipant,
		messageStateRepo:        messageStateRepo,
		messageReact:            messageReact,
		messageRepo:             messageRepo,
		events:                  events,
		pool:                    pool,
	}
}

func (c *ConsumerMessage) ConsumerMessage(ctx context.Context) error {
	err := c.events.Subscribe(ctx, constants.TopicMessage.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data, ok := event.Payload.(*req.MessageRequest)
		if !ok {
			return errors.New("invalid event payload")
		}
		switch data.EventType {
		case constants.Created:
			// Xử lý logic khi nhận được event tạo message
			err := c.pool.Run(ctx, func() {
				dataconversation, err := c.conversationRepo.GetConversationByID(ctx, data.ConversationID)
				if err != nil {
					log.Printf("Error fetching conversation: %v", err)
					return
				}
				if dataconversation == nil {
					log.Printf("Conversation not found: %s", data.ConversationID)
					return
				}

				err = c.conversationRepo.UpdateConversation(ctx, dataconversation)
				messageEntity := mapper.ToEntityMessage(data.MessageReq, dataconversation.ID.Hex())
				err = c.messageRepo.CreateMessage(ctx, messageEntity)
				if err != nil {
					log.Printf("Error creating message: %v", err)
					return
				}
				// Tạo ConversationReadState cho tất cả người tham gia cuộc trò chuyện
				dataconversation.LastMessage.MessageID = messageEntity.MessageID.String()
				dataconversation.LastMessage.Content = messageEntity.Content
				dataconversation.LastMessage.SenderID = messageEntity.SenderID.String()
				dataconversation.LastMessage.CreatedAt = messageEntity.CreatedAt
				dataconversation.LastMessage.Type = messageEntity.Type
				err = c.conversationRepo.UpdateConversation(ctx, dataconversation)
				if err != nil {
					log.Printf("Error updating conversation: %v", err)
					return
				}
			})
			if err != nil {
				log.Printf("Error running task in worker pool: %v", err)
			}
		case constants.Updated:
			// Xử lý logic khi nhận được event cập nhật message
			err := c.pool.Run(ctx, func() {
				dataconversation, err := c.conversationRepo.GetConversationByID(ctx, data.ConversationID)
				if err != nil {
					log.Printf("Error fetching conversation: %v", err)
					return
				}
				if dataconversation == nil {
					log.Printf("Conversation not found: %s", data.ConversationID)
					return
				}
				dataGetMessage, err := c.messageRepo.GetMessagesByConversationID(ctx, data.ConversationID, data.Bucket)
				if err != nil {
					log.Printf("Error fetching message: %v", err)
					return
				}
				if dataGetMessage == nil {
					log.Printf("Message not found: %s", data.MessageID)
					return
				}
				mapper.UpdateToEntityMessage(data.MessageReq, dataGetMessage)
				dataGetMessage.IsEdited = true // Đánh dấu là đã chỉnh sửa
				err = c.messageRepo.UpdateMessage(ctx, dataGetMessage)
				if err != nil {
					log.Printf("Error updating message: %v", err)
					return
				}
			})
			if err != nil {
				log.Printf("Error running task in worker pool: %v", err)
			}
		case constants.Deleted:
			err := c.pool.Run(ctx, func() {
				dataconversation, err := c.conversationRepo.GetConversationByID(ctx, data.ConversationID)
				if err != nil {
					log.Printf("Error fetching conversation: %v", err)
					return
				}
				if dataconversation == nil {
					log.Printf("Conversation not found: %s", data.ConversationID)
					return
				}
				dataGetMessage, err := c.messageRepo.GetMessagesByConversationID(ctx, data.ConversationID, data.Bucket)
				if err != nil {
					log.Printf("Error fetching message: %v", err)
					return
				}
				if dataGetMessage == nil {
					log.Printf("Message not found: %s", data.MessageID)
					return
				}
				dataGetMessage.IsRevoked = true // Đánh dấu là đã bị thu hồi
				err = c.messageRepo.UpdateMessage(ctx, dataGetMessage)
				if err != nil {
					log.Printf("Error updating message: %v", err)
					return
				}
			})
			if err != nil {
				log.Printf("Error running task in worker pool: %v", err)
			}

		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func (c *ConsumerMessage) FailedMessage(ctx context.Context) error {
	return nil
}

func (c *ConsumerMessage) ConsumerStateMessage(ctx context.Context) error {
	err := c.events.Subscribe(ctx, constants.TopicStateMessage.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data, ok := event.Payload.(*req.MessageStateRequest)
		if !ok {
			return errors.New("invalid event payload")
		}
		switch data.EventType {
		case constants.Created:
			dataconversation, err := c.conversationRepo.GetConversationByID(ctx, data.ConversationID)
			if err != nil {
				log.Printf("Error fetching conversation: %v", err)
				return nil
			}
			if dataconversation == nil {
				log.Printf("Conversation not found: %s", data.ConversationID)
				return nil
			}
			dataParticipant, err := c.conversationParticipant.GetConversationParticipant(ctx, data.ConversationID, data.UserID)
			if err != nil {
				log.Printf("Error fetching conversation participant: %v", err)
				return nil
			}
			if dataParticipant == nil {
				log.Printf("Conversation participant not found: ConversationID=%s, UserID=%s", data.ConversationID, data.UserID)
				return nil
			}
			dataParticipant.LastSeenMessageID = data.MessageID
			err = c.conversationParticipant.UpdateConversationParticipant(ctx, dataParticipant)
			if err != nil {
				log.Printf("Error updating conversation participant: %v", err)
				return nil
			}
			convertcql, err := gocql.ParseUUID(data.MessageID)
			if err != nil {
				log.Printf("Error parsing MessageID to UUID: %v", err)
				return nil
			}
			convertuseridcql, err := gocql.ParseUUID(data.UserID)
			if err != nil {
				log.Printf("Error parsing UserID to UUID: %v", err)
				return nil
			}
			entityReadState := &entity.ConversationReadState{
				ConversationID:    data.ConversationID,
				UserID:            convertuseridcql,
				LastReadMessageID: convertcql,
				LastReadAt:        data.LastReadAt,
			}
			err = c.messageStateRepo.UpsertConversationReadState(ctx, entityReadState)
			if err != nil {
				log.Printf("Error updating conversation read state: %v", err)
				return nil
			}
		case constants.Updated:
		case constants.Deleted:
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
func (c *ConsumerMessage) FailedStateMessage(ctx context.Context) error {
	return nil
}

func (c *ConsumerMessage) ConsumerReactMessage(ctx context.Context) error {
	err := c.events.Subscribe(ctx, constants.TopicReactMessage.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data, ok := event.Payload.(*req.ReactMessageRequest)
		if !ok {
			return errors.New("invalid event payload")
		}
		switch data.EventType {
		case constants.Created:
			entity := mapper.ToEntityMessageReaction(data.MessageReactionReq)
			err := c.messageReact.CreateReaction(ctx, entity)
			if err != nil {
				log.Printf("Error creating message reaction: %v", err)
				return nil
			}
		case constants.Updated:
			convertConversationID, err := gocql.ParseUUID(data.ConversationID)
			if err != nil {
				log.Printf("Error parsing ConversationID to UUID: %v", err)
				return nil
			}
			convertUserID, err := gocql.ParseUUID(data.UserID.String())
			if err != nil {
				log.Printf("Error parsing UserID to UUID: %v", err)
				return nil
			}
			data, err := c.messageReact.GetReactionByUser(ctx, convertConversationID.String(), data.MessageID.String(), convertUserID.String())
			if err != nil {
				log.Printf("Error fetching message reaction: %v", err)
				return nil
			}
			data.ReactionCode = data.ReactionCode
			err = c.messageReact.UpdateReaction(ctx, data)
			if err != nil {
				log.Printf("Error updating message reaction: %v", err)
				return nil
			}
			// Xử lý logic khi nhận được event cập nhật reaction
			// Tương tự như phần Created nhưng thay vì tạo mới thì cập nhật reaction hiện có
			// Có thể cần thêm logic để kiểm tra nếu reaction đã tồn tại thì cập nhật, nếu chưa tồn tại thì tạo mới
		case constants.Deleted:
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func (c *ConsumerMessage) FailedReactMessage(ctx context.Context) error
