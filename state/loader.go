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
	"strings"
)

func getMetadata(data json.RawMessage) (Metadata, error) {
	metadata := Metadata{}
	raw := make(map[string]json.RawMessage)
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return metadata, err
	}
	for key, value := range raw {
		if key == "metadata" {
			err := json.Unmarshal(value, &metadata)
			if err != nil {
				return metadata, err
			}
		}
	}
	return metadata, nil
}

func StatesToJson(states map[string][]State) ([]byte, error) {
	var data []byte
	data, err := json.Marshal(states)
	if err != nil {
		return data, err
	}
	return data, nil
}

func StatesFromProcessedJson(data []byte) (map[string][]State, error) {
	raw := make(map[string][]json.RawMessage)
	states := make(map[string][]State)
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return nil, err
	}
	for name, value := range raw {
		for _, entry := range value {
			metadata, err := getMetadata(entry)
			if err != nil {
				return nil, err
			}
			state, err := StateFactory(metadata, entry)
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
	return states, nil
}

func StatesFromJson(data []byte) (map[string][]State, error) {

	raw := make(map[string]map[string]json.RawMessage)
	states := make(map[string][]State)

	err := json.Unmarshal(data, &raw)

	if err != nil {
		return nil, err
	}

	for name, value := range raw {
		for keyword, data := range value {
			split := strings.Split(keyword, ".")
			metadata := Metadata{name, split[0], split[1]}
			state, err := StateFactory(metadata, data)
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

	if err != nil {
		return nil, err
	}

	return states, nil
}

func StatesFromYaml(data []byte) (map[string][]State, error) {
	j, err := yaml.YAMLToJSON(data)

	if err != nil {
		return nil, err
	}
	return StatesFromJson(j)

}

func StatesFromYamlPath(path string) (map[string][]State, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return StatesFromYaml(data)
}
