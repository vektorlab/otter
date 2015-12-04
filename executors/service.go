package executors

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/go-systemd/dbus"
	"github.com/vektorlab/otter/helpers"
	"github.com/vektorlab/otter/state"
)

type Service struct {
	state state.Service
}

func (service *Service) Consistent() (bool, error) {
	running, err := ServiceRunning(service.state.Name)
	if err != nil {
		return false, err
	}
	return service.state.Running == running, nil
}

func (service *Service) Execute() (Result, error) {
	return Result{}, nil
}

func (service *Service) Load(data []byte) error {
	service.state = state.Service{}
	err := json.Unmarshal(data, &service.state)
	return err
}

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
