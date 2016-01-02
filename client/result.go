package clients

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	etcd "github.com/coreos/etcd/client"
	"github.com/vektorlab/otter/state"
	"golang.org/x/net/context"
	"time"
)

func (otter *Otter) SaveResults(id string, results []state.Result) error {
	data, err := state.ResultsToJson(results)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("/result/%s", id)
	_, err = otter.etcdKeysApi.Set(context.Background(), key, string(data), &etcd.SetOptions{})
	log.Printf("Saved results to %s", key)
	if err != nil {
		return err
	}
	return nil
}

/*
Wait for a specific event to occur in the /event/<hostname> keyspace.
*/
func (otter *Otter) WaitForResult(id string) (string, string, error) {

	key := fmt.Sprintf("/result/%s", id)

	result, value, err := otter.WaitForChange(key, false, 10*time.Second)

	if err != nil {
		return "", "", err
	}

	return result, value, nil

}
