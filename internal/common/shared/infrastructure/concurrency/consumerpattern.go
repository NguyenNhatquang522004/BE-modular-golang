package concurrency

// import (
// 	"context"
// 	"log"

// 	"github.com/IBM/sarama"
//     // Import pkg worker pool của bạn
// 	"your-project/pkg/concurrency" 
// 	"your-project/internal/modules/order/usecase"
// )

// // 1. Định nghĩa Config struct cho Consumer này (Wire sẽ inject cái này)
// type OrderConsumerConfig struct {
// 	WorkerLimit int64
// }

// type OrderConsumer struct {
// 	pool    *concurrency.WorkerPool
// 	useCase usecase.OrderUseCase
//     cfg     OrderConsumerConfig // Lưu config để log hoặc debug nếu cần
// }

// // 2. Constructor chuẩn DI (Wire sẽ gọi hàm này)
// func NewOrderConsumer(cfg OrderConsumerConfig, uc usecase.OrderUseCase) *OrderConsumer {
// 	return &OrderConsumer{
// 		// Khởi tạo pool dựa trên Config được inject
// 		pool:    concurrency.NewWorkerPool(cfg.WorkerLimit),
// 		useCase: uc,
//         cfg:     cfg,
// 	}
// }

// // 3. Implement Sarama ConsumerGroupHandler (Bắt buộc phải đúng interface)
// func (c *OrderConsumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
// func (c *OrderConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// func (c *OrderConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
// 	for msg := range claim.Messages() {
// 		// Dùng hàm xử lý riêng để code gọn
// 		c.processMessage(sess.Context(), msg)
		
// 		// Mark message ngay vì chúng ta xử lý Async (cân nhắc trade-off)
// 		sess.MarkMessage(msg, "") 
// 	}
// 	return nil
// }

// // 4. Logic xử lý chính (Đã tối ưu Context và Log)
// func (c *OrderConsumer) processMessage(ctx context.Context, msg *sarama.ConsumerMessage) {
// 	// Task này chạy Async trong WorkerPool
// 	err := c.pool.Run(ctx, func() {
// 		// A. Tách Context: Dùng WithoutCancel (Go 1.21+) để giữ Tracing nhưng không bị Cancel
//         // Nếu chưa lên Go 1.21, dùng context.Background() là chấp nhận được, 
//         // nhưng tốt nhất là copy values từ ctx cha.
// 		processCtx := context.WithoutCancel(ctx) 
		
// 		// B. Xử lý nghiệp vụ & Log lỗi (KHÔNG ĐƯỢC NUỐT LỖI)
// 		err := c.useCase.ProcessOrder(processCtx, msg.Value)
// 		if err != nil {
// 			log.Printf("[OrderConsumer] Error processing msg offset %d: %v", msg.Offset, err)
//             // TODO: Đẩy vào Dead Letter Queue (DLQ) nếu cần
// 		}
// 	})

// 	if err != nil {
// 		log.Printf("[OrderConsumer] Failed to schedule task (Pool full/closed): %v", err)
// 	}
// }

// // 5. Graceful Shutdown Hook
// func (c *OrderConsumer) Stop() {
// 	log.Println("[OrderConsumer] Waiting for workers to finish...")
// 	c.pool.Wait()
// 	log.Println("[OrderConsumer] All workers finished.")
// }
// func (c *OrderConsumer) Start(ctx context.Context) error {
//     // Gọi Sarama Consumer Group Loop ở đây
//     // Giả sử bạn đã có client Sarama inject vào consumer
//     return c.consumerGroup.Consume(ctx, []string{"order-topic"}, c)
// }

// func (c *OrderConsumer) Name() string {
// 	return "OrderKafkaConsumer"
// }