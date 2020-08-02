package publishers

import (
	"cjavellana.me/launchpad/agent/messaging/protobuf"
	"cjavellana.me/launchpad/agent/metrics"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

var logsTopic = "launchpad.agent.logs"

type MetricsPublisher struct {
	producer *kafka.Producer
}

func NewMetricsPublisher(k *kafka.Producer) *MetricsPublisher {
	return &MetricsPublisher{producer: k}
}

func (m *MetricsPublisher) PublishServerLogs(logs string) error {
	simpleLog := &protobuf.SimpleLog{Timestamp: ptypes.TimestampNow(), Message: logs}
	if simpleLogAsBytes, err := proto.Marshal(simpleLog); err != nil {
		return err
	} else {
		err := m.producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &logsTopic, Partition: kafka.PartitionAny},
			Value:          simpleLogAsBytes,
		}, nil)

		if err != nil {
			return err
		}
	}

	return nil
}

func publishMetrics(metrics metrics.Metrics, producer *kafka.Producer, topic string) error {
	cpuMetrics := protobuf.Metrics{
		Timestamp: ptypes.TimestampNow(),
		Type: string(metrics.Type),
		Label: metrics.Label,
		Value: float32(metrics.Value),
	}

	metricsAsBytes, err := proto.Marshal(&cpuMetrics)
	if err != nil {
		return err
	}

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          metricsAsBytes,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}