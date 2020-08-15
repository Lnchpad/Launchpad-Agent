package sync

import (
	"cjavellana.me/launchpad/agent/internal/pkg/cfg"
	"cjavellana.me/launchpad/agent/internal/pkg/servers"
	"log"
)

// The file system updater. This component is responsible for updating the
// root directory. This component ensures that there is only a single thread
// that updates the directory.

type Job struct {
	AppName string
}

type FsUpdater struct {
	cfg cfg.FsUpdaterConfig
	// channel used for sending job to this worker
	jobChannel chan Job
	webServer  servers.WebServer
}

func NewFsUpdater(cfg cfg.FsUpdaterConfig, webServer  servers.WebServer) FsUpdater {
	fs := FsUpdater{
		cfg:        cfg,
		jobChannel: make(chan Job, 100),
		webServer:  webServer,
	}

	// run the worker
	go fs.start()

	return fs
}

func (fs *FsUpdater) EnqueueJob(job Job) {
	fs.jobChannel <- job
}

func (fs *FsUpdater) start() {
	for job := range fs.jobChannel {

		err := fs.webServer.RegisterApp(cfg.PortalApp{
			AppName: job.AppName,
		})
		if err != nil {
			log.Println(err)
			return
		}

		err = fs.webServer.Reload()
		if err != nil {
			log.Fatalf("Unable to restart server %s %v", job.AppName, err)
		}
	}
}
