package collectors

import (
	"cjavellana.me/launchpad/agent/internal/pkg/messaging/api"
	"cjavellana.me/launchpad/agent/internal/pkg/stats/pb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"log"
	"os"
)

type WebStatsCollector struct {
	messageSender api.MessageProducer
}

func NewWebStatsCollector(messageSender api.MessageProducer) WebStatsCollector {
	return WebStatsCollector{messageSender: messageSender}
}

func (w *WebStatsCollector) Update(text string) {
	hostname, _ := os.Hostname()
	stats, err := parseNginxStatus(text)

	var webStats pb.WebServerStats
	if err != nil {
		webStats = pb.WebServerStats{
			Timestamp: ptypes.TimestampNow(),
			Service:   "web",
			Type:      "nginxstats",
			Hostname:  hostname,
			RawData:   text,
			Error:     true,
		}
	} else {
		webStats = pb.WebServerStats{
			Timestamp:         ptypes.TimestampNow(),
			Service:           "web",
			Type:              "nginxstats",
			Hostname:          hostname,
			ActiveConnections: stats.ActiveConnections,
			Accepted:          stats.Accepted,
			Handled:           stats.Handled,
			Reading:           stats.Reading,
			Writing:           stats.Writing,
			Requests:          stats.Requests,
			Waiting:           stats.Waiting,
			RawData:           text,
			Error:             false,
		}
	}

	pbLog, err := proto.Marshal(&webStats)
	if err != nil {
		log.Printf("Unable to marshal log %v", err)
		return
	}

	err = w.messageSender.Send(pbLog)
	if err != nil {
		log.Printf("Unable to send stats %v", err)
	}
}
