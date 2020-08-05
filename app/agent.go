package app

import (
	"cjavellana.me/launchpad/agent/app/cfg"
	"cjavellana.me/launchpad/agent/app/messaging"
	"cjavellana.me/launchpad/agent/app/servers"
	"cjavellana.me/launchpad/agent/app/view"
	"log"
)

type Agent struct {
	AppCfg cfg.AppConfig

	// these are set on agent start
	Broker messaging.Broker
	Nginx  servers.Nginx
}

func NewAgent() Agent {
	return Agent{AppCfg: cfg.Get()}
}

func (agent *Agent) Start() {
	appCfg := agent.AppCfg

	broker := messaging.NewBroker(appCfg.BrokerConfig)
	nginx := servers.NewNginx(appCfg.ServerConfig)

	if err := nginx.Start(); err != nil {
		log.Fatalf("Unable to start Nginx %v", err)
	}

	agent.Nginx = nginx
	agent.Broker = broker

	agent.initProbes()
	agent.initView()
}

func (agent *Agent) Terminate() {
	// stop server. Ignore error if already stopped
	_ = agent.Nginx.Stop()
}

func (agent *Agent) initProbes() {}

func (agent *Agent) initView() {
	viewType := agent.AppCfg.ViewType

	switch viewType {
	case cfg.ViewTypeNone:
		agent.Nginx.Process.Stdout.Observe(&view.SimpleStdoutPrinter{})
	}
}
