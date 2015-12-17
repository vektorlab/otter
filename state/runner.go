package state

import (
	"fmt"
)

/*
Run each state's consistency check and load save the results
*/
func Consistent(states map[string][]State) ([]Result, error) {
	results := make([]Result, 0)
	for _, groups := range states {
		for _, state := range groups {

			consistent, err := state.Consistent() // TODO: Differentiate between results and errors
			if err != nil {
				return results, err
			}
			meta := state.Meta()
			result := Result{
				Consistent: consistent,
				Metadata:   &meta,
				Message:    "",
			}
			results = append(results, result)
		}
	}
	return results, nil
}

/*
Execute each state
*/
func Execute(states map[string][]State) error {
	for _, groups := range states {
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
Return all of the requirements for each loaded state
*/
func requirements(entry []State) []string {
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
func Validate(states map[string][]State) error {
	for name, entry := range states {
		reqs := requirements(entry)
		for req := range reqs {
			other, exists := states[reqs[req]]
			if !exists {
				return fmt.Errorf("Unable to find requirement: %s", reqs[req])
			}
			otherReqs := requirements(other)
			for req := range otherReqs {
				if name == otherReqs[req] {
					return fmt.Errorf("Detected circular requirement: %s", name)
				}
			}
		}
	}
	return nil
}
