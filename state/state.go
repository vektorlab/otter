package state

import (
	"os"
	"fmt"
	"os/user"
	"strings"
	"io/ioutil"
	"encoding/json"
	"github.com/ghodss/yaml"
)

type State interface {
	Consistent() (bool, error) // Check to see if the state is consistent with the operating system's state
	Dump() ([]byte, error)     // Dump the state to a JSON byte array
	Execute() error            // Execute the state if it is not already Consistent
	Initialize() error         // Initialize the state validating loaded fields
	Meta() Metadata            // Return the state's metadata ("Name", "Type", and "state")
	Requirements() []string    // Return the state's requirements // TODO: use in ordering of the state's execution and do not execute on failure
}

type StateMap struct {
	States map[string][]State
}

func (sm *StateMap) Add(name string, entry State) {
	sm.States[name] = append(sm.States[name], entry)
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
*/
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

func (sm *StateMap) ToJson() ([]byte, error) {
	data, err := json.Marshal(sm.States)
	if err != nil {
		return []byte(``), err
	}
	return data, nil
}

func StateMapFromProcessedJson(data []byte) (StateMap, error) {
	sm := StateMap{
		States: make(map[string][]State),
	}
	raw := make(map[string][]json.RawMessage)
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return sm, err
	}
	for name, value := range raw {
		for _, entry := range value {
			metadata, err := getMetadata(entry)
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
			sm.Add(name, state)
		}
	}
	return sm, nil
}

func StateMapFromJson(data []byte) (StateMap, error) {
	sm := StateMap{
		States: make(map[string][]State),
	}
	raw := make(map[string]map[string]json.RawMessage)
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return sm, err
	}
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
			sm.Add(name, state)
		}
	}
	return sm, nil
}

func StateMapFromYaml(data []byte) (StateMap, error) {
	var sm StateMap
	j, err := yaml.YAMLToJSON(data)
	if err != nil {
		return sm, err
	}
	return StateMapFromJson(j)

}

/*
Load a YAML file from a given path, if the file doesn't exist default to ~/.otter
 */
func StateMapFromYamlPath(path string) (StateMap, error) {
	var sm StateMap
	if _, err := os.Stat(path); os.IsNotExist(err) {
		user, err := user.Current()
		if err != nil {
			return sm, err
		}
		path = user.HomeDir + "/.otter"
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return sm, err
	}
	return StateMapFromYaml(data)
}
