package servers

import (
	"cjavellana.me/launchpad/agent/internal/pkg/cfg"
	"cjavellana.me/launchpad/agent/internal/pkg/servers/nginx"
	"cjavellana.me/launchpad/agent/internal/pkg/system"
)

type WebServer interface {
	Start() (*system.Process, error)

	Stop() error

	Reload() error

	RegisterApp(app cfg.PortalApp) error
}

func NewWebServer(config cfg.ServerConfig) WebServer {
	ws := nginx.NewNginx(config)
	return ws
}
