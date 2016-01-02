package state

import (
	"encoding/json"
)

type Result struct {
	Consistent bool      // The state is consistent with the operating system
	Metadata   *Metadata // The metadata of the state which returned this result
	Message    string    // A message returned by the state
}

func ResultsFromJson(data []byte) ([]Result, error) {
	results := make([]Result, 0)
	err := json.Unmarshal(data, &results)
	if err != nil {
		return results, err
	}
	return results, nil
}

func ResultsToJson(results []Result) ([]byte, error) {
	var data []byte
	data, err := json.Marshal(results)
	if err != nil {
		return data, err
	}
	return data, nil
}

type ResultMap struct {
	Results map[string][]Result
}

func (rm *ResultMap) Add(host string, result Result) {
	rm.Results[host] = append(rm.Results[host], result)
}

func (rm *ResultMap) ToJSON() ([]byte, error) {
	data, err := json.Marshal(rm.Results)
	if err != nil {
		return []byte(``), err
	}
	return data, nil
}
