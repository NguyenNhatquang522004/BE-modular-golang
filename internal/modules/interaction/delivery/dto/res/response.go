package res

type CommentDetailResponse struct {
	Comment   *CommentRes
	Edit      []*CommentEditLogRes
	Reactions []*EntityReactionRes
}
