package kafka

import (
	"log"
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
)

// EnsureTopicExists kiểm tra và tạo topic nếu chưa tồn tại
func EnsureTopicExists(kafkaURL string, topic string, partitions int, replicationFactor int) error {
	conn, err := kafka.Dial("tcp", kafkaURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	var controllerConn *kafka.Conn

	// Kết nối thẳng tới Leader/Controller để tạo topic
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     partitions,
			ReplicationFactor: replicationFactor,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		// Bỏ qua lỗi nếu topic đã tồn tại
		if err.Error() == "topic already exists" {
			log.Printf("Topic '%s' already exists", topic)
			return nil
		}
		return err
	}

	log.Printf("Topic '%s' created with %d partitions", topic, partitions)
	return nil
}
