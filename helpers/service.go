package helpers

import (
	"fmt"
	"github.com/coreos/go-systemd/dbus"
)

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
	distro, err := GetDistro()
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
