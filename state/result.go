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

type ResultMap struct {
	Results map[string][]*Result
	Host string
}

/*
Add a new Result to the ResultMap
 */
func (rm *ResultMap) Add(result *Result) {
	if result.Host == "" {
		result.Host = rm.Host // Assume this result was generated on this host if it is not specified.
	}
	if rm.Exists(result) {
		return // Result already is added to this map, silently ignore it
	} else {
		if _, hostExists := rm.Results[result.Host]; !hostExists {
			rm.Results[result.Host] = make([]*Result, 0) // Host entry does not exist, create an empty array
			rm.Add(result) // Add it again with the new map entry
			return
		}
		rm.Results[result.Host] = append(rm.Results[result.Host], result) // Add the Result to an existing host array
	}
}

/*
Check to see if there is a result for the specified host and metadata
 */
func (rm *ResultMap) Exists(other *Result) bool {
	if results, exists := rm.Results[other.Host]; exists {
		for _, result := range results {
			if result.Metadata.Equal(other.Metadata) {
				return true
			}
		}
	}
	return false
}

func (rm *ResultMap) Merge(other *ResultMap) {
	for _, results := range other.Results {
		for _, result := range results {
			rm.Add(result)
		}
	}
}
/*
Dump the ResultMap to JSON
 */
func (rm *ResultMap) ToJSON() ([]byte, error) {
	data, err := json.Marshal(rm.Results)
	if err != nil {
		return []byte(``), err
	}
	return data, nil
}

/*
Get a new ResultMap object
 */
func NewResultMap() *ResultMap {
	resultMap := &ResultMap{
		Results: make(map[string][]*Result),
		Host: helpers.GetHostName(),
	}
	return resultMap
}

/*
Create a ResultMap from JSON byte array
 */
func ResultMapFromJson(data []byte) (*ResultMap, error) {
	resultMap := NewResultMap()
	err := json.Unmarshal(data, &resultMap.Results)
	if err != nil {
		return nil, err
	}
	return resultMap, nil
}
