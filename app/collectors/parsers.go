package collectors

import (
	"cjavellana.me/launchpad/agent/app/stats"
	"cjavellana.me/launchpad/agent/app/utils/stringutils"
	"errors"
	"regexp"
	"strings"
	"time"
)

var patternNginxStatus = regexp.MustCompile("^Active connections: (?P<active>\\d+)" +
	"\\s*server accepts handled requests" +
	"\\s*(?P<accepts>\\d+) (?P<handled>\\d+) (?P<requests>\\d+)" +
	"\\s*Reading: (?P<reading>\\d+) Writing: (?P<writing>\\d+) Waiting: (?P<waiting>\\d+)\\s*$")

// Parse the standard nginx response status
// Active connections: 1
// server accepts handled requests
// 8 8 8
// Reading: 0 Writing: 1 Waiting: 0
func parseNginxStatus(responseBody string) (stats.WebServerStats, error) {
	result := patternNginxStatus.FindAllStringSubmatch(
		// strip newlines. Makes it easier to match pattern
		strings.Replace(responseBody, "\n", "", -1),
		-1,
		)

	if len(result) > 0 {
		return stats.WebServerStats{
			Timestamp:         time.Now(),
			ActiveConnections: stringutils.ToUint8(result[0][1]),
			Accepted:          stringutils.ToUint8(result[0][2]),
			Handled:           stringutils.ToUint8(result[0][3]),
			Requests:          stringutils.ToUint8(result[0][4]),
			Reading:           stringutils.ToUint8(result[0][5]),
			Writing:           stringutils.ToUint8(result[0][6]),
			Waiting:           stringutils.ToUint8(result[0][7]),
			RawData:           responseBody,
		}, nil
	}

	return stats.WebServerStats{}, errors.New("cannot parse response body")
}
