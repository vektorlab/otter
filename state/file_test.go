package state

import (
	"fmt"
	"testing"
)

func TestFileConsistent(t *testing.T) {
	state := stateSetup(simpleFileMeta, simpleFile, t)
	result := state.Consistent()
	if result.Consistent != false {
		fmt.Println("Detected non-existant file: ", result.Metadata.Name)
	}
	fmt.Println(result.Consistent)
}
