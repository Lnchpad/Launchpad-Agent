package publishers

import (
	"cjavellana.me/launchpad/agent/metrics"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

type MetricsPublisher struct {
	producer *kafka.Producer
}

func NewMetricsPublisher(k *kafka.Producer) *MetricsPublisher  {
	return &MetricsPublisher{producer: k}
}

func (m *MetricsPublisher) PublishCpuMetrics(ch chan metrics.Metrics) {
}

func (m *MetricsPublisher) PublishMemoryMetrics(ch chan metrics.Metrics)  {
}

func (m *MetricsPublisher) PublishNginxStatus(ch chan metrics.NginxStatus) {
}

func (m *MetricsPublisher) PublishServerLogs(ch chan string) {
	topic := "launchpad.agent.logs"

	go func() {
		for {
			msg := <- ch
			err := m.producer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value: []byte(msg),
			}, nil)

			if err != nil {
				log.Print(err)
			}
		}
	}()
}

