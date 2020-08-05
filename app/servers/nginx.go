package servers

import (
	"cjavellana.me/launchpad/agent/app/cfg"
	"cjavellana.me/launchpad/agent/app/system"
	"errors"
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

	serverCfg cfg.ServerConfig
}

func NewNginx(cfg cfg.ServerConfig) Nginx {
	return Nginx{serverCfg: cfg}
}

func (n *Nginx) Start() error {
	if err := n.createOrUpdateConfigFile(); err != nil {
		return err
	}

	// TODO: Check if we already have running process

	process := system.Process{}

	cmd := n.getCmd()
	if reader, err := cmd.StdoutPipe(); err != nil {
		return err
	} else {
		process.Stdout = system.Stdout{
			Reader: reader,
		}
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

func (n *Nginx) createOrUpdateConfigFile() error {
	return ioutil.WriteFile(n.serverCfg.ConfigLocation, []byte(n.serverCfg.ConfigTemplate), 0644)
}

func (n *Nginx) getCmd() *exec.Cmd {
	serverCfg := n.serverCfg

	switch serverCfg.ExecutablePath {
	case "local":
		return exec.Command("tail", "-f", "/tmp/sample.txt")
	default:
		return exec.Command(serverCfg.ExecutablePath, "-c", serverCfg.ConfigLocation)
	}
}
