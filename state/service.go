/*
A service represents a running daemonized process on an operating system
States -
  running: The service is running
  stopped: The service is stopped
*/
package state

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/go-systemd/dbus"
	"github.com/vektorlab/otter/helpers"
)

type Service struct {
	Name     string   `json:"name"`
	Running  bool     `json:"running"`
	Metadata Metadata `json:"metadata"`
}

func (service *Service) Meta() Metadata {
	return service.Metadata
}

func (service *Service) State() *Result {
	result := &Result{
		Metadata:   &service.Metadata,
		Consistent: false,
	}
	switch service.Metadata.State {
	case "running":
		running, err := service.serviceRunning(service.Name)
		if err != nil {
			result.Message = err.Error()
			return result
		}
		result.Consistent = running == true
		return result
	case "stopped":
		running, err := service.serviceRunning(service.Name)
		if err != nil {
			result.Message = err.Error()
			return result
		}
		result.Consistent = running == false
	}
	return result
}

func (service *Service) Apply() *Result {
	result := service.State()
	if result.Consistent == true {
		return result
	}
	switch service.Metadata.State {
	case "running":
		err := service.changeServiceState(service.Name, service.Running)
		if err != nil {
			result.Message = err.Error()
			return result
		}
		result.Message = "Service is running"
		result.Consistent = true
	case "stopped":
		err := service.changeServiceState(service.Name, service.Running)
		if err != nil {
			result.Message = err.Error()
			return result
		}
		result.Message = "Service is stopped"
		result.Consistent = true
	}
	return result
}

/*
Create and validate a new Service State
*/
func newService(metadata Metadata, data []byte) (*Service, error) {
	service := &Service{}
	err := json.Unmarshal(data, &service)
	if err != nil {
		return nil, err
	}
	service.Metadata = metadata
	switch metadata.State {
	case "running":
		service.Running = true
	case "stopped":
		service.Running = false
	default:
		return nil, fmt.Errorf("Invalid service state: %s", metadata.State)
	}
	if service.Name == "" {
		service.Name = metadata.Name
	}
	return service, nil
}

/*
Check if the specified service name is running on the operating system.
*/
func (service *Service) serviceRunning(name string) (bool, error) {
	distro, err := helpers.GetDistro()
	if err != nil {
		return false, err
	}
	switch distro.InitSystem {
	case "systemd":
		return service.systemdUnitRunning(name)
	default:
		return false, fmt.Errorf("Unsupported init system %s", distro.InitSystem)
	}
}

/*
Update the service state, may be "running" or "stopped"
*/
func (service *Service) changeServiceState(name string, running bool) error {
	distro, err := helpers.GetDistro()
	if err != nil {
		return err
	}
	switch distro.InitSystem {
	case "systemd":
		if running {
			return service.startSystemdUnit(name)
		} else {
			return service.stopSystemdUnit(name)
		}
	default:
		return fmt.Errorf("Unsupported init system %s", distro.InitSystem)
	}
}

/*
Check if the specified unit name is running with Systemd
*/
func (service *Service) systemdUnitRunning(name string) (bool, error) {
	conn, err := dbus.New()
	defer conn.Close()
	if err != nil {
		return false, err
	}
	units, err := conn.ListUnits()
	if err != nil {
		return false, err
	}
	for _, unit := range units {
		if unit.Name == name+".service" {
			return unit.ActiveState == "active", nil
		}
	}
	return true, nil
}

/*
Start a Systemd unit and wait for it to return. This method may block.
*/
func (service *Service) startSystemdUnit(name string) error {
	conn, err := dbus.New()
	defer conn.Close()
	if err != nil {
		return err
	}
	c := make(chan string)
	conn.StartUnit(name, "replace", c)
	done := <-c
	if done != "finished" {
		return fmt.Errorf("Problem starting systemd unit, dbus responded: %s", done)
	}
	return nil
}

/*
Stop a Systemd unit and wait for it to return, this method may block.
*/
func (service *Service) stopSystemdUnit(name string) error {
	conn, err := dbus.New()
	defer conn.Close()
	if err != nil {
		return err
	}
	c := make(chan string)
	conn.StopUnit(name, "replace", c)
	done := <-c
	if done != "finished" {
		return fmt.Errorf("Problem starting systemd unit, dbus responded: %s", done)
	}
	return nil
}
