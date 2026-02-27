package client

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/concurrency"

func ProvideLifecycleManager(
// notifConsumer *notifKafka.NotificationConsumer,
// paymentConsumer *paymentKafka.PaymentConsumer,
// Cứ thêm consumer mới vào đây...
) *concurrency.Manager {
	// Trả về Manager chứa danh sách service
	return concurrency.NewManager(
	// notifConsumer,
	// paymentConsumer,	
	)
}
