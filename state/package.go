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
	"github.com/mitchellh/mapstructure"
	"github.com/vektorlab/otter/helpers"
)

type Package struct {
	Name     string   `mapstructure:"name"`
	Version  string   `mapstructure:"version"`
	Metadata Metadata `mapstructure:"-"`
	Require  []string `mapstructure:"require"`
}

func (pkg *Package) Initialize() error {
	if pkg.Name == "" {
		pkg.Name = pkg.Metadata.Name
	}
	state := pkg.Metadata.State
	if ! (state == "installed" || state == "removed") {
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
	return nil
}

func PackageFromStructure(metadata Metadata, structure interface{}) (*Package, error) {
	var err error

	pkg := Package{
		Metadata: metadata,
	}

	err = mapstructure.Decode(structure, &pkg)
	if err != nil {
		return nil, err
	}

	err = pkg.Initialize()
	if err != nil {
		return nil, err
	}

	return &pkg, nil
}
