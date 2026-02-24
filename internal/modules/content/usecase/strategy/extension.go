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

type PostExtensionStrategy struct {
	postExtensionRepo IRepositoryMongodb.IPostExtensionRepository
}

func NewPostExtensionStrategy(postExtensionRepo IRepositoryMongodb.IPostExtensionRepository) *PostExtensionStrategy {
	return &PostExtensionStrategy{
		postExtensionRepo: postExtensionRepo,
	}
}

func (p *PostExtensionStrategy) HandlePublishPost(ctx context.Context, request any, postid primitive.ObjectID) error {
	// Implement the logic to handle post extension
	req, ok := request.(*req.PostExtensionReq)
	if !ok {
		return errors.New("invalid request type")
	}
	entity := mapper.ToPostExtensionEntityPostExtension(req)
	entity.PostID = postid
	err := p.postExtensionRepo.CreatePostExtension(ctx, entity)
	if err != nil {
		return err
	}
	return nil
}
func (p *PostExtensionStrategy) GetType() reflect.Type {
	return reflect.TypeOf(req.PostExtensionReq{})
}
func (p *PostExtensionStrategy) HandlePublishDelete(ctx context.Context, request *req.DeletePostRequest) error {
	// Implement the logic to handle post deletion
	err := p.postExtensionRepo.DeleteByPostID(ctx, request.PostID.PostID)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostExtensionStrategy) GetDeleteType() reflect.Type {
	return reflect.TypeOf(p.postExtensionRepo)
}
