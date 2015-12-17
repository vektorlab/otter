/*

The state package provides translation for otter's declarative YAML configuration file
into State objects which are used to configure a host server.

An example YAML configuration might look like this:

/root/readme.md: <-- Name
  file.rendered: <-- {type}.{state}
    mode: 644 <-- Key/Value arguments
    source: git:///git@github.com/vektorlab/otter/README.md
*/

package state

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
)

type Loader struct {
	stateRaw map[string]map[string]json.RawMessage
	State    map[string][]State
	Results  []Result
}

/*
Dump the entire state out to JSON
*/
func (loader *Loader) Dump() ([]byte, error) {
	return json.Marshal(loader.State)
}

/*
Return the total number of state objects loaded
*/
func (loader *Loader) Count() int {
	var count int
	for _, value := range loader.State {
		count += len(value)
	}
	return count
}

/*
Return all of the requirements for each loaded state
*/
func (loader *Loader) requirements(entry []State) []string {
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
func (loader *Loader) validate() error {
	for name, entry := range loader.State {
		reqs := loader.requirements(entry)
		for req := range reqs {
			other, exists := loader.State[reqs[req]]
			if !exists {
				return fmt.Errorf("Unable to find requirement: %s", reqs[req])
			}
			otherReqs := loader.requirements(other)
			for req := range otherReqs {
				if name == otherReqs[req] {
					return fmt.Errorf("Detected circular requirement: %s", name)
				}
			}
		}
	}
	return nil
}

/*
Run each state's consistency check and load save the results
*/
func (loader *Loader) Consistent() error {
	for _, groups := range loader.State {
		for _, state := range groups {

			consistent, err := state.Consistent() // TODO: Differentiate between results and errors
			if err != nil {
				return err
			}
			meta := state.Meta()
			result := Result{
				Consistent: consistent,
				Metadata:   &meta,
				Message:    "",
			}
			loader.Results = append(loader.Results, result)
		}
	}
	return nil
}

/*
Execute each state
*/
func (loader *Loader) Execute() error {
	for _, groups := range loader.State {
		for _, state := range groups {
			err := state.Execute()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

/*
Load a state objects from a byte array
*/

func FromYaml(data []byte) (*Loader, error) {

	var err error

	loader := Loader{
		State:   make(map[string][]State),
		stateRaw: make(map[string]map[string]json.RawMessage),
		Results: make([]Result, 0),
	}

	j, err := yaml.YAMLToJSON(data)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(j, &loader.stateRaw)

	if err != nil {
		return nil, err
	}

	for name, value := range loader.stateRaw {
		for keyword, data := range value {
			state, err := StateFactory(name, keyword, data)
			if err != nil {
				return nil, err
			}
			err = state.Initialize()
			if err != nil {
				return nil, err
			}
			loader.State[state.Meta().Name] = append(loader.State[name], state)
		}
	}

	err = loader.validate()

	if err != nil {
		return nil, err
	}

	return &loader, nil
}

/*
Load state objects from a file path
*/

func FromPathToYaml(path string) (*Loader, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return FromYaml(data)
}
