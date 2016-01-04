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
	stateMap     *state.StateMap
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
	case "apply":
		cli.Run = cli.apply
	case "daemon":
		cli.Run = cli.daemon
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
Apply the loaded state to hosts
*/
func (cli *OtterCLI) apply() error {
	resultMap, err := cli.otter.SubmitCommands("*", "apply")
	if err != nil {
		return err
	}
	DumpResults(resultMap)
	return nil
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
	resultMap, err := cli.otter.SubmitCommands("*", "state")
	if err != nil {
		return err
	}
	DumpResults(resultMap)
	return nil
}
