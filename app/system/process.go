package system

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

type Stdout struct {
	Reader io.Reader

	observers []TextObserver
	done      chan bool
	status    Status
}

type Process struct {
	PID    int
	Stdout Stdout
}

func NewStdout(reader io.Reader) Stdout {
	return Stdout{
		Reader: reader,
		status: Stopped,
	}
}

func (out *Stdout) StopObserving() error {
	if out.done == nil {
		return errors.New("nothing is being observed. are you sure you called observe()")
	}

	out.status = Stopped
	out.done <- true
	return nil
}

func (out *Stdout) Observe(observer TextObserver) {
	if out.done == nil {
		out.done = make(chan bool)
	}

	out.observers = append(out.observers, observer)

	if out.status == Stopped {
		out.status = Running
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
}
