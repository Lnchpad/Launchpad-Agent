package system

type Probe interface {
	Observe(metric Metric)
	StopObserving()
}
