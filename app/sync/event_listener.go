package sync

import (
	"cjavellana.me/launchpad/agent/app/messaging/api"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
)

type PortalEventListener struct {
	consumer  api.MessageConsumer
	fsUpdater FsUpdater
}

func NewPortalEventListener(consumer api.MessageConsumer, updater FsUpdater) *PortalEventListener {
	return &PortalEventListener{
		consumer: consumer,
		fsUpdater: updater,
	}
}

func (p *PortalEventListener) StartListening() {
	_ = p.consumer.Subscribe(onMessage)
}

func onMessage(cb []byte) {
	app := &PortalAppDeployment{}
	if err := proto.Unmarshal(cb, app); err != nil {
		log.Printf("Unable to unmarshall %v", err)
	}

	fmt.Println(app)
}
