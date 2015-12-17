package daemon

import (
	"fmt"
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"log"
	"os"
	"strings"
	"time"
)

type Daemon struct {
	etcd     client.KeysAPI
	index    uint64
	firstRun bool
}

func (daemon *Daemon) register() error {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf(err.Error())
	}
	options := client.SetOptions{
		TTL: 60 * time.Second,
	}
	response, err := daemon.etcd.Set(context.Background(), fmt.Sprintf("/hosts/%s/host", hostname), hostname, &options)
	if err != nil {
		return err
	}

	if daemon.firstRun {
		daemon.index = response.Index
		daemon.firstRun = false
	} else {
		if daemon.index+1 != response.Index {
			return fmt.Errorf("Another host is already registered with this hostname")
		} else {
			daemon.index = response.Index
		}
	}
	return nil
}

func (daemon *Daemon) synchronize() {
	err := daemon.register()
	if err != nil {
		log.Fatalf(err.Error())
	}
	time.Sleep(15 * time.Second)
	daemon.Run()
}

func (daemon *Daemon) maintainState() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	response, err := daemon.etcd.Get(context.Background(), fmt.Sprintf("/hosts/%s/state", hostname), nil)
	if err != nil {
		log.Println("Etcd error: ", err.Error())
	} else {
		fmt.Println(response.Node.Value)
	}
}

func (daemon *Daemon) LoadState(state string) error {

	hosts, err := daemon.ListHosts()

	if err != nil {
		return err
	}

	for _, host := range hosts {
		response, err := daemon.etcd.Set(context.Background(), fmt.Sprintf("/hosts/%s/state", host), state, &client.SetOptions{})
		if err != nil {
			return err
		}
		log.Printf("Set host: %s %s", response.Action, host)
	}

	return nil
}

func (daemon *Daemon) ListHosts() ([]string, error) {

	var hosts []string

	response, err := daemon.etcd.Get(context.Background(), "/hosts", &client.GetOptions{Recursive: true})

	if err != nil {
		return hosts, err
	}

	for _, node := range response.Node.Nodes {
		hosts = append(hosts, strings.Split(node.Key, "/")[2])
	}

	return hosts, nil
}

func (daemon *Daemon) Run() {
	go daemon.synchronize()
	go daemon.maintainState()
	select {}
}

func NewDaemon(servers []string) (*Daemon, error) {
	daemon := Daemon{
		firstRun: true,
	}

	cfg := client.Config{
		Endpoints: servers,
		Transport: client.DefaultTransport,
	}

	c, err := client.New(cfg)

	if err != nil {
		return nil, err
	}

	daemon.etcd = client.NewKeysAPI(c)

	return &daemon, nil
}
