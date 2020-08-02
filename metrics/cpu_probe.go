package metrics

import (
	"github.com/shirou/gopsutil/cpu"
	"log"
	"time"
)

type CpuProbe struct {
	SamplingInterval time.Duration

	probeStatus    ProbeStatus
	subscribers []Subscriber
}

func NewCpuProbe(samplingInterval time.Duration) *CpuProbe {
	return &CpuProbe{SamplingInterval: samplingInterval}
}

func (c *CpuProbe) SubscribeMany(subscribers []Subscriber) {
	c.subscribers = subscribers
}

func (c *CpuProbe) Subscribe(subscriber Subscriber) {
	c.subscribers = append(c.subscribers, subscriber)
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

			for _, sub := range c.subscribers {
				if err := sub.Update(Metrics{
					Timestamp: time.Now(),
					Type:  TypeCpu,
					Label: "CpuProbe Utilization",
					Value: v[0],
				}); err != nil {
					log.Printf("Unable to update %v, reason %s\n", sub, err)
				}
			}
		}
	}()
}

func (c *CpuProbe) StopObserver()  {
	c.probeStatus = Stopped
}
