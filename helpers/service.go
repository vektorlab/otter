package helpers

import (
	"fmt"
	"github.com/coreos/go-systemd/dbus"
)

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
	distro, err := GetDistro()
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
