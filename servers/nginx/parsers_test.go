package nginx

import (
	"reflect"
	"testing"
)

func Test_parseStandardStatus(t *testing.T) {
	type args struct {
		responseBody string
	}
	tests := []struct {
		name string
		args args
		want StandardStatus
	}{
		{"itShouldReturnNginxStatus",
			args{`Active connections: 1
 server accepts handled requests
 10 10 10
 Reading: 0 Writing: 1 Waiting: 0
 `},
			StandardStatus{
				1,
				10,
				10,
				10,
				0,
				1,
				0,
				`Active connections: 1
 server accepts handled requests
 10 10 10
 Reading: 0 Writing: 1 Waiting: 0
 `,
				""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseStandardStatus(tt.args.responseBody); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseStandardStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
