package nginx

import (
	"cjavellana.me/launchpad/agent/utils/stringutils"
	"fmt"
	"regexp"
	"strings"
)

var standardStatusPattern = regexp.MustCompile("^Active connections: (?P<active>\\d+)" +
"\\s*server accepts handled requests" +
"\\s*(?P<accepts>\\d+) (?P<handled>\\d+) (?P<requests>\\d+)" +
"\\s*Reading: (?P<reading>\\d+) Writing: (?P<writing>\\d+) Waiting: (?P<waiting>\\d+)\\s*$")

// Parse the standard nginx response status
// Active connections: 1
// server accepts handled requests
// 8 8 8
// Reading: 0 Writing: 1 Waiting: 0
func parseStandardStatus(responseBody string) StandardStatus {
	responseBodyWithoutNewLines := strings.Replace(responseBody, "\n", "", -1)
	result := standardStatusPattern.FindAllStringSubmatch(responseBodyWithoutNewLines, -1)

	if len(result) > 0 {
		return StandardStatus{
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

	return StandardStatus{ErrorMessage: fmt.Sprintf("Unable to parse\n%s", responseBody)}
}


