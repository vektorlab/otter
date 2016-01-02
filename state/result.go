package state

import (
	"encoding/json"
)

type Result struct {
	Consistent bool      // The state is consistent with the operating system
	Metadata   *Metadata // The metadata of the state which returned this result
	Message    string    // A message returned by the state
}

type Metadata struct {
	Name  string // Unique name to associate with a state
	Type  string // The type of state "package", "file", etc.
	State string // The desired state "installed", "rendered", etc.
}


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

func ResultsToJson(results []Result) ([]byte, error) {
	data, err := json.Marshal(results)
	if err != nil {
		return []byte(``), err
	}
	return data, nil
}

func ResultsFromJson(data []byte) ([]Result, error) {
	results := make([]Result, 0)
	err := json.Unmarshal(data, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}
