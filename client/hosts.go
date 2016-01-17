package clients

import (
	"fmt"
	etcd "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"strings"
	"time"
)

/*
List all registered remote hosts.
*/
func (otter *Otter) ListHosts() ([]string, error) {
	var hosts []string
	response, err := otter.etcdKeysApi.Get(context.Background(), "/ping", &etcd.GetOptions{Recursive: true})
	if err != nil {
		if strings.Contains(err.Error(), "Key not found") {
			return hosts, nil
		} else {
			return hosts, err
		}
	}
	for _, node := range response.Node.Nodes {
		hosts = append(hosts, strings.Split(node.Key, "/")[2])
	}
	return hosts, nil
}

/*
Update the list of remote hosts with the specified hostname.
*/
func (otter *Otter) UpdateHost(hostname, previous string) (string, error) {
	options := etcd.SetOptions{
		TTL: 60 * time.Second,
	}
	now := time.Now().String()
	if previous != "" {
		options.PrevValue = previous
	}
	_, err := otter.etcdKeysApi.Set(context.Background(), fmt.Sprintf("/ping/%s", hostname), now, &options)
	if err != nil {
		return now, err
	}
	return now, nil
}
