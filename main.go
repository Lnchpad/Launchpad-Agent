package main

import (
	"cjavellana.me/launchpad/agent/errors"
	"cjavellana.me/launchpad/agent/metrics"
	"cjavellana.me/launchpad/agent/os"
	"cjavellana.me/launchpad/agent/servers/nginx"
	"context"
	"flag"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/terminal/termbox"
	"log"
	"time"
)

func startServer(n *nginx.Server, ch chan<- os.Process) {
	ch <- n.Start()
}

type StartUpArgs struct {
	config string
}

func main() {
	args := parseArgs()

	server := GetServerConfigFrom(args.config)

	serverStartChannel := make(chan os.Process)
	go startServer(server, serverStartChannel)

	if stdout := <-serverStartChannel; stdout.Err != nil {
		log.Fatalf("Unable to start server: %s\n", stdout.Err)
	} else {
		ctx, _ := context.WithCancel(context.Background())

		// we start the tailer here that will Observe the nginx stdout
		// for build content
		outChannel := make(chan string)
		stdout.TailStdout(outChannel)

		cpuObserver := metrics.NewCpuObserver(1 * time.Second)
		memObserver := metrics.NewMemoryObserver(1 * time.Second)
		nginxObserver := nginx.NewStatusObserver(
			server.Monitor.MonitoringUrl,
			time.Duration(server.Monitor.PollingIntervalSecs) * time.Second,
			time.Duration(server.Monitor.InitialDelaySecs) * time.Second,
			)

		observers := []metrics.Observable{cpuObserver, memObserver, nginxObserver}
		metrics.Observe(observers)

		// The Dashboard View
		t, err := termbox.New()
		errors.CheckFatal(err)
		defer t.Close()

		stdoutWidget := NewServerStdoutWindow(outChannel)
		nginxMetricsWidget := NewNginxStatusWindow(nginxObserver.Channel)
		cpuMetricsWidget := NewLineChart(cpuObserver.Channel, 15, time.Second)
		memoryMetricsWidget := NewLineChart(memObserver.MetricsChannel, 15, time.Second)

		dashboard := StandardDashboardBuilder(t).
			WithCpuWidget(cpuMetricsWidget).
			WithMemoryWidget(memoryMetricsWidget).
			WithStdoutWidget(stdoutWidget).
			WithNginxMetrics(nginxMetricsWidget).
			build()

		if err := termdash.Run(ctx, t, dashboard); err != nil {
			panic(err)
		}
	}
}

func parseArgs() *StartUpArgs {
	var cfg string
	flag.StringVar(&cfg, "config", "config.yaml", "The configuration file to use")
	flag.Parse()

	return &StartUpArgs{
		config: cfg,
	}
}