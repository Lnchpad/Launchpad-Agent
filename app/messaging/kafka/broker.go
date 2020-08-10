package kafka

import (
	"cjavellana.me/launchpad/agent/app/messaging/api"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"strings"
)

// kafka specific broker.
type Broker struct {
	brokers []string
	config api.BrokerConfig
}

func (k *Broker) NewConsumer(consumerId string) api.MessageConsumer {
	panic("implement me")
}

func (k *Broker) NewProducer(producerId string) api.MessageProducer {
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

func NewKafkaBroker(config api.BrokerConfig) api.Broker {
	servers := config.Hosts
	if servers == nil {
		log.Fatal("Unable to find \"boostrap.servers\" parameter")
	}

	return &Broker{
		brokers: servers,
		config: config,
	}
}

