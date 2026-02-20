package IRepositoryMongodb

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
)

type IPostRepository interface {
	CreatePost(post *entity.Post) (*entity.Post, error)
	CreateBulkPosts(posts []*entity.Post) ([]*entity.Post, error)
	GetPostByID(postID string) (*entity.Post, error)
	GetPostsBulkByIDs(postIDs []string) ([]*entity.Post, error)
	GetPostsByUserID(userID string) ([]*entity.Post, error)
	UpdatePost(post *entity.Post) (*entity.Post, error)
	UpdateBulkPosts(posts []*entity.Post) (int64, []*dto.BulkError, error)
	DeletePost(postID string) error
	DeleteBulkPosts(postIDs []string) (int64, []*dto.BulkError, error)
	PanigationPosts(userID string, cursor string, limit int) (*dto.PaginationRes, error)
}
