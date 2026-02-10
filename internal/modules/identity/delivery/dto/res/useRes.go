package res

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"

type PanigationUsersRes struct {
	Users      []*entity.User `json:"users"`
	NextCursor string         `json:"next_cursor,omitempty"`
}
