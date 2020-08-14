package sync

import "fmt"

// The file system updater. This component is responsible for updating the
// root directory. This component ensures that there is only a single thread
// that updates the directory.

type Job struct {
	AppName string
}

type FsUpdaterConfig struct {
	// the location of all the static files in the
	// local file system i.e. in the container
	RootDirectory string

	// the url of the nexus repository
	NexusUrl string
}

type FsUpdater struct {
	cfg FsUpdaterConfig

	jobChannel chan Job
}

func NewFsUpdater(cfg FsUpdaterConfig) FsUpdater {
	fs := FsUpdater{
		cfg:        cfg,
		jobChannel: make(chan Job, 100),
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
		fmt.Println(job.AppName)
	}
}
