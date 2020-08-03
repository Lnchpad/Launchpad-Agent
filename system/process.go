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

func (p *Process) StdOut() chan string {
	out := make(chan string)

	go func(reader io.Reader) {
		scanner := bufio.NewScanner(reader)

		for {
			select {
			case <-p.terminateFollow:
				return
			default:
				if scanner.Scan() {
					out <- fmt.Sprintf("%s\n", scanner.Text())
				}
			}
		}

	}(p.Stdout)

	return out
}

func (p *Process) StopTail() {
	if p.isTailActive {
		p.terminateFollow <- true
	}
}
