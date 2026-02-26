package res

type CommentDetailResponse struct {
	Comment   *CommentRes
	Edit      []*CommentEditLogRes
	Reactions []*EntityReactionRes
}
type EnumReponse struct {
	Value int    `json:"value"`
	Label string `json:"label"`
}
