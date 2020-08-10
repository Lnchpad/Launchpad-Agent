package sync

import "cjavellana.me/launchpad/agent/app/messaging"

type PortalEventListener struct {
	msgConsumer *messaging.MessageConsumer
}

func NewPortalEventListener(msgConsumer *messaging.MessageConsumer) PortalEventListener {
	return PortalEventListener{
		msgConsumer: msgConsumer,
	}
}

func (p *PortalEventListener) onMessage() {
}
