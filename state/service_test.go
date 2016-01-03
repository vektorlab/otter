package state

import (
	"fmt"
	"testing"
)

func TestServiceConsistent(t *testing.T) {
	state := stateSetup(simpleServiceMeta, simpleService, t)
	result := state.Consistent()
	if result.Consistent != false {
		fmt.Println("Detected running non-existant service: ", result.Metadata.Name)
	}
	fmt.Println(result.Consistent)
}

/*  TODO: systemd query times out when service is not found
func TestServiceExecute(t *testing.T) {
	state := stateSetup(simpleServiceMeta, simpleService, t)
	result := state.Execute()
	if result.Consistent != false {
		fmt.Println("Started non-existant service: ", result.Metadata.Name)
	}
}
*/
