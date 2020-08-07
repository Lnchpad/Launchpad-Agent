package system

import (
	"cjavellana.me/launchpad/agent/app/cfg"
	"github.com/shirou/gopsutil/cpu"
	"log"
	"time"
)

type CpuProbe struct {
	observers        []ProbeObserver
	samplingInterval time.Duration

	// channel used to terminate the sampling routine
	done             chan bool
	status           Status
}

func (c *CpuProbe) Type() cfg.ProbeType {
	return cfg.CpuProbe
}

func (c *CpuProbe) init(samplingInterval time.Duration) {
	c.samplingInterval = samplingInterval
}

func (c *CpuProbe) Observe(observer ProbeObserver) {
	if c.done == nil {
		c.done = make(chan bool)
	}

	c.observers = append(c.observers, observer)

	if c.status == Stopped {
		c.status = Running
		go c.poll()
	}
}

func (c *CpuProbe) poll() {
	for {
		select {
		case <-c.done:
			log.Print("Terminating CPU Probe")
			return
		default:
			v, _ := cpu.Percent(c.samplingInterval, false)

			for _, o := range c.observers {
				// Fixes: Loop variable `o` captured by func literal
				observer := o

				// Update the observers asynchronously
				go func() {
					observer.Update(Metric{
						Timestamp: time.Now(),
						Label:     string(c.Type()),
						Value:     v[0],
					})
				}()
			}
		}
	}
}

func (c *CpuProbe) StopObserving() {
	c.status = Stopped
	c.done <- true
}
