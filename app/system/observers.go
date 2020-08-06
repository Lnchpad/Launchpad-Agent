package system

import "time"

type TextObserver interface {
	Update(text string)
}

type Metric struct {
	// the time this metric instance was taken
	Timestamp time.Time

	// the type of this metric. e.g. cpu, network, or memory utilization
	Label string
	Value float64
}

type MetricObserver interface {
	Update(metric Metric)
}
