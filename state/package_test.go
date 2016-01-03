package state

import (
	"fmt"
	"testing"
)

func TestPackageConsistent(t *testing.T) {
	state := stateSetup(simplePackageMeta, simplePackage, t)
	result := state.Consistent()
	if result.Consistent != false {
		fmt.Println("Detected non-existant package: ", result.Metadata.Name)
	}
	fmt.Println(result.Consistent)
}

func TestPackageExecute(t *testing.T) {
	state := stateSetup(simplePackageMeta, simplePackage, t)
	result := state.Execute()
	if result.Consistent != false {
		fmt.Println("Installed non-existant package: ", result.Metadata.Name)
	}
}
