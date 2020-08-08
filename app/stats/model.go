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
	ActiveConnections uint8
	Accepted          uint8
	Handled           uint8
	Requests          uint8
	Reading           uint8
	Writing           uint8
	Waiting           uint8
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
