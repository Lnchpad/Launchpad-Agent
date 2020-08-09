package stats

import "time"

// Structure derived from the result of nginx_status
// curl http://127.0.0.1:1337/nginx_status
// Active connections: 1
// server accepts handled requests
// 8 8 8
// Reading: 0 Writing: 1 Waiting: 0
type WebServerStats struct {
	Timestamp         time.Time
	ActiveConnections uint32
	Accepted          uint32
	Handled           uint32
	Requests          uint32
	Reading           uint32
	Writing           uint32
	Waiting           uint32
	RawData           string
	ErrorMessage      string
}

type Stats struct {
	// the time this metric instance was taken
	Timestamp time.Time

	// the type of this metric. e.g. cpu, network, or memory utilization
	Label string
	Value float64
}
