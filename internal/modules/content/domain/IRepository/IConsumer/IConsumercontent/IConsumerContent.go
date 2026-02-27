package IConsumercontent

import "context"

type IContentConsumer interface {
	// Các phương thức liên quan đến việc tiêu thụ nội dung, ví dụ:

	ConsumerReactionPost(ctx context.Context) error
	ConsumerCommentPost(ctx context.Context) error
	ConsumerContentDeletePublishPost(ctx context.Context) error
	ConsumerSharePost(ctx context.Context) error
	// Các phương thức khác tùy theo yêu cầu nghiệp vụ
}
