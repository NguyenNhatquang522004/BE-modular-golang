package res

type CommentDetailResponse struct {
	Comment   *CommentRes
	Edit      []*CommentEditLogRes
	Reactions []*EntityReactionRes
}
type EnumReponse struct {
	Value int    `json:"value"`
	Label string `json:"label"`
}

type FailedCommentCountResponse struct {
	CommentID    string `json:"comment_id,omitempty"` // Nếu có comment_id thì đếm reply của comment đó, không có thì đếm comment của post
	ReplyCount   int    `json:"reply_count,omitempty"`
	MentionCount int    `json:"mention_count,omitempty"`
	ReportCount  int    `json:"report_count,omitempty"`
	EventType    string `json:"event_type,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}
