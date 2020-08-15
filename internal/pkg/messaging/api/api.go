package api

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
	//     topic: "launchpad.internal.deployment"
	//   # This refers to consumer id `topictwo`
	//   topictwo:
	//     topic: "launchpad.internal.somethingelse"
	Consumers map[string]map[string]interface{} `yaml:"consumers,omitempty"`
}

type Callback func(message []byte)

// An instance of a `MessageConsumer` represents a consumer of
// a single topic.
type MessageConsumer interface {

	// Subscribes a callback function to a topic (represented by this MessageConsumer)
	// Callback function will be invoked for every message received by the consumer
	//
	// Returns an error when there is an error in subscribing. Not when an error is encountered
	// in the callback
	Subscribe(cb Callback) error
}

type MessageProducer interface {
	Send(message []byte) error
}

type Broker interface {
	NewConsumer(consumerId string) MessageConsumer

	NewProducer(producerId string) MessageProducer
}
