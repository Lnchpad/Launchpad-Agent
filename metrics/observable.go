package metrics

type ObserverStatus int

const (
	Stopped ObserverStatus = -1
	Running                = 0
)

type Observable interface {
	Observe()

	StopObserver()
}

func Observe(m []Observable) {
	for _, monitor := range m {
		monitor.Observe()
	}
}

func StopObserver(m []Observable) {
	for _, observer := range m {
		observer.StopObserver()
	}
}