package mediaInContent

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"

type DeleteMediaAssetsPayload struct {
	PostID string `json:"post_id"`
	UserID string `json:"user_id"`
}

// mapper DeletePostRequest
func DeleteMediaAssetsPayloadRequestToPayload(request *req.DeletePostRequest) *DeleteMediaAssetsPayload {
	return &DeleteMediaAssetsPayload{
		PostID: request.PostID.PostID,
		UserID: request.UserID,
	}
}
