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
