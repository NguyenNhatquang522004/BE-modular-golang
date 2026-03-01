package mediaEvent

type ReplyStoryPayload struct {
	StoryID string `json:"story_id"`
	ReplyID int    `json:"reply_id"`
}
