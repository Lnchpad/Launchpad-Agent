package cfg

import (
	"cjavellana.me/launchpad/agent/app/messaging/api"
	"reflect"
	"testing"
)

func Test_parseConfigYaml(t *testing.T) {
	type args struct {
		yamlConfig []byte
	}
	tests := []struct {
		name string
		args args
		want AppConfig
	}{
		{
			name: "itShouldReturnAnAppConfig",
			args: args{
				yamlConfig: []byte(`
apiversion: 1.0
viewtype: none
probeconfig:
  enabled: true
  probetypes:
    - cpu
    - memory
brokerconfig:
  hosts: 
    - localhost:9092
  consumers:
    sometopic:
      topic: helloworld
`),
			},
			want: AppConfig{
				ViewType:     "none",
				MaxSeriesElements: 15,
				ServerConfig: ServerConfig{},
				ProbeConfig: ProbeConfig{
					Enabled:    true,
					ProbeTypes: []ProbeType{CpuProbe, MemProbe},
				},
				BrokerConfig: api.BrokerConfig{
					Hosts: []string{"localhost:9092"},
					Consumers: map[string]map[string]interface{}{
						// The yaml library unmarshalls the
						// inner map as map[interface{}] interface
						"sometopic": {
							"topic": "helloworld",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseConfigYaml(tt.args.yamlConfig); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseConfigYaml() = %v, want %v", got, tt.want)
			}
		})
	}
}
