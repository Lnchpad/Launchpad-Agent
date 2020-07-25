package listeners

// listens for portal events e.g. add, remove, or update of portal apps

type MessageBrokerVendor string

const (
	Kafka MessageBrokerVendor = "kafka"

	// future
	RabbitMQ = "rabbitmq"
	ZeroMQ   = "zeromq"
	ActiveMQ = "activemq"
)

type MessageBroker interface {
	connect()
}

type EventListenerConfig struct {
	MessageBrokerVendor string

	KafkaConfig KafkaConfig
}

func InitializeEventListener(c EventListenerConfig) MessageBroker {
	switch c.MessageBrokerVendor {
	case "kafka":
		m := &KafkaBroker{}
		m.configure(c.KafkaConfig)
		m.connect()
		return m
	default:
		return nil
	}
}
