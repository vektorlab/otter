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
func (otter *Otter) SubmitCommand(host, command string, resultChan chan []*state.Result) {
	key := fmt.Sprintf("/command/%s/%s", host, command)
	id := helpers.RandomString(8)
	_, err := otter.etcdKeysApi.Set(context.Background(), key, id, &etcd.SetOptions{})
	log.Printf("Submitted command with keyspace: %s", key)
	if err != nil {
		resultChan <- state.ResultsFromError(host, err)
	}
	results, err := otter.WaitForResults(id)
	if err != nil {
		resultChan <- state.ResultsFromError(host, err)
	}
	for _, result := range results {
		result.Host = host
		if result.Message == "" {
			result.Message = "None"
		}
	}
	resultChan <- results
}

/*
Submit a single command against all hosts who's name matches the given hostString.
*/
func (otter *Otter) SubmitCommands(hostString, command string) (*state.ResultMap, error) {
	hosts, err := otter.ListHosts() // TODO: Match hostString
	resultChan := make(chan []*state.Result)
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
func (otter *Otter) CollectResults(results chan []*state.Result, count int) *state.ResultMap {
	resultMap := state.NewResultMap()
	for i := 0; i < count; i++ {
		log.Printf("Recieved %d of %d results sets", i, count)
		r := <-results
		for _, result := range r {
			log.Printf("[%s] Saving result %s", result.Host, result.Metadata.Name)
			resultMap.Add(result)
		}
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
