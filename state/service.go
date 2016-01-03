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

func (service *Service) Consistent() *Result {
	result := &Result{
		Metadata: &service.Metadata,
		Consistent: false,
	}
	running, err := helpers.ServiceRunning(service.Name)
	if err != nil {
		result.Message = err.Error()
		return result
	}
	result.Consistent = service.Running == running
	return result
}

func (service *Service) Execute() *Result {
	result := &Result{
		Metadata: &service.Metadata,
		Consistent: service.Consistent().Consistent,
	}
	if result.Consistent == false {
		err := helpers.ChangeServiceState(service.Name, service.Running)
		if err != nil {
			result.Message = err.Error()
			return result
		}
		result.Consistent = service.Consistent().Consistent
	}
	return result
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
