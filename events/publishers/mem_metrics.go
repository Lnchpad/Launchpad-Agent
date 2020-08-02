package publishers

import (
	"cjavellana.me/launchpad/agent/metrics"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var memMetricsTopic = "launchpad.agent.memory"

type MemMetricsPublisher struct {
	producer *kafka.Producer
}

func NewMemMetricsPublisher(producer *kafka.Producer) *MemMetricsPublisher {
	return &MemMetricsPublisher{
		producer: producer,
	}
}

func (memMetricsPub *MemMetricsPublisher) Update(metrics metrics.Metrics) error {
	return publishMetrics(metrics, memMetricsPub.producer, memMetricsTopic)
}