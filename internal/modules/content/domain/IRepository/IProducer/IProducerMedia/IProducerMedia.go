package IProducerMedia

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/contentEvent/mediaInContent"
)

type IProducerMedia interface {
	// CreatePostMedia tạo mới document PostMedia cho một Post (chứa nhiều media item)
	ProducerPublishPostCreateMediaAssets(ctx context.Context, postID string, items []*mediaInContent.CreateMediaAssetsPayload) error
	ProducerPublishPostDeleteMediaAssets(ctx context.Context, item *mediaInContent.DeleteMediaAssetsPayload) error
}
