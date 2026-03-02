package producerMedia

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/contentEvent/mediaInContent"
)

type MediaProducer struct {
	eventBus events.EventBus
}

func NewMediaProducer(eventBus events.EventBus) *MediaProducer {
	return &MediaProducer{
		eventBus: eventBus,
	}
}

// Implement các phương thức của IProducerMedia tại đây
func (p *MediaProducer) ProducerPublishPostCreateMediaAssets(ctx context.Context, postID string, items *mediaInContent.CreateMediaAssetsPayload) error {
	// Logic để publish sự kiện tạo media assets cho một post
	
	err := p.eventBus.Publish(ctx, constants.TopicContentPostPublishMediaAssets.String(), postID, constants.Created.String(), items)
	if err != nil {
		return err
	}
	return nil
}
func (p *MediaProducer) ProducerPublishPostDeleteMediaAssets(ctx context.Context, item *mediaInContent.DeleteMediaAssetsPayload) error {
	// Logic để publish sự kiện xóa media assets cho một post
	err := p.eventBus.Publish(ctx, constants.TopicContentPostPublishMediaAssets.String(), item.PostID, constants.Deleted.String(), item)
	if err != nil {
		return err
	}
	return nil
}
