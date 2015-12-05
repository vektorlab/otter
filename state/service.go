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
	"github.com/mitchellh/mapstructure"
	"github.com/vektorlab/otter/helpers"
)

type Service struct {
	Name     string   `mapstructure:"name"`
	Running  bool     `mapstructure:"running"`
	Metadata Metadata `mapstructure:"-"`
	Require  []string `mapstructure:"require"`
}

func (service *Service) Initialize() error {
	if service.Name == "" {
		service.Name = service.Metadata.Name
	}
	state := service.Metadata.State
	if ! (state == "running" || state == "stopped") {
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

func ServiceFromStructure(metadata Metadata, structure interface{}) (*Service, error) {
	var err error

	service := Service{
		Metadata: metadata,
	}

	err = mapstructure.Decode(structure, &service)
	if err != nil {
		return nil, err
	}

	err = service.Initialize()
	if err != nil {
		return nil, err
	}

	return &service, nil
}
