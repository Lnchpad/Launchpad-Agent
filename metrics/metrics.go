package metrics

import "time"

type Type string

const (
	TypeCpu    Type = "cpu"
	TypeMemory      = "mem"
)

type Metrics struct {
	Timestamp time.Time
	Type      Type
	Label     string
	Value     float64
}

type Subscriber interface {
	Update(metrics Metrics) error
}

type Publisher interface {
	SubscribeMany([]Subscriber)
	Subscribe(Subscriber)
}