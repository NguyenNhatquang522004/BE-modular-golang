package contentEvent

type SharePostPayload struct {
	UserID  string  `json:"user_id"`
	PostID  string  `json:"post_id"`
	GroupID *string `json:"group_id,omitempty"` // Nếu share vào group, có thể có group_id
	PageID  *string `json:"page_id,omitempty"`  // Nếu share vào page, có thể có page_id

}

type DeleteContentRelationTargetPayload struct {
	TargetID string `json:"target_id"`
}
