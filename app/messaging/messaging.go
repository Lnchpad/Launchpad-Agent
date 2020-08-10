package messaging

import (
	"cjavellana.me/launchpad/agent/app/messaging/api"
	"cjavellana.me/launchpad/agent/app/messaging/kafka"
)

func NewBroker(config api.BrokerConfig) api.Broker {
	switch config.BrokerType {
	case "kafka":
		return kafka.NewKafkaBroker(config)
	}

	return nil
}
