package res

type FailedPrivateConversationResponse struct {
	ConversationID      string `json:"conversation_id"`
	UserIDSenderFirst   string `json:"user_sender_first,omitempty"`
	UserIDReceiverFirst string `json:"user_receiver_first,omitempty"`
	ErrorMessage        error  `json:"error_message"`
}
type FailedChannelGroupConversationResponse struct {
	ChannelID      string `json:"channel_id"`
	ConversationID string `json:"conversation_id"`
	UserID         string `json:"user_id,omitempty"`
	ErrorMessage   error  `json:"error_message"`
}

type FailedParticipantResponse struct {
	ConversationID string `json:"conversation_id"`
	UserID         string `json:"user_id,omitempty"`
	ErrorMessage   error  `json:"error_message"`
}

type FailedMessageResponse struct {
	ConversationID string `json:"conversation_id"`
	UserID         string `json:"user_id,omitempty"`
	MessageID      string `json:"message_id,omitempty"`
	ErrorMessage   error  `json:"error_message"`
}

type FailedReactMessageResponse struct {
	ConversationID string `json:"conversation_id"`
	MessageID      string `json:"message_id"`
	UserID         string `json:"user_id,omitempty"`
	ReactionCode   string `json:"reaction_code,omitempty"`
	ErrorMessage   error  `json:"error_message"`
}
type FailedReplyStoryResponse struct {
	TargetID     string `json:"target_id"`
	UserID       string `json:"user_id,omitempty"`
	Content      string `json:"content,omitempty"`
	ErrorMessage error  `json:"error_message"`
}
