package messaging

type BrokerType string

const (
	Kafka    BrokerType = "kafka"
	RabbitMQ            = "rabbitmq"
	ZeroMQ              = "zeromq"
	ActiveMQ            = "activemq"
)

type BrokerConfig struct {
	BrokerType BrokerType

	// a hostname/ip-port pair. e.g. localhost:9092
	Hosts []string

	// A hierarchical-structured configuration in the form of:
	// producers:
	//   sometopictopublishto:
	//	   # The topic where all agent metrics are published
	//     topic: "launchpad.metrics"
	//
	Producers map[string]map[string]interface{} `yaml:"producers,omitempty"`

	// A hierarchical-structured configuration in the form of:
	// consumers:
	//   # This refers to consumer id `topicone`
	//   topicone:
	//     topic: "launchpad.app.deployment"
	//   # This refers to consumer id `topictwo`
	//   topictwo:
	//     topic: "launchpad.app.somethingelse"
	Consumers map[string]map[string]interface{} `yaml:"consumers,omitempty"`
}

type MessageConsumer interface {
}

type MessageProducer interface {
	Send(message []byte) error
}

type Broker interface {
	NewConsumer(consumerId string) MessageConsumer

	NewProducer(producerId string) MessageProducer
}

func NewBroker(config BrokerConfig) Broker {
	switch config.BrokerType {
	case "kafka":
		return newKafkaBroker(config)
	}

	return nil
}
