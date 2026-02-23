package strategy

import (
	"context"
	"errors"
	"reflect"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryMongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SettingStrategy struct {
	// Thêm các repository hoặc service cần thiết để xử lý logic liên quan đến setting
	settingRepo IRepositoryMongodb.IPostSettingRepository
}

func NewSettingStrategy(settingRepo IRepositoryMongodb.IPostSettingRepository) *SettingStrategy {
	return &SettingStrategy{
		settingRepo: settingRepo,
	}
}

func (s *SettingStrategy) HandlePublishPost(ctx context.Context, request any, postid primitive.ObjectID) error {
	// Implement the logic to handle post settings
	// Chuyển đổi request sang định dạng entity phù hợp và lưu vào repository nếu cần thiết
	data, ok := request.(*req.PostSettingReq) // Giả sử có req.PostSettingReq
	if !ok {
		return errors.New("invalid request format") // Hoặc trả về lỗi nếu request không đúng định dạng
	}
	dataEntity := mapper.ToPostSettingEntityPostSetting(data)
	dataEntity.PostID = postid // Gán postID từ tham số vào entity
	_, err := s.settingRepo.CreatePostSetting(ctx, dataEntity)
	if err != nil {
		return err
	}
	return nil
}
func (s *SettingStrategy) GetType() reflect.Type {
	return reflect.TypeOf(req.PostSettingReq{})
}
