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
)

type Package struct {
	Name     string   `mapstructure:"name"`
	Version  string   `mapstructure:"version"`
	Metadata Metadata `mapstructure:"-"`
	Require  []string `mapstructure:"require"`
}

func (pkg *Package) validateState() error {
	switch pkg.Metadata.State {
	case "installed":
		return nil
	case "removed":
		return nil
	default:
		return fmt.Errorf("Invalid state choosen for package: %s, should be one of 'installed' or 'removed'", pkg.Metadata.State)
	}
}

func (pkg *Package) Initialize() error {
	if pkg.Name == "" {
		pkg.Name = pkg.Metadata.Name
	}
	return pkg.validateState()
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
