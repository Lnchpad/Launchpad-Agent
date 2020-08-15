package servers

import (
	"cjavellana.me/launchpad/agent/app/cfg"
	"cjavellana.me/launchpad/agent/app/system"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
)

type Status int

const (
	Running Status = iota
	Stopped
)

type Nginx struct {
	// The process id
	Process system.Process

	Status Status

	ServerCfg cfg.ServerConfig
}

func NewNginx(cfg cfg.ServerConfig) Nginx {
	return Nginx{ServerCfg: cfg}
}

// returns an error if `p.appName` already exist
func (n *Nginx) RegisterApp(p cfg.PortalApp) error {
	for _, app := range n.ServerCfg.Applications {
		if app.AppName == p.AppName {
			return errors.New(fmt.Sprintf("App %s already exists", p.AppName))
		}
	}

	return nil
}

func (n *Nginx) Start() error {
	err := n.rebuildAndUpdateConfig()
	if err != nil {
		return err
	}

	// TODO: Check if we already have running process

	process := system.Process{}

	cmd := n.getCmd()
	if reader, err := cmd.StdoutPipe(); err != nil {
		return err
	} else {
		process.Stdout = system.NewStdout(reader)
	}

	if err := cmd.Start(); err != nil {
		return err
	} else {
		// Fixme: Get PID of the newly spawned process
		process.PID = 0

		n.Status = Running
	}

	n.Process = process

	// no error
	return nil
}

func (n *Nginx) Stop() error {
	if n.Status != Running {
		return errors.New("cannot stop a stopped process")
	}

	n.Process.PID = -1
	n.Status = Stopped

	return nil
}

func (n *Nginx) Restart() error {
	return nil
}

func (n *Nginx) rebuildAndUpdateConfig() error {
	builder := newConfigBuilder(n.ServerCfg)
	configString, _ := builder.Build()
	return writeToFile(n.ServerCfg.ConfigLocation, configString)
}

func writeToFile(fileLocation string, config string) error {
	return ioutil.WriteFile(fileLocation, []byte(config), 0644)
}

func (n *Nginx) getCmd() *exec.Cmd {
	serverCfg := n.ServerCfg

	switch serverCfg.ExecutablePath {
	case "local":
		return exec.Command("tail", "-f", "/tmp/sample.txt")
	default:
		return exec.Command(serverCfg.ExecutablePath, "-c", serverCfg.ConfigLocation)
	}
}
