package req

import "time"

type FriendShipRequest struct {
	ID               string    `json:"id"`
	Requester_UserID string    `json:"requester_user_id"`
	Requested_UserID string    `json:"requested_user_id"`
	Created_At       time.Time `json:"created_at"`
	Updated_At       time.Time `json:"updated_at"`
	EventType        string    `json:"event_type"`
}
