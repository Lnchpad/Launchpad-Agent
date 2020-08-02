package nginx

import (
	"cjavellana.me/launchpad/agent/errors"
	"cjavellana.me/launchpad/agent/models"
	"cjavellana.me/launchpad/agent/system"
	"fmt"
	"io/ioutil"
	"os/exec"
)

type Config struct {
	Template string

	// the location are where pages are served
	RootDir      string
	Location     string
	Applications []models.PortalApplication
}

type Monitor struct {
	InitialDelaySecs    int
	PollingIntervalSecs int
	MonitoringUrl       string
}

type Server struct {
	Config    *Config
	NginxPath string
	Monitor   *Monitor
}

func (n *Server) Start() (system.Process, error) {
	fmt.Println("Starting nginx server")

	err := n.createOrUpdateConfigFile()
	errors.CheckFatal(err)

	cmd := n.getCmd()
	reader, err := cmd.StdoutPipe()
	errors.CheckFatal(err)

	return system.Process{Stdout: reader}, cmd.Start()
}

func (n *Server) Stop() {
}

func (n *Server) Reload() system.Process {
	cmd := exec.Command(n.NginxPath, "-s", "reload")
	reader, err := cmd.StdoutPipe()
	errors.CheckFatal(err)

	return system.Process{Stdout: reader, Err: cmd.Start()}
}

func (n *Server) Status() string {
	return ""
}

func (n *Server) createOrUpdateConfigFile() error {
	return ioutil.WriteFile(n.Config.Location, []byte(n.Config.Template), 0644)
}

func (n *Server) getCmd() *exec.Cmd {
	switch n.NginxPath {
	case "local":
		return exec.Command("tail", "-f", "/tmp/sample.txt")
	default:
		return exec.Command(n.NginxPath, "-c", n.Config.Location)
	}
}
