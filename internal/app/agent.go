package app

import (
	"cjavellana.me/launchpad/agent/internal/pkg/cfg"
	"cjavellana.me/launchpad/agent/internal/pkg/messaging"
	"cjavellana.me/launchpad/agent/internal/pkg/messaging/api"
	"cjavellana.me/launchpad/agent/internal/pkg/servers"
	"cjavellana.me/launchpad/agent/internal/pkg/stats/collectors"
	"cjavellana.me/launchpad/agent/internal/pkg/sync"
	"cjavellana.me/launchpad/agent/internal/pkg/system"
	"cjavellana.me/launchpad/agent/internal/pkg/view"
	"cjavellana.me/launchpad/agent/internal/pkg/view/dashboard"
	"context"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/terminal/termbox"
	"log"
	"time"
)

type Agent struct {
	AppCfg cfg.AppConfig

	// these are set on agent start
	Broker           api.Broker
	webServer        servers.WebServer
	webServerProcess *system.Process
	Probes           []system.Probe
	WebStats         system.WebStatsProbe
}

func NewAgent() Agent {
	return Agent{AppCfg: cfg.Get()}
}

func (agent *Agent) Start() {
	agent.initWebServer()
	agent.initMessageBroker()
	agent.initProbes()
	agent.initLogsAndStatsCollector()
	agent.initSyncListener()
	agent.initView()
}

func (agent *Agent) initSyncListener() {
	msgConsumer := agent.Broker.NewConsumer("gblevent")
	fsUpdater := sync.NewFsUpdater(agent.AppCfg.FsUpdaterConfig, agent.webServer)

	p := sync.NewPortalEventListener(msgConsumer, fsUpdater)
	p.StartListening()
}

func (agent *Agent) Terminate() {
	// stop server. Ignore error if already stopped
	_ = agent.webServer.Stop()
}

func (agent *Agent) initWebServer() {
	appCfg := agent.AppCfg

	webServer := servers.NewWebServer(appCfg.ServerConfig)
	if webServerProcess, err := webServer.Start(); err != nil {
		log.Fatalf("Unable to start webServer %v", err)
	} else {
		agent.webServer = webServer
		agent.webServerProcess = webServerProcess
	}

}

func (agent *Agent) initMessageBroker() {
	appCfg := agent.AppCfg
	broker := messaging.NewBroker(appCfg.BrokerConfig)
	agent.Broker = broker
}

// registers message producers to system probes for the purpose of
// sending them to central diagnostics collection system
func (agent *Agent) initLogsAndStatsCollector() {
	// cpu & memory utilization
	statsCollector := collectors.NewStatsCollector(
		agent.Broker.NewProducer("stats"))
	for _, probe := range agent.Probes {
		probe.Observe(&statsCollector)
	}

	// nginx logs
	logCollector := collectors.NewLogCollector(agent.Broker.NewProducer("logs"))
	agent.webServerProcess.Stdout.Observe(&logCollector)

	// nginx stats
	webStatsCollector := collectors.NewWebStatsCollector(agent.Broker.NewProducer("webstats"))
	agent.WebStats.Observe(&webStatsCollector)
}

func (agent *Agent) initProbes() {
	probeCfg := agent.AppCfg.ProbeConfig
	if probeCfg.Enabled {
		enabledProbes := probeCfg.ProbeTypes
		if len(enabledProbes) < 1 {
			log.Fatalf("no probes found")
		}

		for _, p := range probeCfg.ProbeTypes {
			probe := system.NewProbe(p, time.Duration(probeCfg.SamplingInterval)*time.Second)
			agent.Probes = append(agent.Probes, probe)
		}
	}

	agent.WebStats = system.NewWebStatsProbe(
		probeCfg.WebStatsConfig.StatsUrl,
		time.Duration(probeCfg.SamplingInterval)*time.Second,
		time.Duration(probeCfg.WebStatsConfig.InitialDelay)*time.Second,
	)
}

func (agent *Agent) initView() {
	viewType := agent.AppCfg.ViewType

	// Point server log to the address of Stdout otherwise,
	// serverLog will hold a copy of the `Stdout` structure
	// resulting to the dashboard not being able to stream server logs
	serverLog := &agent.webServerProcess.Stdout
	webStats := &agent.WebStats

	switch viewType {
	case cfg.ViewTypeNone:
		serverLog.Observe(&view.SimpleStdoutPrinter{})
		agent.StartReceivingLogs()
	case cfg.ViewTypeDashboardSimple:
		if len(agent.Probes) < 1 {
			log.Fatalf("unable to initialize simple display, no probes found")
		}

		t, err := termbox.New()
		if err != nil {
			log.Fatal(err)
		}
		defer t.Close()

		display := dashboard.NewSimpleDashboardBuilder(dashboard.SimpleDashboardConfig{
			AppCfg:    agent.AppCfg,
			Probes:    agent.Probes,
			ServerLog: serverLog,
			WebStats:  webStats,
		}).Build(t)

		agent.StartReceivingLogs()

		ctx, _ := context.WithCancel(context.Background())
		if err := termdash.Run(ctx, t, display); err != nil {
			panic(err)
		}
	}
}

func (agent *Agent) StartReceivingLogs() {
	agent.webServerProcess.Stdout.StartObserving()
}
