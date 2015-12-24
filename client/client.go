/*
Key spaces:
	/ping/<hostname> - All remote servers will update the ping keyspace every 15 seconds.
	/state/<hostname> - Latest requested state for specified hostname.
	/command/<hostname>/<type> - Requested action to be performed on the remote host.
	/result/<id> - The key to save the result of a command in.
*/

package clients

import (
	log "github.com/Sirupsen/logrus"
	etcd "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"os"
	"time"
)

type Otter struct {
	etcdKeysApi etcd.KeysAPI
	Hostname    string
}

func NewOtterClient(servers []string) (*Otter, error) {

	client := Otter{}

	cfg := etcd.Config{
		Endpoints:               servers,
		Transport:               etcd.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	etcdClient, err := etcd.New(cfg)

	if err != nil {
		return nil, err
	}

	api := etcd.NewKeysAPI(etcdClient)

	client.etcdKeysApi = api

	hostname, err := os.Hostname()

	if err != nil {
		return nil, err
	}

	client.Hostname = hostname

	return &client, nil
}

/*
Wait for a key to be added or modified in etcd
*/
func (otter *Otter) WaitForChange(key string, recurse bool) (string, string, error) {

	log.Printf("Waiting for change to keyspace: %s", key)

	options := etcd.WatcherOptions{Recursive: recurse}
	watcher := otter.etcdKeysApi.Watcher(key, &options)

	response, err := watcher.Next(context.Background())

	if err != nil {
		return "", "", err
	}

	log.Printf("Got change from key: %s", response.Node.Key)

	return response.Node.Key, response.Node.Value, nil
}
