/*
A package represents the state of a package on an operating system.
States -
  installed: The package is installed on the operating system. // TODO
  removed: The package is removed from the operating system. // TODO
*/

package state

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/vektorlab/otter/helpers"
	"os/exec"
	"strings"
)

type Package struct {
	Name     string   `json:"name"`
	Version  string   `json:"version"`
	Metadata Metadata `json:"metadata"`
	Require  []string `json:"require"`
}

func (pkg *Package) Initialize() error {
	if pkg.Name == "" {
		pkg.Name = pkg.Metadata.Name
	}
	state := pkg.Metadata.State
	if !(state == "installed" || state == "removed") {
		return fmt.Errorf("Invalid package state: %s", state)
	}
	return nil
}

func (pkg *Package) Dump() ([]byte, error) {
	return json.Marshal(pkg)
}

func (pkg *Package) Requirements() []string {
	return pkg.Require
}

func (pkg *Package) Meta() Metadata {
	return pkg.Metadata
}

func (pkg *Package) Consistent() *Result {
	result := &Result{
		Metadata:   &pkg.Metadata,
		Consistent: false,
	}
	status, err := GetPackageStatus(pkg.Name)
	if err != nil {
		result.Message = err.Error()
		return result
	}
	result.Consistent = status == pkg.Metadata.State
	return result
}

func (pkg *Package) Execute() *Result {
	result := &Result{
		Metadata:   &pkg.Metadata,
		Consistent: pkg.Consistent().Consistent,
	}
	if result.Consistent == false {
		err := InstallPackage(pkg.Name)
		if err != nil {
			result.Message = err.Error()
			return result
		}
		result.Consistent = pkg.Consistent().Consistent
	}
	return result
}

/*
Get the status of a package in DPKG
*/
func GetDpkgPackage(name string) (string, error) {
	out, err := exec.Command("dpkg", "-l", name).Output() // TODO: Security
	if err != nil {                                       // dpkg doesn't know about the package
		return "rc", nil
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

/*
Check to see if a package is installed on the operating system
*/
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

/*
Install a package with apt-get
*/
func installAptPackage(name string) error {
	out, err := exec.Command("apt-get", "update").CombinedOutput()
	if err != nil {
		log.Warningln(string(out))
		return err
	}
	log.Warningln(string(out))
	out, err = exec.Command("apt-get", "install", "-y", name).CombinedOutput()
	if err != nil {
		return err
	}
	log.Warningln(string(out))
	return nil
}

/*
Install a package on the operating system
*/
func InstallPackage(name string) error {
	distro, err := helpers.GetDistro()
	if err != nil {
		return err
	}
	switch distro.Family {
	case "debian":
		err := installAptPackage(name)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unsupported operating system: %s", distro.Family)
	}
	return nil
}
