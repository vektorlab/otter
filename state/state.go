package state

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
)

type State interface {
	Consistent() *Result    // Check to see if the state is consistent with the operating system's state
	Dump() ([]byte, error)  // Dump the state to a JSON byte array
	Execute() *Result       // Execute the state if it is not already Consistent
	Initialize() error      // Initialize the state validating loaded fields
	Meta() Metadata         // Return the state's metadata ("Name", "Type", and "state")
	Requirements() []string // Return the state's requirements // TODO: use in ordering of the state's execution and do not execute on failure
}

type StateMap struct {
	States []State
}

/*
Add a new state to the StateMap
*/
func (sm *StateMap) Add(entry State) error {
	if sm.Exists(entry.Meta(), false) {
		return fmt.Errorf("Detected duplicate state entry: %s", entry.Meta())
	}
	requirements := entry.Requirements()
	for _, requirement := range requirements {
		if !sm.Exists(Metadata{Name: requirement}, true) {
			return fmt.Errorf("Unable to find requirement: %s", requirement)
		}
	}
	sm.States = append(sm.States, entry)
	return nil
}

/*
Add multiple states to the StateMap
*/
func (sm *StateMap) AddMany(entries []State, attempts, max int) error {
	attempts++
	retry := make([]State, 0)
	for _, entry := range entries {
		err := sm.Add(entry)
		if err != nil {
			retry = append(retry, entry)
		}
	}
	if len(retry) >= 1 {
		if attempts > max {
			return fmt.Errorf("Unable to load %d states: %s", len(retry), retry)
		}
		return sm.AddMany(retry, attempts, max)
	}
	return nil
}

/*
Check to see if a given state already exists
*/
func (sm *StateMap) Exists(other Metadata, nameOnly bool) bool {
	for _, state := range sm.States {
		metadata := state.Meta()
		if nameOnly {
			if metadata.Name == other.Name {
				return true
			}
		}
		if metadata.Equal(&other) {
			return true
		}
	}
	return false
}

/*
Apply all states loaded in the StateMap
*/
func (sm *StateMap) Apply() *ResultMap {
	resultMap := NewResultMap()
	for _, state := range sm.States {
		resultMap.Add(state.Execute())
	}
	return resultMap
}

/*
Check if all states loaded in the StateMap are consistent
*/
func (sm *StateMap) Consistent() *ResultMap {
	resultMap := NewResultMap()
	for _, state := range sm.States {
		resultMap.Add(state.Consistent())
	}
	return resultMap
}

/*
Return all of the requirements for each loaded state
*/
func (sm *StateMap) requirements(entry []State) []string {
	reqs := make([]string, 0)
	for s := range entry {
		entryReqs := entry[s].Requirements()
		for r := range entryReqs {
			reqs = append(reqs, entryReqs[r])
		}
	}
	return reqs
}

/*
Validate that all requirements in each state exist and that there are no circular requirements
func (sm *StateMap) Validate() error {
	for name, entry := range sm.States {
		reqs := sm.requirements(entry)
		for req := range reqs {
			other, exists := sm.States[reqs[req]]
			if !exists {
				return fmt.Errorf("Unable to find requirement: %s", reqs[req])
			}
			otherReqs := sm.requirements(other)
			for req := range otherReqs {
				if name == otherReqs[req] {
					return fmt.Errorf("Detected circular requirement: %s", name)
				}
			}
		}
	}
	return nil
}
*/

/*
Dump this StateMap to a JSON byte array
*/
func (sm *StateMap) ToJson() ([]byte, error) {
	data, err := json.Marshal(sm.States)
	if err != nil {
		return []byte(``), err
	}
	return data, nil
}

func NewStateMap() *StateMap {
	sm := &StateMap{
		States: make([]State, 0),
	}
	return sm
}

/*
Load a StateMap from a JSON Byte array which has already been initialized (contains State Metadata)
*/
func StateMapFromProcessedJson(data []byte) (*StateMap, error) {
	sm := NewStateMap()
	raw := make(map[string][]json.RawMessage)
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return sm, err
	}
	states := make([]State, 0)
	for _, value := range raw {
		for _, entry := range value {
			metadata, err := MetadataFromJSON(entry)
			if err != nil {
				return sm, err
			}
			state, err := StateFactory(metadata, entry)
			if err != nil {
				return sm, err
			}
			err = state.Initialize()
			if err != nil {
				return sm, err
			}
			states = append(states, state)
		}
	}
	err = sm.AddMany(states, 0, len(states))
	return sm, err
}

/*
Load a StateMap from a JSON byte array
*/
func StateMapFromJson(data []byte) (*StateMap, error) {
	sm := NewStateMap()
	raw := make(map[string]map[string]json.RawMessage)
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return nil, err
	}
	states := make([]State, 0)
	for name, value := range raw {
		for keyword, data := range value {
			split := strings.Split(keyword, ".")
			metadata := Metadata{name, split[0], split[1]}
			state, err := StateFactory(metadata, data)
			if err != nil {
				return sm, err
			}
			err = state.Initialize()
			if err != nil {
				return sm, err
			}
			states = append(states, state)
		}
	}
	err = sm.AddMany(states, 0, len(states))
	return sm, err
}

/*
Load a StateMap from a YAML byte array
*/
func StateMapFromYaml(data []byte) (*StateMap, error) {
	j, err := yaml.YAMLToJSON(data)
	if err != nil {
		return nil, err
	}
	return StateMapFromJson(j)

}

/*
Load a YAML file from a given path, if the file doesn't exist default to ~/.otter
*/
func StateMapFromYamlPath(path string) (*StateMap, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		user, err := user.Current()
		if err != nil {
			return nil, err
		}
		path = user.HomeDir + "/.otter"
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return StateMapFromYaml(data)
}
