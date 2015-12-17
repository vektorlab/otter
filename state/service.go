/*
A service represents a running daemonized process on an operating system

States -
  running: The service is running // TODO
  stopped: The service is stopped // TODO
*/
package state

import (
	"encoding/json"
	"fmt"
	"github.com/vektorlab/otter/helpers"
)

type Service struct {
	Name     string   `json:"name"`
	Running  bool     `json:"running"`
	Metadata Metadata `json:"metadata"`
	Require  []string `json:"require"`
}

func (service *Service) Initialize() error {
	if service.Name == "" {
		service.Name = service.Metadata.Name
	}
	state := service.Metadata.State
	if !(state == "running" || state == "stopped") {
		return fmt.Errorf("Invalid service state: %s", state)
	}
	return nil
}

func (service *Service) Consistent() (bool, error) {
	running, err := helpers.ServiceRunning(service.Name)
	if err != nil {
		return false, err
	}
	return service.Running == running, nil
}

func (service *Service) Execute() error {
	consistent, err := service.Consistent()
	if err != nil {
		return err
	}
	if !consistent {
		err := helpers.ChangeServiceState(service.Name, service.Running)
		if err != nil {
			return err
		}
	}
	return nil
}

func (service *Service) Dump() ([]byte, error) {
	return json.Marshal(service)
}

func (service *Service) Requirements() []string {
	return service.Require
}

func (service *Service) Meta() Metadata {
	return service.Metadata
}
