package producercontent

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/contentEvent"
)

type ContentProducer struct {
	eventBus events.EventBus
}

func NewContentProducer(eventBus events.EventBus) *ContentProducer {
	return &ContentProducer{
		eventBus: eventBus,
	}
}

// Implement các phương thức của IProducerContent tại đây
func (p *ContentProducer) PublishContentPostStats(ctx context.Context, payload ...*contentEvent.PostStatsPayload) error {
	for _, items := range payload {
		err := p.eventBus.Publish(ctx, constants.TopicContentPostPublish.String(), items.PostID, constants.Updated.String(), items)
		if err != nil {
			return err
		}
	}
	return nil
}
func (p *ContentProducer) PublishContentDeletePublishPost(ctx context.Context, payload ...*contentEvent.PostDeletePayload) error {
	for _, items := range payload {
		err := p.eventBus.Publish(ctx, constants.TopicContentPostPublish.String(), items.PostID, constants.Deleted.String(), items)
		if err != nil {
			return err
		}
	}
	return nil
}
