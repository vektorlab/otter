package daemon

import (
	"fmt"
	//"github.com/vektorlab/otter/state"
	"github.com/coreos/etcd/client"
	"github.com/vektorlab/otter/clients"
	"log"
	"os"
	"time"
)

type Daemon struct {
	conn     client.KeysAPI
	last     string
	firstRun bool
}

func (daemon *Daemon) register() error {

	hostname, err := os.Hostname()

	if err != nil {
		log.Fatalf(err.Error())
	}

	last, err := clients.Ping(daemon.conn, hostname, daemon.last)

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

func (daemon *Daemon) monitorEvents() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	change, err := clients.WaitForStateChange(daemon.conn, hostname)
	if err != nil {
		panic(err)
	}
	fmt.Println(change.Action, change.Index, change.Node, change.PrevNode)
	daemon.monitorEvents()
	/*
		response, err := daemon.etcd.Get(context.Background(), fmt.Sprintf("/hosts/%s/state", hostname), nil)
		if err != nil {
			log.Println("Etcd error: ", err.Error())
		} else {
			states, err := state.StatesFromProcessedJson([]byte(response.Node.Value))
			if err != nil {
				panic(err)
			}
			err = state.Execute(states)
			if err != nil {
				log.Println("Caught exception applying state: ", err.Error())
			}
		}
	*/
}

func (daemon *Daemon) Run() {
	go daemon.synchronize()
	go daemon.monitorEvents()
	select {}
}

func NewDaemon(servers []string) (*Daemon, error) {
	var err error

	daemon := Daemon{
		firstRun: true,
	}

	daemon.conn, err = clients.NewKeysApi(servers)

	if err != nil {
		return nil, err
	}

	return &daemon, nil
}
