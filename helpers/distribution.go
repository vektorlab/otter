package helpers

import (
	"fmt"
	"gopkg.in/ini.v1"
	"runtime"
)

type Distro struct {
	Family     string // Major distribution type currently "debian" or "centos"
	InitSystem string // The Linux init system used by this operating system
	Version    string // Major version of the distribution // TODO
}

func (d *Distro) SetFamily(family string) error {

	switch family {
	case "debian":
		d.Family = family
	case "centos":
		d.Family = family
	default:
		return fmt.Errorf("Unknown Linux distribution: %s", family)
	}
	return nil
}

func (d *Distro) SetInitSystem() error {

	switch d.Family {
	case "debian":
		d.InitSystem = "systemd" // TODO: Base off distribution.Version
	case "centos":
		d.InitSystem = "sysv"
	default:
		return fmt.Errorf("Unknown Init system for distribution: %s", d.Family)
	}
	return nil
}

func GetDistro() (*Distro, error) {
	d := Distro{}
	if runtime.GOOS == "linux" {
		i, err := ini.Load([]byte(""), "/etc/os-release")
		section, err := i.Section("").GetKey("ID_LIKE")
		if err != nil {
			return nil, err
		}
		err = d.SetFamily(section.String())
		if err != nil {
			return nil, err
		}
		err = d.SetInitSystem()
		if err != nil {
			return nil, err
		}
	}
	return &d, nil
}
