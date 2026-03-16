package IConsumer

import "context"

type IConsumerAIAnalysisMediaAsset interface {
	// ConsumeAIAnalysisMediaAsset: Tiêu thụ dữ liệu phân tích AI từ Kafka và lưu vào Redis
	ConsumerAIAnalysisMediaAsset(ctx context.Context) error
	ConsumerFailedAIAnalysisMediaAsset(ctx context.Context) error
}
