package IStrategy

import (
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IPublishPostStrategy interface {
	HandlePublishPost(ctx context.Context, request any, postid primitive.ObjectID) error
	GetType() reflect.Type
}
