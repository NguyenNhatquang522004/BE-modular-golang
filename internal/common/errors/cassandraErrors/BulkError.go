package cassandraErrors

import "time"

type ReactionBulkError struct {
	TargetID string `json:"target_id"` // PostID/CommentID bị lỗi
	UserID   string `json:"user_id"`   // User bị lỗi
	Error    string `json:"error"`     // Chi tiết lỗi
}

type StoryViewBulkError struct {
	StoryID string `json:"story_id"` // StoryID bị lỗi
	UserID  string `json:"user_id"`  // User bị lỗi
	Error   string `json:"error"`    // Chi tiết lỗi
}
type LiveCommentBulkError struct {
	StreamID  string `json:"stream_id"`  // StreamID bị lỗi
	CommentID string `json:"comment_id"` // CommentID bị lỗi
	UserID    string `json:"user_id"`    // User bị lỗi
	Error     string `json:"error"`      // Chi tiết lỗi
}

type MessageBulkError struct {
	ConversationID string `json:"conversation_id"` // ConversationID bị lỗi
	MessageID      string `json:"message_id"`      // MessageID bị lỗi
	Error          string `json:"error"`           // Chi tiết lỗi
}

type MessageReactionBulkError struct {
	ConversationID string `json:"conversation_id"` // ConversationID bị lỗi
	MessageID      string `json:"message_id"`      // MessageID bị lỗi
	UserID         string `json:"user_id"`         // UserID bị lỗi
	Error          string `json:"error"`           // Chi tiết lỗi
}

type ConversationReadStateBulkError struct {
	ConversationID string `json:"conversation_id"` // ConversationID bị lỗi
	UserID         string `json:"user_id"`         // UserID bị lỗi
	Error          string `json:"error"`           // Chi tiết lỗi
}

type PageDailyMetricsBulkError struct {
	PageID     string    `json:"page_id"`     // PageID bị lỗi
	MetricDate time.Time `json:"metric_date"` // MetricDate bị lỗi
	Error      string    `json:"error"`       // Chi tiết lỗi
}
type NotificationBulkError struct {
	UserID         string    `json:"user_id"`         // UserID bị lỗi
	CreatedAt      time.Time `json:"created_at"`      // CreatedAt bị lỗi
	NotificationID string    `json:"notification_id"` // NotificationID bị lỗi
	Error          string    `json:"error"`           // Chi tiết lỗi
}
