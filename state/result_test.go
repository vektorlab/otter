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
	data, err := resultMap.ToJSON()
	if err != nil {
		fmt.Println("Failed to marshal ResultMap to JSON: ", err)
		t.Fail()
	}
	resultMap, err = ResultMapFromJson(data)
	if err != nil {
		fmt.Println("Failed to unmarshal ResultMap from JSON: ", err)
		t.Fail()
	}
}

func TestResultMapMergeNoDuplicate (t *testing.T) {
	state1 := stateSetup(simpleFileMeta, simpleFile, t)
	state2 := stateSetup(simplePackageMeta, simplePackage, t)
	resultMap1 := NewResultMap()
	resultMap1.Add(state1.Consistent())
	resultMap2 := NewResultMap()
	resultMap2.Add(state2.Consistent())
	resultMap1.Merge(resultMap2)
	fmt.Println(resultMap1.Results)
	if len(resultMap1.Results[helpers.GetHostName()]) != 2 {
		fmt.Println("Bad stateMap merge: ", len(resultMap1.Results[helpers.GetHostName()]))
		t.Fail()
	}
}

func TestResultMapMergeDuplicate (t *testing.T) {
	state1 := stateSetup(simpleFileMeta, simpleFile, t)
	state2 := stateSetup(simpleFileMeta, simpleFile, t)
	resultMap1 := NewResultMap()
	resultMap1.Add(state1.Consistent())
	resultMap2 := NewResultMap()
	resultMap2.Add(state2.Consistent())
	resultMap1.Merge(resultMap2)
	fmt.Println(len(resultMap1.Results[helpers.GetHostName()]))
	if len(resultMap1.Results[helpers.GetHostName()]) != 1 {
		fmt.Println("Bad stateMap merge: ", len(resultMap1.Results[helpers.GetHostName()]))
		t.Fail()
	}
}
