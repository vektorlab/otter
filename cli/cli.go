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
	stateMap     state.StateMap
	statesAsJson []byte
	Run          func() error
}

func NewOtterCLI(command, statePath string, EtcdURL []string) (*OtterCLI, error) {

	var err error

	cli := OtterCLI{}

	cli.stateMap, err = state.StateMapFromYamlPath(statePath)

	if err != nil {
		return nil, err
	}

	cli.otter, err = clients.NewOtterClient(EtcdURL)

	if err != nil {
		return nil, err
	}

	cli.etcdURL = EtcdURL

	switch command {
	case "daemon":
		cli.Run = cli.daemon
	case "execute":
		cli.Run = cli.execute
	case "load":
		cli.Run = cli.load
	case "ls":
		cli.Run = cli.ls
	default:
		return nil, fmt.Errorf("Unknown command: %s", command)
	}

	return &cli, nil
}

/*
Run Otter as a daemon (slave-mode)
*/
func (cli *OtterCLI) daemon() error {
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
func (cli *OtterCLI) execute() error {
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
func (cli *OtterCLI) load() error {
	raw, err := cli.stateMap.ToJson()
	if err != nil {
		return err
	}
	return cli.otter.SubmitState(string(raw))
}

/*
List remote hosts
*/
func (cli *OtterCLI) ls() error {
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
