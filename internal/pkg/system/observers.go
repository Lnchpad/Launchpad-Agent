package system

import (
	"cjavellana.me/launchpad/agent/internal/pkg/stats"
)

type TextObserver interface {
	Update(text string)
}

type ProbeObserver interface {
	Update(stats stats.Stats)
}