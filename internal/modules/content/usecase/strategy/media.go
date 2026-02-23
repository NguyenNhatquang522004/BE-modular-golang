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

type MediaStrategy struct {
	mediaRepo IRepositoryMongodb.IPostMediaRepository // Thay 'any' bằng interface repository tương ứng khi có
}

func NewMediaStrategy(mediaRepo IRepositoryMongodb.IPostMediaRepository) *MediaStrategy {
	return &MediaStrategy{
		mediaRepo: mediaRepo,
	}
}

func (m *MediaStrategy) HandlePublishPost(ctx context.Context, request any, postid primitive.ObjectID) error {
	// Implement the logic to handle media
	// Chuyển đổi request sang định dạng entity phù hợp và lưu vào repository
	req, ok := request.(*req.PostMediaReq) // Giả sử có req.PostMediaReq
	if !ok {
		return errors.New("invalid request type for media strategy")
	}
	entity := mapper.ToPostMediaEntityPostMedia(req)
	entity.PostID = postid // Gán postID từ tham số vào entity
	err := m.mediaRepo.CreatePostMedia(ctx, entity)
	if err != nil {
		return err
	}
	return nil
}
func (m *MediaStrategy) GetType() reflect.Type {
	return reflect.TypeOf(req.PostMediaReq{})
}
