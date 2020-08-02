package system

import (
	"bufio"
	"fmt"
	"io"
)

//
type Process struct {
	Stdout io.Reader
	Err    error

	isTailActive    bool
	terminateFollow chan bool
}

func (s *Process) TailStdout(out chan string) {
	go func(reader io.Reader) {
		scanner := bufio.NewScanner(reader)

		for {
			select {
			case <-s.terminateFollow:
				return
			default:
				if scanner.Scan() {
					out <- fmt.Sprintf("%s\n", scanner.Text())
				}
			}
		}

	}(s.Stdout)
}

func (s *Process) StopTail() {
	if s.isTailActive {
		s.terminateFollow <- true
	}
}
