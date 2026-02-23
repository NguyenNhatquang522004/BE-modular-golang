package strategy

import (
	"context"
	"errors"
	"reflect"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryCassandra"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostInsightStrategy struct {
	postInsightRepo IRepositoryCassandra.IPostInsights
}

func NewPostInsightStrategy(postInsightRepo IRepositoryCassandra.IPostInsights) *PostInsightStrategy {
	return &PostInsightStrategy{
		postInsightRepo: postInsightRepo,
	}
}

func (p *PostInsightStrategy) HandlePublishPost(ctx context.Context, request any, postid primitive.ObjectID) error {
	req, ok := request.(*req.PostInsightReqv1)
	if !ok {
		return errors.New("invalid request type")
	}
	entity, err := mapper.ToEntityPostInsight(req)
	if err != nil {
		return err
	}
	err = p.postInsightRepo.CreatePostInsightInitPost(ctx, entity)
	if err != nil {
		return err
	}
	return nil
}
func (p *PostInsightStrategy) GetType() reflect.Type {
	return reflect.TypeOf(req.PostInsightReqv1{})
}
