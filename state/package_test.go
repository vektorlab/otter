package state

import (
	"fmt"
	"testing"
)

func TestPackageConsistent(t *testing.T) {
	state := stateSetup(simplePackageMeta, simplePackage, t)
	result := state.State()
	if result.Consistent != false {
		fmt.Println("Detected non-existant package: ", result.Metadata.Name)
	}
	fmt.Println(result.Consistent)
}

func TestPackageExecute(t *testing.T) {
	state := stateSetup(simplePackageMeta, simplePackage, t)
	result := state.Apply()
	if result.Consistent != false {
		fmt.Println("Installed non-existant package: ", result.Metadata.Name)
	}
}
