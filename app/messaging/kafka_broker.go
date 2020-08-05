package messaging

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"strings"
)

type KafkaBroker struct {
	brokers []string

	config BrokerConfig
}

func (k *KafkaBroker) NewConsumer(consumerId string) MessageConsumer {
	panic("implement me")
}

func (k *KafkaBroker) NewProducer(producerId string) MessageProducer {
	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": strings.Join(k.brokers[:], ","),
	})
	if err != nil {
		log.Fatal(err)
	}

	return &Producer{
		producer: kafkaProducer,
	}
}

func newKafkaBroker(config BrokerConfig) Broker {
	servers := config.Hosts
	if servers == nil {
		log.Fatal("Unable to find \"boostrap.servers\" parameter")
	}

	return &KafkaBroker{
		brokers: servers,
		config: config,
	}
}

