package IRepositoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
)

type IPostSettingRepository interface {
	CreatePostSetting(ctx context.Context, Postid string, postSetting *entity.PostSetting) (*entity.PostSetting, error)
	CreateBulkPostSetting(ctx context.Context, postSettings []*entity.PostSetting) (int64, []*dto.BulkError, error)
	GetPostSettingByPostID(ctx context.Context, Postid string) (*entity.PostSetting, error)
	GetPostSettingBulkByPostID(ctx context.Context, Postids []string) ([]*entity.PostSetting, error)
	UpdatePostSetting(ctx context.Context, Postid string, postSetting *entity.PostSetting) (*entity.PostSetting, error)
	UpdateBulkPostSettings(ctx context.Context, postSettings []*entity.PostSetting) (int64, []*dto.BulkError, error)
	DeletePostSetting(ctx context.Context, Postid string) error
	DeleteBulkPostSettings(ctx context.Context, Postids []string) (int64, []*dto.BulkError, error)
	PanigationPostSettings(ctx context.Context, postid string, cursor string, limit int) (*dto.PaginationRes, error)
}
