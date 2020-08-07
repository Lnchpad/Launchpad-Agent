package system

import (
	"cjavellana.me/launchpad/agent/app/cfg"
	"time"
)

type Status int

const (
	Running Status = iota
	Stopped
)

type Probe interface {
	// ~ Exported functions =====================

	// Returns the type of Probe implementing structs implement
	Type() cfg.ProbeType
	Observe(probeObserver ProbeObserver)
	StopObserving()

	// ~ Package Private methods =====================

	init(samplingInterval time.Duration)
}

func NewProbe(probeType cfg.ProbeType, samplingInterval time.Duration) Probe {
	var p Probe

	switch probeType {
	case cfg.CpuProbe:
		p = &CpuProbe{}
	case cfg.MemProbe:
		p = &MemoryProbe{}
	default:
		return nil
	}

	p.init(samplingInterval)
	return p
}
