package messaging

import (
	"reflect"
	"testing"
)

func TestNewKafkaBroker(t *testing.T) {
	type args struct {
		config BrokerConfig
	}

	producerConfig := map[string]interface{}{
		"metricsproducer": map[string]interface{}{
			"topic": "launchpad.metrics",
		},
	}

	brokerCfg := BrokerConfig{
		BrokerType: Kafka,
		Hosts:      []string{"localhost"},
		Producers:  producerConfig,
	}

	tests := []struct {
		name string
		args args
		want Broker
	}{
		{
			name: "itShouldReturnNewKafkaBroker",
			args: args{
				config: brokerCfg,
			},
			want: &KafkaBroker{
				brokers: []string{"localhost"},
				config: brokerCfg,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newKafkaBroker(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newKafkaBroker() = %v, want %v", got, tt.want)
			}
		})
	}
}
