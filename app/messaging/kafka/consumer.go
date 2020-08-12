package kafka

import (
	"cjavellana.me/launchpad/agent/app/messaging/api"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

type Consumer struct {
	kafkaConsumer *kafka.Consumer

	// the group this consumer belongs
	consumerGroup string

	// the topic this consumer is listening to
	topic string
}

func (c *Consumer) Subscribe(cb api.Callback) error {
	err := c.kafkaConsumer.Subscribe(c.topic, nil)
	if err != nil {
		return err
	}

	go func() {
		for {
			msg, err := c.kafkaConsumer.ReadMessage(0)
			if err != nil {
				log.Printf("Error reading message from %s, %v", c.topic, err)
				continue
			}

			cb(msg.Value)
		}
	}()

	return nil
}
