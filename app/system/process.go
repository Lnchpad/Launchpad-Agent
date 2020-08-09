package system

import (
	"bufio"
	"fmt"
	"github.com/matryer/runner"
	"io"
	"time"
)

type Stdout struct {
	Reader io.ReadCloser

	observers []TextObserver
	done      chan struct{}
	status    Status

	pollTask *runner.Task
}

type Process struct {
	PID    int
	Stdout Stdout
}

func NewStdout(reader io.ReadCloser) Stdout {
	return Stdout{
		Reader: reader,
		status: Stopped,
		done:   make(chan struct{}),
	}
}

func (out *Stdout) StopObserving() {
	close(out.done)
	out.pollTask.Stop()

	select {
	case <-out.pollTask.StopChan():
		// task successfully stopped
	case <-time.After(2 * time.Second):
		// task didn't stop in time
	}
}

func (out *Stdout) StartObserving() {
	out.done = make(chan struct{})
	out.pollTask = runner.Go(func(shouldStop runner.S) error {
		scanner := bufio.NewScanner(out.Reader)
		for {
			select {
			case <-out.done:
				break
			default:
				if scanner.Scan() {
					for _, o := range out.observers {
						o.Update(fmt.Sprintf("%s\n", scanner.Text()))
					}
				}
			}
		}
	})
}

func (out *Stdout) Observe(observer TextObserver) {
	out.observers = append(out.observers, observer)

	if out.pollTask != nil && out.pollTask.Running() {
		out.StopObserving()
	}
}
