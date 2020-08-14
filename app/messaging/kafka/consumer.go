package kafka

import (
	msg "cjavellana.me/launchpad/agent/app/messaging/api"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"time"
)

//

type Consumer struct {
	kafkaConsumer *kafka.Consumer

	// the group this consumer belongs
	consumerGroup string

	// the topic this consumer is listening to
	topic string
}

func (c *Consumer) Subscribe(cb msg.Callback) error {
	err := c.kafkaConsumer.Subscribe(c.topic, nil)
	if err != nil {
		return err
	}

	go func() {
		for {
			newMessage, err := c.kafkaConsumer.ReadMessage(-1)
			if err != nil {
				k := err.(kafka.Error)
				if k.Code() == kafka.ErrTimedOut {
					time.Sleep(10 * time.Second)
					continue
				}

				log.Printf("Error reading message from %s, %v", c.topic, err)

				// if error is fatal, stop this go-routine
				// TODO: Propagate error
				break
			}

			cb(newMessage.Value)
		}
	}()

	return nil
}
