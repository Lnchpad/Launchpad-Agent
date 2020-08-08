package messaging

import (
	"fmt"
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
	if kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": strings.Join(k.brokers[:], ","),
	}); err != nil {
		log.Fatal(err)
	} else {
		if topic := k.config.Producers[producerId]; topic != nil {
			return &Producer{
				producer: kafkaProducer,
				topic: fmt.Sprintf("%v", topic["topic"]),
			}
		}
	}

	return nil
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

