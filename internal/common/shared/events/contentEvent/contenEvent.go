package contentEvent

type SharePostPayload struct {
	UserID string `json:"user_id"`
	PostID string `json:"post_id"`
}

type DeleteContentRelationTargetPayload struct {
	TargetID string `json:"target_id"`
}
