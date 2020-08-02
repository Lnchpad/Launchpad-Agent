package metrics

import (
	"github.com/shirou/gopsutil/mem"
	"log"
	"time"
)

type MemoryProbe struct {
	SamplingInterval time.Duration
	MetricsChannel   chan Metrics

	probeStatus ProbeStatus
}

func NewMemoryProbe(samplingInterval time.Duration) *MemoryProbe {
	ch := make(chan Metrics)
	return &MemoryProbe{SamplingInterval: samplingInterval, MetricsChannel: ch, probeStatus: Running}
}

func (m *MemoryProbe) Observe() {
	// calling observe will again start the observer
	m.probeStatus = Running

	go func() {
		for {
			if m.probeStatus == Stopped {
				log.Print("Stopping Memory Observer...")
				return
			}

			v, _ := mem.VirtualMemory()
			m.MetricsChannel <- Metrics{
				Timestamp: time.Now(),
				Type: TypeMemory,
				Label: "Memory Utilization",
				Value: v.UsedPercent,
			}
			time.Sleep(m.SamplingInterval)
		}
	}()
}

func (m *MemoryProbe) StopObserver() {
	m.probeStatus = Stopped
}
