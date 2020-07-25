package metrics

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"log"
	"time"
)

type Metrics struct {
	Label string
	Value float64
}

// ~ Cpu Observer ==============================

type CpuObserver struct {
	SamplingInterval time.Duration

	// The channel with which metrics will be transmitted
	Channel          chan Metrics
	observerStatus ObserverStatus
}

func NewCpuObserver(samplingInterval time.Duration) *CpuObserver {
	ch := make(chan Metrics)
	return &CpuObserver{SamplingInterval: samplingInterval, Channel: ch}
}

func (c *CpuObserver) Observe() {
	// calling observe will again start the observer
	c.observerStatus = Running

	go func() {
		for {
			if c.observerStatus == Stopped {
				log.Print("Stopping Memory Observer...")
				return
			}

			v, _ := cpu.Percent(c.SamplingInterval, false)
			c.Channel <- Metrics{Label: "CpuObserver Utilization", Value: v[0]}
		}
	}()
}

func (c *CpuObserver) StopObserver()  {
	c.observerStatus = Stopped
}

// ~ Memory Observer ==============================

type MemoryObserver struct {
	SamplingInterval time.Duration
	MetricsChannel   chan Metrics

	observerStatus ObserverStatus
}

func NewMemoryObserver(samplingInterval time.Duration) *MemoryObserver {
	ch := make(chan Metrics)
	return &MemoryObserver{SamplingInterval: samplingInterval, MetricsChannel: ch, observerStatus: Running}
}

func (m *MemoryObserver) Observe() {
	// calling observe will again start the observer
	m.observerStatus = Running

	go func() {
		for {
			if m.observerStatus == Stopped {
				log.Print("Stopping Memory Observer...")
				return
			}

			v, _ := mem.VirtualMemory()
			m.MetricsChannel <- Metrics{Label: "Memory Utilization", Value: v.UsedPercent}
			time.Sleep(m.SamplingInterval)
		}
	}()
}

func (m *MemoryObserver) StopObserver()  {
	m.observerStatus = Stopped
}