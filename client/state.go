package clients

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	etcd "github.com/coreos/etcd/client"
	"github.com/vektorlab/otter/state"
	"golang.org/x/net/context"
)

/*
Retrieve state data for this server from etcd, check if it is consistent, then return the result.
*/

func (otter *Otter) RetrieveStateMap() (map[string][]state.State, error) {
	key := fmt.Sprintf("/state/%s", otter.Hostname)
	response, err := otter.etcdKeysApi.Get(context.Background(), key, &etcd.GetOptions{})
	if err != nil {
		log.Printf("Unable to load state from key %s", key)
		return nil, err
	} else {
		raw := response.Node.Value
		loaded, err := state.StatesFromProcessedJson([]byte(raw))
		if err != nil {
			log.Fatalf("Bad JSON state payload: %s", raw)
		}
		return loaded, nil
	}
}

/*
Submit a local state to all registered hosts.
*/
func (otter *Otter) SubmitState(state string) error {

	hosts, err := otter.ListHosts()

	if err != nil {
		return err
	}

	for _, host := range hosts {
		_, err := otter.etcdKeysApi.Set(context.Background(), fmt.Sprintf("/state/%s", host), state, &etcd.SetOptions{})
		if err != nil {
			return err
		}
		log.Printf("Updated state for host: %s", host)
	}

	return nil
}

/*
Run each state's consistency check and load save the results
*/
func (otter *Otter) CheckConsistent(s state.State) (state.Result, error) {
	metadata := s.Meta()
	result := state.Result{
		Metadata: &metadata,
	}
	consistent, err := s.Consistent() // TODO: Differentiate between results and errors
	if err != nil {
		result.Consistent = false
		result.Message = err.Error()
	} else {
		result.Consistent = consistent
		result.Message = ""
	}
	return result, nil
}

/*
Execute each state
*/
func (otter *Otter) ExecuteState(states map[string][]state.State) ([]state.Result, error) {
	results := make([]state.Result, 0)
	for _, groups := range states {
		for _, s := range groups {
			metadata := s.Meta()
			log.Printf("Applying state: %s.%s.%s", metadata.Name, metadata.State, metadata.Type)
			err := s.Execute()
			result := state.Result{Metadata: &metadata}
			if err != nil {
				result.Consistent = false
				result.Message = err.Error()
				log.Warnf("State failed: %s - %s", metadata.State, result.Message)
			} else {
				result.Consistent = true
			}
			results = append(results, result)
		}
	}
	return results, nil
}
