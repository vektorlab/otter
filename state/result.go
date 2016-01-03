package state

import (
	"encoding/json"
	"github.com/vektorlab/otter/helpers"
)

type Result struct {
	Host       string    // The host that produced this Result object
	Consistent bool      // The state is consistent with the operating system
	Metadata   *Metadata // The metadata of the state which returned this result
	Message    string    // A message returned by the state
}

func ResultsFromJson(data []byte) ([]*Result, error) {
	var results []*Result
	raw := make([]json.RawMessage, 0)
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return results, err
	}
	for _, rawResult := range raw {
		result := &Result{} // Result.Host must be modified by the Otter client
		err := json.Unmarshal(rawResult, result)
		if err != nil {
			return results, err
		}
		results = append(results, result)
	}
	return results, nil
}

/*
Return an array with a single result in state "Faulted", add the error message to the Result.
*/
func ResultsFromError(host string, err error) []*Result {
	results := make([]*Result, 1)
	results[0] = &Result{
		Host:    host,
		Message: err.Error(),
		Metadata: &Metadata{
			Name:  "Faulted",
			Type:  "Faulted",
			State: "Faulted",
		},
	}
	return results
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
	Results map[string][]*Result
	Host string
}

func (rm *ResultMap) Add(result *Result) {
	result.Host = rm.Host
	if _, exists := rm.Results[result.Host]; !exists {
		rm.Results[result.Host] = make([]*Result, 0)
	}
	rm.Results[result.Host] = append(rm.Results[result.Host], result)
}

func (rm *ResultMap) ToJSON() ([]byte, error) {
	data, err := json.Marshal(rm.Results)
	if err != nil {
		return []byte(``), err
	}
	return data, nil
}

func NewResultMap() *ResultMap {
	resultMap := &ResultMap{
		Results: make(map[string][]*Result),
		Host: helpers.GetHostName(),
	}
	return resultMap
}
