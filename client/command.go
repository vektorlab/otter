package clients

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	etcd "github.com/coreos/etcd/client"
	"github.com/vektorlab/otter/helpers"
	"github.com/vektorlab/otter/state"
	"golang.org/x/net/context"
	"strings"
	"time"
)

/*
Run a command against remote hosts
*/

func (otter *Otter) SubmitCommand(host, command string) ([]state.Result, error) {

	key := fmt.Sprintf("/command/%s/%s", host, command)
	id := helpers.RandomString(8)

	_, err := otter.etcdKeysApi.Set(context.Background(), key, id, &etcd.SetOptions{})

	log.Printf("Submitted command with keyspace: %s", key)

	if err != nil {
		return nil, err
	}

	key, value, err := otter.WaitForResult(id)

	if err != nil {
		return nil, err
	}

	results, err := state.ResultsFromJson([]byte(value))

	return results, nil

}

/*
Wait for a command and then return it.
*/
func (otter *Otter) WaitForCommand(hostname string) (string, string, error) {

	key := fmt.Sprintf("/command/%s")

	key, id, err := otter.WaitForChange(fmt.Sprintf("/command/%s", hostname), true, 0 * time.Second)

	if err != nil {
		return "", "", err
	}

	split := strings.Split(key, "/")

	return split[len(split)-1], id, nil
}
