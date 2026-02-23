package req

type PublishPostRequest struct {
	Post          *PostReq
	PostSetting   *PostSettingReq
	PostMedia     *PostMediaReq
	PostExtension *PostExtensionReq
}

type EditPostRequest struct {
}

type HideOrUnhidePostRequest struct {
	Post     *PostReq
	PostEdit *PostEntityEditLogReq
}
type DeletePostRequest struct {
}

type GetPostByUserIDRequest struct{}

type GetPostDetailRequest struct {
}
