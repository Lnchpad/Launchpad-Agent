package messaging

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

type KafkaConfig struct {
	Broker string
}

type KafkaBroker struct {
	broker string
}

func NewKafkaBroker(c KafkaConfig) *KafkaBroker {
	return &KafkaBroker{broker: c.Broker}
}

func (k *KafkaBroker) NewProducer() *kafka.Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": k.broker,
	})

	if err != nil {
		log.Print(err)
		return nil
	}

	return p
}

func (k *KafkaBroker) NewConsumer() *kafka.Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": k.broker,
	})

	if err != nil {
		log.Print(err)
		return nil
	}

	return c
}