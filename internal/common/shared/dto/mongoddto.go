package dto

type BulkError struct {
	ID     string
	Reason string
}

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
