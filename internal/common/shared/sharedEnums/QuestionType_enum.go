package sharedEnums

//go:generate enumer -type=QuestionType -json -transform=snake -trimprefix=QuestionType
type QuestionType int

const (
	QuestionTypeText           QuestionType = iota // 'text' (Tự luận)
	QuestionTypeMultipleChoice                     // 'multiple_choice' (Chọn 1 trong nhiều)
	QuestionTypeCheckbox                           // 'checkbox' (Chọn nhiều)
)
