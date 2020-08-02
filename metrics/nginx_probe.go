package metrics

import (
	"cjavellana.me/launchpad/agent/utils/stringutils"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Structure derived from the result of nginx_status
// curl http://127.0.0.1:1337/nginx_status
// Active connections: 1
// server accepts handled requests
// 8 8 8
// Reading: 0 Writing: 1 Waiting: 0
type NginxStatus struct {
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

type NginxProbe struct {
	// The url to be called to obtain server metrics
	StatusUrl        string
	SamplingInterval time.Duration

	// The channel with which metrics will be transmitted
	StatsChannel chan NginxStatus

	// a delay to prevent probing the server before it has fully started
	InitialDelay time.Duration

	// whether the probe is running or stopped
	probeStatus ProbeStatus
}

func NewNginxProbe(statusUrl string, samplingInterval time.Duration, initialDelay time.Duration) *NginxProbe {
	return &NginxProbe{
		StatusUrl:        statusUrl,
		SamplingInterval: samplingInterval,
		StatsChannel:     make(chan NginxStatus),
		InitialDelay:     initialDelay,
	}
}

func (s *NginxProbe) Observe() {
	// calling observe will again start the observer
	s.probeStatus = Running

	go func() {
		// wait for initial delay to lapse before polling to give time
		// to the server to start
		time.Sleep(s.InitialDelay)

		for {
			if s.probeStatus == Stopped {
				log.Print("Stopping NginxStatus Observer...")
				return
			}

			if resp, err := http.Get(s.StatusUrl); err != nil {
				s.StatsChannel <- NginxStatus{ErrorMessage: fmt.Sprintf("%s\n", err)}
			} else {
				s.StatsChannel <- parseStandardStatus(readResponse(resp))
			}

			time.Sleep(s.SamplingInterval)
		}
	}()
}

func (s *NginxProbe) StopObserver() {
	s.probeStatus = Stopped
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

var standardStatusPattern = regexp.MustCompile("^Active connections: (?P<active>\\d+)" +
	"\\s*server accepts handled requests" +
	"\\s*(?P<accepts>\\d+) (?P<handled>\\d+) (?P<requests>\\d+)" +
	"\\s*Reading: (?P<reading>\\d+) Writing: (?P<writing>\\d+) Waiting: (?P<waiting>\\d+)\\s*$")

// Parse the standard nginx response status
// Active connections: 1
// server accepts handled requests
// 8 8 8
// Reading: 0 Writing: 1 Waiting: 0
func parseStandardStatus(responseBody string) NginxStatus {
	responseBodyWithoutNewLines := strings.Replace(responseBody, "\n", "", -1)
	result := standardStatusPattern.FindAllStringSubmatch(responseBodyWithoutNewLines, -1)

	if len(result) > 0 {
		return NginxStatus{
			time.Now(),
			stringutils.ToUint8(result[0][1]),
			stringutils.ToUint8(result[0][2]),
			stringutils.ToUint8(result[0][3]),
			stringutils.ToUint8(result[0][4]),
			stringutils.ToUint8(result[0][5]),
			stringutils.ToUint8(result[0][6]),
			stringutils.ToUint8(result[0][7]),
			responseBody,
			"",
		}
	}

	return NginxStatus{Timestamp: time.Now(), ErrorMessage: fmt.Sprintf("Unable to parse\n%s", responseBody)}
}
