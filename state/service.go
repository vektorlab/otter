/*
A service represents a running daemonized process on an operating system

States -
  running: The service is running // TODO
  stopped: The service is stopped // TODO
*/
package state

import (
	"fmt"
	"encoding/json"
	"github.com/mitchellh/mapstructure"
)

type Service struct {
	Name     string   `mapstructure:"name"`
	Running  bool     `mapstructure:"running"`
	Metadata Metadata `mapstructure:"-"`
	Require  []string `mapstructure:"require"`
}

func (service *Service) validateState() error {
	switch service.Metadata.State {
	case "running":
		return nil
	case "stopped":
		return nil
	default:
		return fmt.Errorf("Service state %s should be one of 'running' or 'stopped'", service.Metadata.State)
	}
}

func (service *Service) Initialize() error {
	if service.Name == "" {
		service.Name = service.Metadata.Name
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
