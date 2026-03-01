package req

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
type CreateMessageRequest struct {
}
type DeleteMessageRequest struct {
}
type UpdateMessageRequest struct {
}
type ReactMessageRequest struct {
}
type ReadMessageRequest struct {
}
type UnreadMessageRequest struct {
}
type MessageReplyStoryRequest struct {
}
type GetCoversationListRequest struct {
}
