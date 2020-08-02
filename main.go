package main

import (
	"cjavellana.me/launchpad/agent/errors"
	"cjavellana.me/launchpad/agent/events/publishers"
	"cjavellana.me/launchpad/agent/messaging"
	"cjavellana.me/launchpad/agent/metrics"
	"cjavellana.me/launchpad/agent/view"
	"cjavellana.me/launchpad/agent/view/widgets"
	"context"
	"flag"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/terminal/termbox"
	"log"
	"time"
)

type StartUpArgs struct {
	config string
}

func main() {
	args := parseArgs()

	appCfg := GetServerConfigFrom(args.config)
	server := &appCfg.Server

	serverProcess, err := server.Start()
	errors.CheckFatal(err)

	ctx, _ := context.WithCancel(context.Background())

	// we start the tailer here that will Observe the nginx stdout
	// for Build content
	logsChannel := serverProcess.StdOut()

	cpuProbe := metrics.NewCpuProbe(1 * time.Second)
	memProbe := metrics.NewMemoryProbe(1 * time.Second)
	nginxProbe := metrics.NewNginxProbe(
		server.Monitor.MonitoringUrl,
		time.Duration(server.Monitor.PollingIntervalSecs)*time.Second,
		time.Duration(server.Monitor.InitialDelaySecs)*time.Second,
	)

	observers := []metrics.Probe{cpuProbe, memProbe, nginxProbe}
	metrics.Observe(observers)

	// Get Kafka Producer
	broker := messaging.NewKafkaBroker(appCfg.Messaging)
	kafkaProducer := broker.NewProducer()
	metricsPublisher := publishers.NewMetricsPublisher(kafkaProducer)

	// The Dashboard View
	t, err := termbox.New()
	errors.CheckFatal(err)
	defer t.Close()

	stdoutWidget := widgets.NewRollContentDisplay()
	nginxMetricsWidget := widgets.NewRollContentDisplay()

	memMetricsWidget := widgets.NewLineChart(15, time.Second)
	memMetricsKafkaPublisher := publishers.NewMemMetricsPublisher(kafkaProducer)
	memProbe.SubscribeMany([]metrics.Subscriber{
		memMetricsWidget,
		memMetricsKafkaPublisher,
	})

	cpuMetricsWidget := widgets.NewLineChart(15, time.Second)
	cpuMetricsKafkaPublisher := publishers.NewCpuMetricsPublisher(kafkaProducer)
	cpuProbe.SubscribeMany([]metrics.Subscriber{
		cpuMetricsWidget,
		cpuMetricsKafkaPublisher,
	})

	// Metrics Dispatcher
	go func(out chan string, n *metrics.NginxProbe) {
		for {
			select {
			case logs := <-out:
				stdoutWidget.Update(logs)

				if err := metricsPublisher.PublishServerLogs(logs); err != nil {
					log.Println(err)
				}
			case stat := <-n.StatsChannel:
				message := stat.RawData
				if stat.ErrorMessage != "" {
					message = stat.ErrorMessage
				}
				nginxMetricsWidget.Update(message)
			}
		}
	}(logsChannel, nginxProbe)

	dashboard := view.SimpleDashboardBuilder().
		WithCpuWidget(cpuMetricsWidget.LineChart).
		WithMemoryWidget(memMetricsWidget.LineChart).
		WithStdoutWidget(stdoutWidget.Display).
		WithNginxMetrics(nginxMetricsWidget.Display).
		Build(t)

	if err := termdash.Run(ctx, t, dashboard); err != nil {
		panic(err)
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
