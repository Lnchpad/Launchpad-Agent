package main

import (
	"cjavellana.me/launchpad/agent/errors"
	"cjavellana.me/launchpad/agent/events/publishers"
	"cjavellana.me/launchpad/agent/messaging"
	"cjavellana.me/launchpad/agent/metrics"
	"cjavellana.me/launchpad/agent/servers/nginx"
	"cjavellana.me/launchpad/agent/system"
	"cjavellana.me/launchpad/agent/view"
	"cjavellana.me/launchpad/agent/view/widgets"
	"context"
	"flag"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/terminal/termbox"
	"log"
	"time"
)

type SystemProbes struct {
	cpu        *metrics.CpuProbe
	mem        *metrics.MemoryProbe
	serverStat *metrics.NginxProbe
}

type StartUpArgs struct {
	config string
}

func parseArgs() *StartUpArgs {
	var cfg string
	flag.StringVar(&cfg, "config", "config.yaml", "The configuration file to use")
	flag.Parse()

	return &StartUpArgs{
		config: cfg,
	}
}

func initializeProbes(serverCfg nginx.Server) SystemProbes {
	probes := SystemProbes{
		cpu: metrics.NewCpuProbe(1 * time.Second),
		mem: metrics.NewMemoryProbe(1 * time.Second),
		serverStat: metrics.NewNginxProbe(
			serverCfg.Monitor.MonitoringUrl,
			time.Duration(serverCfg.Monitor.PollingIntervalSecs)*time.Second,
			time.Duration(serverCfg.Monitor.InitialDelaySecs)*time.Second,
		),
	}

	// Start the Probes
	observers := []metrics.Probe{probes.cpu,
		probes.mem,
		probes.serverStat,
	}
	system.Observe(observers)
	return probes
}

func main() {
	args := parseArgs()

	appCfg := GetServerConfigFrom(args.config)
	server := appCfg.Server

	serverProcess, err := server.Start()
	errors.CheckFatal(err)

	// we start the tailer here that will Observe the nginx stdout
	// for Build content
	logsChannel := serverProcess.StdOut()

	probes := initializeProbes(server)

	// Get Kafka Producer
	broker := messaging.NewKafkaBroker(appCfg.Messaging)
	kafkaProducer := broker.NewProducer()

	stdoutWidget := widgets.NewRollContentDisplay()
	nginxMetricsWidget := widgets.NewRollContentDisplay()

	memMetricsWidget := widgets.NewLineChart(15, time.Second)
	memMetricsKafkaPublisher := publishers.NewMemMetricsPublisher(kafkaProducer)
	probes.mem.SubscribeMany([]metrics.Subscriber{
		memMetricsWidget,
		memMetricsKafkaPublisher,
	})

	cpuMetricsWidget := widgets.NewLineChart(15, time.Second)
	cpuMetricsKafkaPublisher := publishers.NewCpuMetricsPublisher(kafkaProducer)
	probes.cpu.SubscribeMany([]metrics.Subscriber{
		cpuMetricsWidget,
		cpuMetricsKafkaPublisher,
	})

	// Metrics Dispatcher
	metricsPublisher := publishers.NewMetricsPublisher(kafkaProducer)
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
	}(logsChannel, probes.serverStat)

	// The Dashboard View
	t, err := termbox.New()
	errors.CheckFatal(err)
	defer t.Close()

	dashboard := view.SimpleDashboardBuilder().
		WithCpuWidget(cpuMetricsWidget).
		WithMemoryWidget(memMetricsWidget).
		WithStdoutWidget(stdoutWidget.Display).
		WithNginxMetrics(nginxMetricsWidget.Display).
		Build(t)

	ctx, _ := context.WithCancel(context.Background())
	if err := termdash.Run(ctx, t, dashboard); err != nil {
		panic(err)
	}

}


