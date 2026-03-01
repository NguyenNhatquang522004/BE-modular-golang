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
