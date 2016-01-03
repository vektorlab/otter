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
	"github.com/vektorlab/otter/helpers"
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

func (pkg *Package) Consistent() (bool, error) {
	var err error
	status, err := helpers.GetPackageStatus(pkg.Name)
	if err != nil {
		return false, err
	}
	return status == pkg.Metadata.State, nil
}

func (pkg *Package) Execute() error {
	consistent, err := pkg.Consistent()
	if err != nil {
		return err
	}
	if !consistent {
		err := helpers.InstallPackage(pkg.Name)
		if err != nil {
			return err
		}
	}
	return nil
}
