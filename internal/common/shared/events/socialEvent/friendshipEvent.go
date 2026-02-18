package socialEvent

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"


type FriendshipsDeletePayload struct {
	Requester_ID string `json:"user_id_1"`
	Recipient_ID string `json:"user_id_2"`
}
type FriendshipsUpdatePayload struct {
	Requester_ID string                `json:"user_id_1"`
	Recipient_ID string                `json:"user_id_2"`
	Status  enum.StatusFriendship `json:"status"`
}