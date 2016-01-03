package clients

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	etcd "github.com/coreos/etcd/client"
	"github.com/vektorlab/otter/state"
	"golang.org/x/net/context"
	"time"
)

func (otter *Otter) SaveResultMap(id string, resultMap *state.ResultMap) error {
	data, err := resultMap.ToJSON()
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
Wait for a new ResultMap to be updated in the /result/<id> keyspace.
*/
func (otter *Otter) WaitForResults(id string) (*state.ResultMap, error) {
	key := fmt.Sprintf("/result/%s", id)
	_, value, err := otter.WaitForChange(key, false, 10*time.Second)
	if err != nil {
		return nil, err
	}
	return state.ResultMapFromJson([]byte(value))
}
