package kafka

import (
	"cjavellana.me/launchpad/agent/internal/pkg/messaging/api"
	"reflect"
	"testing"
)

func TestNewKafkaBroker(t *testing.T) {
	type args struct {
		config api.BrokerConfig
	}

	producerConfig := map[string]map[string]interface{}{
		"metricsproducer": {
			"topic": "launchpad.metrics",
		},
	}

	brokerCfg := api.BrokerConfig{
		BrokerType: api.Kafka,
		Hosts:      []string{"localhost"},
		Producers:  producerConfig,
	}

	tests := []struct {
		name string
		args args
		want api.Broker
	}{
		{
			name: "itShouldReturnNewKafkaBroker",
			args: args{
				config: brokerCfg,
			},
			want: &Broker{
				brokers: []string{"localhost"},
				config: brokerCfg,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewKafkaBroker(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newKafkaBroker() = %v, want %v", got, tt.want)
			}
		})
	}
}
