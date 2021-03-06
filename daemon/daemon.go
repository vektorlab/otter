package daemon

import (
	log "github.com/Sirupsen/logrus"
	"github.com/vektorlab/otter/client"
	"github.com/vektorlab/otter/helpers"
	"os"
	"time"
)

type Daemon struct {
	otter    *clients.Otter
	last     string
	firstRun bool
}

func (daemon *Daemon) register() error {
	last, err := daemon.otter.UpdateHost(daemon.otter.Hostname, daemon.last)
	if err != nil {
		return err
	}
	daemon.last = last
	return nil
}

func (daemon *Daemon) synchronize() {
	err := daemon.register()
	if err != nil {
		log.Fatalf(err.Error())
	}
	time.Sleep(15 * time.Second)
	daemon.synchronize()
}

func (daemon *Daemon) listen() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	command, id, err := daemon.otter.WaitForCommand(hostname)
	if err != nil {
		log.Fatal(err)
	}
	err = daemon.ProcessCommand(command, id)
	if err != nil {
		log.Fatal(err)
	}
	daemon.listen()
}

func (daemon *Daemon) ProcessCommand(command, id string) error {
	log.Printf("Processing command %s (%s)", command, id)
	switch command {
	case "apply":
		stateMap, err := daemon.otter.RetrieveStateMap()
		helpers.FailOnError(err, "Unable to retrieve StateMap from Etcd")
		return daemon.otter.SaveResultMap(id, stateMap.Apply())
	case "state":
		stateMap, err := daemon.otter.RetrieveStateMap()
		helpers.FailOnError(err, "Unable to retrieve StateMap from Etcd")
		return daemon.otter.SaveResultMap(id, stateMap.State())
	default:
		log.Fatalf("Unknown command: %s (%s)", command, id)
	}
	return nil
}

func (daemon *Daemon) Run() {
	go daemon.synchronize()
	go daemon.listen()
	select {}
}

func NewDaemon(servers []string) (*Daemon, error) {
	var err error
	daemon := Daemon{
		firstRun: true,
	}
	daemon.otter, err = clients.NewOtterClient(servers)
	if err != nil {
		return nil, err
	}
	return &daemon, nil
}
