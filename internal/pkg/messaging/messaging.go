package messaging

import (
	"cjavellana.me/launchpad/agent/internal/pkg/messaging/api"
	"cjavellana.me/launchpad/agent/internal/pkg/messaging/kafka"
)

func NewBroker(config api.BrokerConfig) api.Broker {
	switch config.BrokerType {
	case "kafka":
		return kafka.NewKafkaBroker(config)
	}

	return nil
}
