package state

import (
	"fmt"
)

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
