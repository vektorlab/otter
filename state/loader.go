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
	"github.com/ghodss/yaml"
	"io/ioutil"
)

// func StatesFromJson (data []byte) (States, error)

func StatesToJson(states map[string][]State) ([]byte, error) {
	var data []byte
	data, err := json.Marshal(states)
	if err != nil {
		return data, err
	}
	return data, nil
}

func StatesFromYaml(data []byte) (map[string][]State, error) {

	var err error

	raw := make(map[string]map[string]json.RawMessage)
	states := make(map[string][]State)

	j, err := yaml.YAMLToJSON(data)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(j, &raw)

	if err != nil {
		return nil, err
	}

	for name, value := range raw {
		for keyword, data := range value {
			state, err := StateFactory(name, keyword, data)
			if err != nil {
				return nil, err
			}
			err = state.Initialize()
			if err != nil {
				return nil, err
			}
			states[name] = append(states[name], state)
		}
	}

	//err = loader.validate()

	if err != nil {
		return nil, err
	}

	return states, nil
}

func StatesFromYamlPath(path string) (map[string][]State, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return StatesFromYaml(data)
}
