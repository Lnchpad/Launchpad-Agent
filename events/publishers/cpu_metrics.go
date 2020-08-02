package publishers

import (
	"cjavellana.me/launchpad/agent/metrics"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var cpuMetricsTopic = "launchpad.agent.cpu"

type CpuMetricsPublisher struct {
	producer *kafka.Producer
}

func NewCpuMetricsPublisher(producer *kafka.Producer) *CpuMetricsPublisher {
	return &CpuMetricsPublisher{
		producer: producer,
	}
}

func (cpuMetricsPub *CpuMetricsPublisher) Update(cpu metrics.Metrics) error {
	return publishMetrics(cpu, cpuMetricsPub.producer, cpuMetricsTopic)
}