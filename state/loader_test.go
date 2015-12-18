package state

import (
	"encoding/json"
	"fmt"
	"testing"
)

var simple = []byte(`
Really Cool File:
  file.rendered:
    mode: 644
    source: git:///git@github.com/repo/cool_file.txt
    path: file:///root/really-cool-file.txt
docker:
  file.rendered:
    mode: 644
    source: git:///git@github.com/repo/docker.default
    path: /etc/default/docker
  service.running:
    running: true
  package.installed:
    version: 1.9.1
    require:
      - Really Cool File
`)

var processed = []byte(`
{
    "docker-engine": [
        {
            "mode": 644,
            "path": "/etc/default/docker",
            "source": "git:///git@github.com/vektorlab/otter/docker.txt",
            "metadata": {
                "Name": "docker-engine",
                "Type": "file",
                "State": "rendered"
            },
            "require": null
        }
    ]
}
`)

func TestStatesFromYaml(t *testing.T) {
	states, err := StatesFromYaml(simple)
	if err != nil {
		fmt.Println("ERROR:", err)
		t.Fail()
	} else {
		state, err := StatesToJson(states)
		if err != nil {
			fmt.Println("failed to dump state")
			t.Fail()
		}
		fmt.Println(string(state))
		rcf, _ := states["Really Cool File"][0].Dump()
		file := File{}
		json.Unmarshal(rcf, &file)
		if file.Mode != 644 {
			t.Fail()
		}
		if file.Metadata.State != "rendered" {
			t.Fail()
		}
	}
}

func TestStatesFromProcessedJson(t *testing.T) {
	states, err := StatesFromProcessedJson(processed)
	if err != nil {
		fmt.Println("ERROR:", err)
		t.Fail()
	}
	if len(states["docker-engine"]) != 1 {
		fmt.Println("Did not load correct amount of states", states)
	}
}
