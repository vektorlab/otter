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
}

func (pkg *Package) Meta() Metadata {
	return pkg.Metadata
}

func (pkg *Package) State() *Result {
	result := &Result{
		Metadata:   &pkg.Metadata,
		Consistent: false,
	}
	status, version, err := pkg.GetPackageStatus(pkg.Name)
	if err != nil {
		result.Message = err.Error()
		return result
	}
	switch result.Metadata.State {
	case "installed":
		if status == "installed" {
			if pkg.Version != "" {
				result.Consistent = version == pkg.Version
			}
		}
	case "removed":
		result.Consistent = status == ""
	}
	return result
}

func (pkg *Package) Apply() *Result {
	result := pkg.State()
	if result.Consistent == true {
		return result
	}
	switch pkg.Metadata.State {
	case "installed":
		err := pkg.InstallPackage(pkg.Name, pkg.Version)
		if err != nil {
			result.Message = err.Error()
			return result
		}
		result.Message = "Package Installed"
		result.Consistent = true
	case "removed":
		err := pkg.RemovePackage(pkg.Name)
		if err != nil {
			result.Message = err.Error()
			return result
		}
		result.Message = "Package Removed"
		result.Consistent = true
	}
	return result
}

/*
Create and validate a new Package State
*/
func newPackage(metadata Metadata, data []byte) (*Package, error) {
	pkg := &Package{}
	err := json.Unmarshal(data, &pkg)
	if err != nil {
		return nil, err
	}
	pkg.Metadata = metadata
	switch metadata.State {
	case "installed":
	case "removed":
	default:
		return nil, fmt.Errorf("Invalid package state: %s", metadata.State)
	}
	if pkg.Name == "" {
		pkg.Name = metadata.Name
	}
	return pkg, nil
}

/*
Check to see if a package is installed on the operating system
*/
func (pkg *Package) GetPackageStatus(name string) (string, string, error) {
	distro, err := helpers.GetDistro()
	if err != nil {
		return "", "", err
	}

	switch distro.Family {
	case "debian":
		status, version, err := pkg.GetDpkgPackage(name)
		if err != nil {
			return "", "", err
		}
		switch status {
		case "ii":
			return "installed", version, nil
		case "rc":
			return "removed", version, nil
		case "":
			return "purged", version, nil
		default:
			return "", "", fmt.Errorf("Package %s is in an unknown state: %s", name, status)
		}
	default:
		return "", "", fmt.Errorf("Unsupported operating system: %s", distro.Family)
	}
}

/*
Install a package on the operating system
*/
func (pkg *Package) InstallPackage(name, version string) error {
	distro, err := helpers.GetDistro()
	if err != nil {
		return err
	}
	switch distro.Family {
	case "debian":
		if version != "" {
			name = name + fmt.Sprintf("==%s", version)
		}
		err := pkg.installAptPackage(name)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unsupported operating system: %s", distro.Family)
	}
	return nil
}

/*
Remove a package from the operating system
*/
func (pkg *Package) RemovePackage(name string) error {
	distro, err := helpers.GetDistro()
	if err != nil {
		return err
	}
	switch distro.Family {
	case "debian":
		err := pkg.removeAptPackage(name)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unsupported operating system: %s", distro.Family)
	}
	return nil
}

/*
Get the status of a package in DPKG
*/
func (pkg *Package) GetDpkgPackage(name string) (string, string, error) {
	out, err := exec.Command("dpkg", "-l", name).Output() // TODO: Security
	if err != nil {                                       // dpkg doesn't know about the package
		return "rc", "", nil
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		split := strings.Fields(line)
		if len(split) >= 1 {
			f := string(split[0][0])
			if f == "h" || f == "i" || f == "p" || f == "r" || f == "u" {
				return split[0], split[2], nil
			}
		}
	}
	return "", "", nil
}

/*
Install a package with apt-get
*/
func (pkg *Package) installAptPackage(name string) error {
	out, err := exec.Command("apt-get", "update").CombinedOutput()
	if err != nil {
		log.Warningln(string(out))
		return err
	}
	out, err = exec.Command("apt-get", "install", "-y", name).CombinedOutput()
	if err != nil {
		log.Warningln(string(out))
		return err
	}
	log.Println("Installed Apt package: ", string(out))
	return nil
}

/*
Remove a package with apt-get
*/
func (pkg *Package) removeAptPackage(name string) error {
	out, err := exec.Command("apt-get", "remove", name).CombinedOutput()
	if err != nil {
		return err
	}
	log.Println("Removed Apt package: ", out)
	return nil
}
