package metrics

import (
	"github.com/shirou/gopsutil/cpu"
	"log"
	"time"
)

type CpuProbe struct {
	SamplingInterval time.Duration

	// The channel with which metrics will be transmitted
	Channel     chan Metrics
	probeStatus ProbeStatus
}

func NewCpuProbe(samplingInterval time.Duration) *CpuProbe {
	ch := make(chan Metrics)
	return &CpuProbe{SamplingInterval: samplingInterval, Channel: ch}
}

func (c *CpuProbe) Observe() {
	// calling observe will again start the observer
	c.probeStatus = Running

	go func() {
		for {
			if c.probeStatus == Stopped {
				log.Print("Stopping Memory Observer...")
				return
			}

			v, _ := cpu.Percent(c.SamplingInterval, false)
			c.Channel <- Metrics{Label: "CpuProbe Utilization", Value: v[0]}
		}
	}()
}

func (c *CpuProbe) StopObserver()  {
	c.probeStatus = Stopped
}
