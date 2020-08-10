package sync

import (
	"cjavellana.me/launchpad/agent/app/messaging/api"
)

type PortalEventListener struct {
	msgConsumer *api.MessageConsumer
}

func NewPortalEventListener(msgConsumer *api.MessageConsumer) PortalEventListener {
	return PortalEventListener{
		msgConsumer: msgConsumer,
	}
}

func (p *PortalEventListener) onMessage() {
}
