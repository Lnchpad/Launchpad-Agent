package metrics

import (
	"github.com/shirou/gopsutil/mem"
	"log"
	"time"
)

type MemoryProbe struct {
	SamplingInterval time.Duration
	probeStatus ProbeStatus
	subscribers []Subscriber
}

func NewMemoryProbe(samplingInterval time.Duration) *MemoryProbe {
	return &MemoryProbe{SamplingInterval: samplingInterval, probeStatus: Running}
}

func (m *MemoryProbe) SubscribeMany(subscribers []Subscriber) {
	m.subscribers = subscribers
}

func (m *MemoryProbe) Subscribe(subscriber Subscriber) {
	m.subscribers = append(m.subscribers, subscriber)
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
			for _, sub := range m.subscribers {
				if err := sub.Update(Metrics{
					Timestamp: time.Now(),
					Type: TypeMemory,
					Label: "Memory Utilization",
					Value: v.UsedPercent,
				}); err != nil {
					log.Printf("Unable to update %v, reason %s\n", sub, err)
				}
			}

			time.Sleep(m.SamplingInterval)
		}
	}()
}

func (m *MemoryProbe) StopObserver() {
	m.probeStatus = Stopped
}
