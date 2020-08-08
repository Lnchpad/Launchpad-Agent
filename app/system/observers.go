package system

import (
	"cjavellana.me/launchpad/agent/app/stats"
)

type TextObserver interface {
	Update(text string)
}

type ProbeObserver interface {
	Update(stats stats.Stats)
}
