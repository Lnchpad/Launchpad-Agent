package stringutils

import (
	"log"
	"strconv"
)

func ToUint8(s string) uint32 {
	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		log.Print(err)
		return 0
	}

	return uint32(v)
}
