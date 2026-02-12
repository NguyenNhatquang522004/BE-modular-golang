package dto

type PaginationRes struct {
	NextCursor string `json:"next_cursor,omitempty"`
	HasNext    bool   `json:"has_next"`
	Data       any    `json:"data"`
	Limit      int    `json:"-"`
}
