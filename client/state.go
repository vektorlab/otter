package clients

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	etcd "github.com/coreos/etcd/client"
	"github.com/vektorlab/otter/state"
	"golang.org/x/net/context"
)

func (otter *Otter) RetrieveStateMap() (state.StateMap, error) {
	var stateMap state.StateMap
	key := fmt.Sprintf("/state/%s", otter.Hostname)
	response, err := otter.etcdKeysApi.Get(context.Background(), key, &etcd.GetOptions{})
	if err != nil {
		log.Printf("Unable to load state from key %s", key)
		return stateMap, err
	} else {
		raw := response.Node.Value
		stateMap, err := state.StateMapFromProcessedJson([]byte(raw))
		if err != nil {
			log.Fatalf("Bad JSON state payload: %s", raw)
		}
		return stateMap, nil
	}
}

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
