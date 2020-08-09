package collectors

import (
	"cjavellana.me/launchpad/agent/app/messaging"
	"cjavellana.me/launchpad/agent/app/stats/pb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"log"
	"os"
)

type LogCollector struct {
	messageSender messaging.MessageProducer
}

func NewLogCollector(messageSender messaging.MessageProducer) LogCollector {
	return LogCollector{messageSender: messageSender}
}

func (c *LogCollector) Update(text string) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Unable to obtain hostname %v", err)
		return
	}

	pbLog, err := proto.Marshal(&pb.SimpleLog{
		Timestamp: ptypes.TimestampNow(),
		Hostname:  hostname,
		Service:   "web",
		Type:      "nginx",
		Message:   text,
	})
	if err != nil {
		log.Printf("Unable to marshal log %v", err)
		return
	}

	err = c.messageSender.Send(pbLog)
	if err != nil {
		log.Printf("Unable to send stats %v", err)
	}
}
