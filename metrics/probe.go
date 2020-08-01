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

func Observe(m []Probe) {
	for _, monitor := range m {
		monitor.Observe()
	}
}

func StopObserver(m []Probe) {
	for _, observer := range m {
		observer.StopObserver()
	}
}
