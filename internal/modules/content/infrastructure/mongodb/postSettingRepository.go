package mongodb

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostSettingRepository struct {
	client *mongo.Database
}

func NewPostSettingRepository(client *mongo.Database) *PostSettingRepository {
	return &PostSettingRepository{client: client}
}
func (r *PostSettingRepository) CreatePostSetting(Postid string, postSetting *entity.PostSetting) (*entity.PostSetting, error) {
	return nil, nil
}
func (r *PostSettingRepository) CreateBulkPostSetting(postSettings []*entity.PostSetting) (int64, []*dto.BulkError, error) {
	return 0, nil, nil
}
func (r *PostSettingRepository) GetPostSettingByPostID(Postid string) (*entity.PostSetting, error) {
	return nil, nil
}
func (r *PostSettingRepository) UpdatePostSetting(Postid string, postSetting *entity.PostSetting) (*entity.PostSetting, error) {
	return nil, nil
}
func (r *PostSettingRepository) UpdateBulkPostSettings(postSettings []*entity.PostSetting) (int64, []*dto.BulkError, error) {
	return 0, nil, nil
}
func (r *PostSettingRepository) DeletePostSetting(Postid string) error {
	return nil
}
func (r *PostSettingRepository) DeleteBulkPostSettings(postids []string) (int64, []*dto.BulkError, error) {
	return 0, nil, nil
}
