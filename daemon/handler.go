package daemon

import (
	log "github.com/Sirupsen/logrus"
	"github.com/vektorlab/otter/state"
)

func (daemon *Daemon) GetState(id string) error {
	stateMap, err := daemon.otter.RetrieveStateMap()
	if err != nil {
		return err
	}
	results := make([]state.Result, 0)
	for _, value := range stateMap.States { // TODO: Dependency Processing
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
	return nil
}

func (daemon *Daemon) ApplyState(id string) error {
	stateMap, err := daemon.otter.RetrieveStateMap()
	if err != nil {
		return err
	}
	for _, value := range stateMap.States {
		for _, state := range value {
			err := state.Execute()
			if err != nil {
				log.Println("WARNING: ", err.Error())
			}
		}
	}
	return daemon.GetState(id)
}
