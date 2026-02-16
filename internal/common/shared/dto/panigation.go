package dto

type PaginationReq struct {
	Cursor string `form:"cursor"`
	Limit  int    `form:"limit" validate:"gte=1,lte=100"`
}
type PaginationRes struct {
	NextCursor string `json:"next_cursor,omitempty"`
	HasNext    bool   `json:"has_next"`
	Data       any    `json:"data"`
	Limit      int    `json:"-"`
}
