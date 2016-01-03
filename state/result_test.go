package state

import (
	"fmt"
	"testing"
	"github.com/vektorlab/otter/helpers"
)

func TestResultMap (t *testing.T) {
	state := stateSetup(simpleFileMeta, simpleFile, t)
	resultMap := NewResultMap()
	resultMap.Add(state.Execute())
	if resultMap.Results[helpers.GetHostName()][0].Consistent != false {
		fmt.Println("ResultMap failed to load result")
		t.Fail()
	}
}
