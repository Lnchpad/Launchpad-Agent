package collectors

import (
	"cjavellana.me/launchpad/agent/internal/pkg/messaging/api"
	"cjavellana.me/launchpad/agent/internal/pkg/stats"
	"cjavellana.me/launchpad/agent/internal/pkg/stats/pb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"log"
	"os"
)

// collects cpu, memory, network utilization type of stats
type StatsCollector struct {
	messageSender api.MessageProducer
}

func NewStatsCollector(messageSender api.MessageProducer) StatsCollector {
	return StatsCollector{messageSender: messageSender}
}

func (c *StatsCollector) Update(s stats.Stats) {
	hostname, _ := os.Hostname()

	bStats, _ := proto.Marshal(&pb.Metrics{
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
