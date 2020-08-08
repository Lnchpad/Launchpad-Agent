package stats

import (
	"cjavellana.me/launchpad/agent/app/messaging"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"log"
	"os"
)

// collects cpu, memory, network utilization type of stats
type BasicStatsCollector struct {
	messageSender messaging.MessageProducer
}

func NewBasicStatsCollector(messageSender messaging.MessageProducer) BasicStatsCollector {
	return BasicStatsCollector{messageSender: messageSender}
}

func (c *BasicStatsCollector) Update(stats Stats) {
	hostname, _ := os.Hostname()

	bStats, _ := proto.Marshal(&Metrics{
		Timestamp: ptypes.TimestampNow(),
		Type:      "web",
		Hostname:  hostname,
		Label:     stats.Label,
		Value:     float32(stats.Value),
	})

	err := c.messageSender.Send(bStats)
	if err != nil {
		log.Printf("Unable to send stats %v", err)
	}
}
