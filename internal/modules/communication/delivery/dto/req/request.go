package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type CreatePrivateConversationRequest struct {
	*ConversationReq
	UserSenderFirst   *ConversationParticipantReq `json:"user_sender_first,omitempty"`
	UserReceiverFirst *ConversationParticipantReq `json:"user_receiver_first,omitempty"`
}
type UpdatePrivateConversationRequest struct {
	*ConversationReq
}
type DeletePrivateConversationRequest struct {
	ConversationID string `json:"conversation_id"`
	GroupID        string `json:"group_id,omitempty"`
}
type DeletePrivateConversationGroupRequest struct {
	GroupID string `json:"group_id"`
}
type CreateChannelGroupConversationRequest struct {
	*ConversationReq
	User []*ConversationParticipantReq `json:"user_first,omitempty"`
}
type CreateGroupConversationRequest struct {
	*ConversationReq
	User []*ConversationParticipantReq `json:"user_first,omitempty"`
}
type CreateGroupToChannelConversationRequest struct {
	ConversationChannelID string `json:"conversation_channel_id,omitempty"`
	*ConversationReq
	User []*ConversationParticipantReq `json:"user_first,omitempty"`
}
type UpdateChannelConversationRequest struct {
	ChannelID string `json:"channel_id,omitempty"`
	*ConversationReq
}
type DeleteChannelConversationRequest struct {
	ChannelID string `json:"channel_id,omitempty"`
}
type UpdateGroupConversationRequest struct {
	*ConversationReq
}
type DeleteGroupConversationRequest struct {
	ConversationID string `json:"conversation_id"`
}
type AddChannelParticipantRequest struct {
	ChannelID string `json:"channel_id,omitempty"`
	*ConversationParticipantReq
}
type AddGroupToChannelParticipantRequest struct {
	ChannelID string `json:"channel_id,omitempty"`
	GroupID   string `json:"group_id,omitempty"`
	*ConversationParticipantReq
}
type AddGroupParticipantRequest struct {
	GroupID string `json:"group_id"`
	*ConversationParticipantReq
}
type RemoveGroupParticipantRequest struct {
	GroupID string `json:"group_id"`
	UserID  string `json:"user_id"`
}
type UpdateGroupParticipantRequest struct {
	GroupID string `json:"group_id"`
	UserID  string `json:"user_id"`
	*ConversationParticipantReq
}
type MessageRequest struct {
	ConversationID string `json:"conversation_id"`
	*MessageReq
	EventType constants.EventType `json:"event_type"`
}
type MessageStateRequest struct {
	ConversationID string                   `json:"conversation_id"`
	UserID         string                   `json:"user_id"`
	Bucket         int                      `json:"bucket"`
	MessageID      string                   `json:"message_id"`
	ReactionCode   sharedEnums.ReactionCode `json:"reaction_code"`
	EventType      constants.EventType      `json:"event_type"`
}

type MessageRelyStoryRequest struct {
	StoryID string    `json:"story_id"`
	UserID  string    `json:"user_id"`
	Bucket  time.Time `json:"bucket"`
	Content string    `json:"content,omitempty"`
}
type GetCoversationListRequest struct {
}
