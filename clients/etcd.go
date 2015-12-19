package clients

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"strings"
	"time"
)

func NewKeysApi(servers []string) (client.KeysAPI, error) {
	cfg := client.Config{
		Endpoints: servers,
		Transport: client.DefaultTransport,
	}

	c, err := client.New(cfg)

	if err != nil {
		return nil, err
	}

	api := client.NewKeysAPI(c)

	return api, nil
}

func ListHosts(conn client.KeysAPI) ([]string, error) {

	var hosts []string

	response, err := conn.Get(context.Background(), "/ping", &client.GetOptions{Recursive: true})

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

func Ping(conn client.KeysAPI, hostname, previous string) (string, error) {

	options := client.SetOptions{
		TTL: 60 * time.Second,
	}

	now := time.Now().String()

	if previous != "" {
		options.PrevValue = previous
	}

	_, err := conn.Set(context.Background(), fmt.Sprintf("/ping/%s", hostname), now, &options)

	if err != nil {
		return now, err
	}

	return now, nil
}

func SubmitState(conn client.KeysAPI, state string) error {

	hosts, err := ListHosts(conn)

	if err != nil {
		return err
	}

	for _, host := range hosts {
		_, err := conn.Set(context.Background(), fmt.Sprintf("/hosts/%s/state", host), state, &client.SetOptions{})
		if err != nil {
			return err
		}
		log.Printf("Loaded state change for host: %s", host)
	}

	return nil
}

func WaitForStateChange(conn client.KeysAPI, hostname string) (*client.Response, error) {

	options := client.WatcherOptions{Recursive: true}
	watcher := conn.Watcher(fmt.Sprintf("/hosts/%s", hostname), &options)

	response, err := watcher.Next(context.Background())

	if err != nil {
		return nil, err
	}

	log.Println(response.Node.Key)

	return response, nil
}
