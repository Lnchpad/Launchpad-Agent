package system

import (
	"bufio"
	observer "cjavellana.me/launchpad/agent/app/observers"
	"errors"
	"fmt"
	"io"
)

type Stdout struct {
	Reader io.Reader

	observers []observer.TextObserver
	done      chan bool
}

type Process struct {
	PID    int
	Stdout Stdout
}

func (out *Stdout) StopObserving() error {
	if out.done == nil {
		return errors.New("nothing is being observed. are you sure you called observe()")
	}

	out.done <- true
	return nil
}

func (out *Stdout) Observe(observer observer.TextObserver) {
	if out.done == nil {
		out.done = make(chan bool)
	}

	out.observers = append(out.observers, observer)

	go func(reader io.Reader) {
		scanner := bufio.NewScanner(reader)

		for {
			select {
			case <-out.done:
				return
			default:
				if scanner.Scan() {
					for _, o := range out.observers {
						o.Update(fmt.Sprintf("%s\n", scanner.Text()))
					}
				}
			}
		}

	}(out.Reader)
}
