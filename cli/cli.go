package cli

import (
	"fmt"
	"github.com/vektorlab/otter/client"
	"github.com/vektorlab/otter/daemon"
	"github.com/vektorlab/otter/state"
)

type OtterCLI struct {
	otter        *clients.Otter
	etcdURL      []string
	states       map[string][]state.State
	statesAsJson []byte
	Run          func() error
}

func NewOtterCLI(command, statePath string, EtcdURL []string) (*OtterCLI, error) {

	var err error

	cli := OtterCLI{}

	cli.states, err = state.StatesFromYamlPath(statePath)

	if err != nil {
		return nil, err
	}

	cli.otter, err = clients.NewOtterClient(EtcdURL)

	if err != nil {
		return nil, err
	}

	cli.etcdURL = EtcdURL

	cli.statesAsJson, err = state.StatesToJson(cli.states)

	if err != nil {
		return nil, err
	}

	switch command {
	case "daemon":
		cli.Run = cli.Daemon
	case "execute":
		cli.Run = cli.Execute
	case "load":
		cli.Run = cli.Load
	case "ls":
		cli.Run = cli.Ls
	default:
		return nil, fmt.Errorf("Unknown command: %s", command)
	}

	return &cli, nil
}

/*
Run Otter as a daemon (slave-mode)
*/
func (cli *OtterCLI) Daemon() error {
	daemon, err := daemon.NewDaemon(cli.etcdURL)
	if err != nil {
		return err
	}
	daemon.Run()
	return nil
}

/*
Execute a remote command against slaves
*/
func (cli *OtterCLI) Execute() error {
	hosts, err := cli.otter.ListHosts()

	if err != nil {
		return err
	}

	for _, host := range hosts {
		results, err := cli.otter.SubmitCommand(host, "execute")
		if err != nil {
			return err
		}
		DumpResults(results)
	}

	return nil
}

/*
Load a new state definition
*/
func (cli *OtterCLI) Load() error {
	return cli.otter.SubmitState(string(cli.statesAsJson))
}

/*
List remote hosts
*/
func (cli *OtterCLI) Ls() error {
	hosts, err := cli.otter.ListHosts()

	if err != nil {
		return err
	}

	hostList := make(map[string]bool)

	for _, host := range hosts {
		results, err := cli.otter.SubmitCommand(host, "state")
		if err != nil {
			return err
		}
		hostList[host] = isConsistent(results)
	}

	DumpHosts(hostList)

	return nil
}
