package app

import (
	"cjavellana.me/launchpad/agent/app/cfg"
	"cjavellana.me/launchpad/agent/app/messaging"
	"cjavellana.me/launchpad/agent/app/servers"
	"cjavellana.me/launchpad/agent/app/system"
	"cjavellana.me/launchpad/agent/app/view"
	"cjavellana.me/launchpad/agent/app/view/dashboard"
	"context"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/terminal/termbox"
	"log"
	"time"
)

type Agent struct {
	AppCfg cfg.AppConfig

	// these are set on agent start
	Broker messaging.Broker
	Nginx  servers.Nginx
	Probes []system.Probe
}

func NewAgent() Agent {
	return Agent{AppCfg: cfg.Get()}
}

func (agent *Agent) Start() {
	agent.initWebServer()
	agent.initMessageBroker()
	agent.initProbes()
	agent.initDiagnosticsCollector()
	agent.initView()
}

func (agent *Agent) Terminate() {
	// stop server. Ignore error if already stopped
	_ = agent.Nginx.Stop()
}

func (agent *Agent) initWebServer() {
	appCfg := agent.AppCfg

	nginx := servers.NewNginx(appCfg.ServerConfig)
	if err := nginx.Start(); err != nil {
		log.Fatalf("Unable to start Nginx %v", err)
	}

	agent.Nginx = nginx
}

func (agent *Agent) initMessageBroker() {
	appCfg := agent.AppCfg
	broker := messaging.NewBroker(appCfg.BrokerConfig)
	agent.Broker = broker
}

// registers message producers to system probes for the purpose of
// sending them to central diagnostics collection system
func (agent *Agent) initDiagnosticsCollector() {
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
}

func (agent *Agent) initView() {
	viewType := agent.AppCfg.ViewType
	serverLog := agent.Nginx.Process.Stdout

	switch viewType {
	case cfg.ViewTypeNone:
		serverLog.Observe(&view.SimpleStdoutPrinter{})
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
			ServerLog: &serverLog,
		}).Build(t)

		ctx, _ := context.WithCancel(context.Background())
		if err := termdash.Run(ctx, t, display); err != nil {
			panic(err)
		}
	}
}
