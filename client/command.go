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
func (otter *Otter) SubmitCommand(host, command string, resultChan chan *state.ResultMap) {
	key := fmt.Sprintf("/command/%s/%s", host, command)
	id := helpers.RandomString(8)
	_, err := otter.etcdKeysApi.Set(context.Background(), key, id, &etcd.SetOptions{})
	log.Printf("Submitted command with keyspace: %s", key)
	if err != nil {
		resultChan <- state.ResultMapFromError(host, err)
		return
	}
	results, err := otter.WaitForResults(id)
	if err != nil {
		resultChan <- state.ResultMapFromError(host, err)
		return
	}
	resultChan <- results
}

/*
Submit a single command against all hosts who's name matches the given hostString.
*/
func (otter *Otter) SubmitCommands(hostString, command string) (*state.ResultMap, error) {
	hosts, err := otter.ListHosts() // TODO: Match hostString
	resultChan := make(chan *state.ResultMap)
	if err != nil {
		return nil, err
	}
	for _, host := range hosts {
		go otter.SubmitCommand(host, command, resultChan)
	}
	return otter.CollectResults(resultChan, len(hosts)), nil
}

/*
Collect results from a channel and append them to a ResultMap when they are received
*/
func (otter *Otter) CollectResults(results chan *state.ResultMap, count int) *state.ResultMap {
	resultMap := state.NewResultMap()
	for i := 0; i < count; i++ {
		log.Printf("Recieved %d of %d results sets", i, count)
		r := <-results
		resultMap.Merge(r)
	}
	return resultMap
}

/*
Wait for a command and then return it.
*/
func (otter *Otter) WaitForCommand(hostname string) (string, string, error) {
	key := fmt.Sprintf("/command/%s")
	key, id, err := otter.WaitForChange(fmt.Sprintf("/command/%s", hostname), true, 0*time.Second)
	if err != nil {
		return "", "", err
	}
	split := strings.Split(key, "/")
	return split[len(split)-1], id, nil
}
