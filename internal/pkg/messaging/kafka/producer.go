package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// An instance of a messaging.Producer for sending messages to Kafka.
//
// A single instance of this producer can only send to one topic.
//

type Producer struct {
	producer *kafka.Producer
	topic    string
}

func (p *Producer) Send(message []byte) error {
	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}
