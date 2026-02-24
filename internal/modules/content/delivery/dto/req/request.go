package req

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"

type PublishPostRequest struct {
	Post          *PostReq
	PostSetting   *PostSettingReq
	PostMedia     *PostMediaReq
	PostExtension *PostExtensionReq
}
type PostIDRequest struct {
	PostID string `json:"post_id"`
}
type EditPostRequest struct {
	EditLog *PostEntityEditLogReq
}

type HideOrUnhidePostRequest struct {
	Post     *PostReq
	PostEdit *PostEntityEditLogReq
}
type DeletePostRequest struct {
	PostID PostIDRequest
	UserID string `json:"user_id"`
}

type GetPostByUserIDRequest struct {
	UserID string `json:"user_id"`
	Post   *dto.PaginationReq
}

type GetPostDataUpdateRequest struct {
	PostID PostIDRequest
}

type GetPostEditLogsByPostIDRequest struct {
	PostID PostIDRequest
}
