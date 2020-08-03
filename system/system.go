package system

import "cjavellana.me/launchpad/agent/metrics"

func StopObserver(m []metrics.Probe) {
	for _, observer := range m {
		observer.StopObserver()
	}
}

func Observe(m []metrics.Probe) {
	for _, monitor := range m {
		monitor.Observe()
	}
}

