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