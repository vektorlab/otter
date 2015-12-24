package daemon

import (
	"github.com/vektorlab/otter/client"
	"github.com/vektorlab/otter/state"
	log "github.com/Sirupsen/logrus"
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
	switch command {

	case "state":

		log.Printf("Processing command %s (%s)", command, id)

		states, err := daemon.otter.RetrieveStateMap()

		if err != nil {
			return err
		}

		results := make([]state.Result, 0)

		for _, value := range states { // TODO: Dependency Processing
			for _, state := range value {
				result, err := daemon.otter.CheckConsistent(state)
				if err != nil {
					return err
				}
				results = append(results, result)
			}
		}
		err = daemon.otter.SaveResults(id, results)
		if err != nil {
			return err
		}

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
