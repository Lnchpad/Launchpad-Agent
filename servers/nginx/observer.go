package nginx

import (
	"cjavellana.me/launchpad/agent/metrics"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Structure derived from the result of nginx_status
// curl http://127.0.0.1:1337/nginx_status
// Active connections: 1
// server accepts handled requests
// 8 8 8
// Reading: 0 Writing: 1 Waiting: 0
type StandardStatus struct {
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

type StatusObserver struct {
	// The url to be called to obtain server metrics
	StatusUrl        string
	SamplingInterval time.Duration

	// The channel with which metrics will be transmitted
	Channel chan StandardStatus

	InitialDelay time.Duration

	observerStatus metrics.ObserverStatus
}

func NewStatusObserver(statusUrl string, samplingInterval time.Duration, initialDelay time.Duration) *StatusObserver {
	return &StatusObserver{
		StatusUrl:        statusUrl,
		SamplingInterval: samplingInterval,
		Channel:          make(chan StandardStatus),
		InitialDelay:     initialDelay,
	}
}

func (s *StatusObserver) Observe() {
	// calling observe will again start the observer
	s.observerStatus = metrics.Running

	go func() {
		// wait for initial delay to lapse before polling to give time
		// to the server to start
		time.Sleep(s.InitialDelay)

		for {
			if s.observerStatus == metrics.Stopped {
				log.Print("Stopping NginxStatus Observer...")
				return
			}

			if resp, err := http.Get(s.StatusUrl); err != nil {
				s.Channel <- StandardStatus{ErrorMessage: fmt.Sprintf("%s\n", err)}
			} else {
				s.Channel <- parseStandardStatus(readResponse(resp))
			}

			time.Sleep(s.SamplingInterval)
		}
	}()
}

func (s *StatusObserver) StopObserver() {
	s.observerStatus = metrics.Stopped
}

func readResponse(resp *http.Response) string {
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		log.Println(err)
		return ""
	}

	return string(body)
}
