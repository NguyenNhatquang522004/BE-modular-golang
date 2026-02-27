package res

type GetPostByUserIDReponse struct {
	PostRes          []*PostRes
	PostMediaRes     []*PostMediaRes
	PostSettingRes   []*PostSettingRes
	PostExtensionRes []*PostExtensionRes
	PostEditlog      []*PostEntityEditLogRes
}

type EnumReponse struct {
	Value int    `json:"value"`
	Label string `json:"label"`
}

type FailSharePostResponse struct {
	UserID  string `json:"user_id"`
	PostID  string `json:"post_id"`
	Message string `json:"message"`
}
