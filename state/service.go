/*
A service represents a running daemonized process on an operating system

States -
  running: The service is running // TODO
  stopped: The service is stopped // TODO
*/
package state

import (
	"fmt"
	"github.com/coreos/go-systemd/dbus"
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
		Metadata:   &service.Metadata,
		Consistent: false,
	}
	running, err := ServiceRunning(service.Name)
	if err != nil {
		result.Message = err.Error()
		return result
	}
	result.Consistent = service.Running == running
	return result
}

func (service *Service) Apply() *Result {
	result := &Result{
		Metadata:   &service.Metadata,
		Consistent: service.Consistent().Consistent,
	}
	if result.Consistent == false {
		err := ChangeServiceState(service.Name, service.Running)
		if err != nil {
			result.Message = err.Error()
			return result
		}
		result.Consistent = service.Consistent().Consistent
	}
	return result
}

func (service *Service) Requirements() []string {
	return service.Require
}

func (service *Service) Meta() Metadata {
	return service.Metadata
}

/*
Check if the specified unit name is running with Systemd
*/
func SystemdUnitRunning(name string) (bool, error) {
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
Check if the specified service name is running on the operating system.
*/
func ServiceRunning(name string) (bool, error) {
	distro, err := helpers.GetDistro()
	if err != nil {
		return false, err
	}
	switch distro.InitSystem {
	case "systemd":
		return SystemdUnitRunning(name)
	default:
		return false, fmt.Errorf("Unsupported init system %s", distro.InitSystem)
	}
}

/*
Start a Systemd unit and wait for it to return. This method may block.
*/
func startSystemdUnit(name string) error {
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
func stopSystemdUnit(name string) error {
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

/*
Update the service state, may be "running" or "stopped"
*/
func ChangeServiceState(name string, running bool) error {
	distro, err := helpers.GetDistro()
	if err != nil {
		return err
	}
	switch distro.InitSystem {
	case "systemd":
		if running {
			return startSystemdUnit(name)
		} else {
			return stopSystemdUnit(name)
		}
	default:
		return fmt.Errorf("Unsupported init system %s", distro.InitSystem)
	}
}
