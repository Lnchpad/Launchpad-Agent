package metrics

type ProbeStatus int

const (
	Stopped ProbeStatus = -1
	Running             = 0
)

type Probe interface {
	Observe()

	StopObserver()
}

