package state

import (
	"encoding/json"
)

type Result struct {
	Consistent bool      // The state is consistent with the operating system
	Metadata   *Metadata // The metadata of the state which returned this result
	Message    string    // A message returned by the state
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
