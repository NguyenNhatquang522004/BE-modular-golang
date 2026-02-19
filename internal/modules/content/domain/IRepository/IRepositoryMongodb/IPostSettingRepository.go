package IRepositoryMongodb

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
)

type IPostSettingRepository interface {
	CreatePostSetting(Postid string, postSetting *entity.PostSetting) (*entity.PostSetting, error)
	CreateBulkPostSetting(postSettings []*entity.PostSetting) (int64, []*dto.BulkError, error)
	GetPostSettingByPostID(Postid string) (*entity.PostSetting, error)
	UpdatePostSetting(Postid string, postSetting *entity.PostSetting) (*entity.PostSetting, error)
	UpdateBulkPostSettings(postSettings []*entity.PostSetting) (int64, []*dto.BulkError, error)
	DeletePostSetting(Postid string) error
	DeleteBulkPostSettings(Postids []string) (int64, []*dto.BulkError, error)
}
