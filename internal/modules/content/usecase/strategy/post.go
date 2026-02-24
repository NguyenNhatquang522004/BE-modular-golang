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

type PostStrategy struct {
	postRepo IRepositoryMongodb.IPostRepository
}

func NewPostStrategy(postRepo IRepositoryMongodb.IPostRepository) *PostStrategy {
	return &PostStrategy{
		postRepo: postRepo,
	}
}
func (p *PostStrategy) HandlePublishPost(ctx context.Context, request any, postid primitive.ObjectID) error {
	// Implement the logic to handle post extension
	req, ok := request.(*req.PostReq)
	if !ok {
		return errors.New("invalid request type")
	}
	entity := mapper.ToEntityPost(req)
	entity.ID = postid // Gán postID từ tham số vào entity
	_, err := p.postRepo.CreatePost(ctx, entity)
	if err != nil {
		return err
	}
	return nil
}
func (p *PostStrategy) GetType() reflect.Type {
	return reflect.TypeOf(req.PostReq{})
}

func (p *PostStrategy) HandlePublishDelete(ctx context.Context, request *req.DeletePostRequest) error {
	// Implement the logic to handle post deletion
	err := p.postRepo.DeletePost(ctx, request.PostID.PostID)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostStrategy) GetDeleteType() reflect.Type {
	return reflect.TypeOf(p.postRepo)
}
