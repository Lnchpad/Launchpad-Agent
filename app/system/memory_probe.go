package system

import (
	"cjavellana.me/launchpad/agent/app/cfg"
	"github.com/shirou/gopsutil/mem"
	"log"
	"time"
)

type MemoryProbe struct {
	observers        []ProbeObserver
	samplingInterval time.Duration

	// channel used to terminate the sampling routine
	done   chan bool
	status Status
}

func (m *MemoryProbe) Type() cfg.ProbeType {
	return cfg.MemProbe
}

func (m *MemoryProbe) init(samplingInterval time.Duration) {
	m.samplingInterval = samplingInterval
}

func (m *MemoryProbe) Observe(observer ProbeObserver) {
	if m.done == nil {
		m.done = make(chan bool)
	}

	m.observers = append(m.observers, observer)

	if m.status == Stopped {
		m.status = Running
		go m.poll()
	}
}

func (m *MemoryProbe) poll() {
	for {
		select {
		case <-m.done:
			log.Print("Terminating Memory Probe")
			return
		default:
			v, _ := mem.VirtualMemory()
			for _, o := range m.observers {
				// Fixes: Loop variable `o` captured by func literal
				observer := o

				// Update the observers asynchronously
				go func() {
					observer.Update(Metric{
						Timestamp: time.Now(),
						Label:     cfg.MemProbe,
						Value:     v.UsedPercent,
					})
				}()
			}

			time.Sleep(m.samplingInterval)
		}
	}
}

func (m *MemoryProbe) StopObserving() {
	m.status = Stopped
	m.done <- true
}
