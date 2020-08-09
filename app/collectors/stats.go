package collectors

import (
	"cjavellana.me/launchpad/agent/app/messaging"
	"cjavellana.me/launchpad/agent/app/stats"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"log"
	"os"
)

// collects cpu, memory, network utilization type of stats
type StatsCollector struct {
	messageSender messaging.MessageProducer
}

func NewStatsCollector(messageSender messaging.MessageProducer) StatsCollector {
	return StatsCollector{messageSender: messageSender}
}

func (c *StatsCollector) Update(s stats.Stats) {
	hostname, _ := os.Hostname()

	bStats, _ := proto.Marshal(&stats.Metrics{
		Timestamp: ptypes.TimestampNow(),
		Service:   "web",
		Type:      s.Label,
		Hostname:  hostname,
		Label:     s.Label,
		Value:     float32(s.Value),
	})

	err := c.messageSender.Send(bStats)
	if err != nil {
		log.Printf("Unable to send stats %v", err)
	}
}
