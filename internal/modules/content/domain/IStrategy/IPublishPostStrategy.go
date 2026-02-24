package IStrategy

import (
	"context"
	"reflect"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IPublishPostStrategy interface {
	HandlePublishPost(ctx context.Context, request any, postid primitive.ObjectID) error
	GetType() reflect.Type
}
type IPublishDeleteStrategy interface {
	HandlePublishDelete(ctx context.Context, request *req.DeletePostRequest) error
	GetDeleteType() reflect.Type
}
