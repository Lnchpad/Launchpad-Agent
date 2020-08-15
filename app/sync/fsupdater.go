package sync

import (
	"cjavellana.me/launchpad/agent/app/cfg"
	"cjavellana.me/launchpad/agent/app/servers"
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
	nginx      *servers.Nginx
}

func NewFsUpdater(cfg cfg.FsUpdaterConfig, nginx *servers.Nginx) FsUpdater {
	fs := FsUpdater{
		cfg:        cfg,
		jobChannel: make(chan Job, 100),
		nginx:      nginx,
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

		err := fs.nginx.RegisterApp(cfg.PortalApp{
			AppName: job.AppName,
		})
		if err != nil {
			log.Println(err)
			return
		}

		if err := fs.nginx.Restart(); err != nil {
			// Error restarting nginx.
			log.Fatalf("Unable to restart server %s %v", job.AppName, err)
		}
	}
}
