package executors

import (
	"fmt"
	"os/exec"
	"strings"
	"encoding/json"
	"github.com/vektorlab/otter/helpers"
	"github.com/vektorlab/otter/state"
)

type Package struct {
	state  state.Package
	status string
}

func (pkg *Package) Consistent() (bool, error) {
	var err error

	pkg.status, err = GetPackageStatus(pkg.state.Name)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (pkg *Package) Execute() (Result, error) {
	return Result{}, nil
}

func (pkg *Package) Load(data []byte) error {
	pkg.state = state.Package{}
	err := json.Unmarshal(data, &pkg.state)
	return err
}

func GetDpkgPackage(name string) (string, error) {
	out, err := exec.Command("dpkg", "-l", name).Output() // TODO: Security
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		split := strings.Fields(line)
		if len(split) >= 1 {
			f := string(split[0][0])
			if f == "h" || f == "i" || f == "p" || f == "r" || f == "u" {
				return split[0], nil
			}
		}
	}
	return "", nil
}

func GetPackageStatus(name string) (string, error) {
	distro, err := helpers.GetDistro()
	if err != nil {
		return "", err
	}

	switch distro.Family {
	case "debian":
		status, err := GetDpkgPackage(name)
		if err != nil {
			return "", err
		}
		switch status {
		case "ii":
			return "installed", nil
		case "rc":
			return "removed", nil
		case "":
			return "purged", nil
		default:
			return "", fmt.Errorf("Package %s is in an unknown state: %s", name, status)
		}
	default:
		return "", fmt.Errorf("Unsupported operating system: %s", distro.Family)
	}
}