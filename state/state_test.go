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

var missing = []byte(`
docker:
  package.installed:
    version: 1.9.1
mesos:
  package.installed:
    require:
      - docker
      - zookeeper
`)

var circular = []byte(`
docker:
  package.installed:
    version: 1.9.1
    require:
      - mesos
mesos:
  package.installed:
    require:
      - docker
`)

func TestStateMapFromYaml(t *testing.T) {
	stateMap, err := StateMapFromYaml(simple)
	if err != nil {
		fmt.Println("ERROR:", err)
		t.Fail()
	} else {
		_, err := stateMap.ToJson()
		if err != nil {
			fmt.Println("failed to dump state")
			t.Fail()
		}
		rcf, _ := stateMap.States["Really Cool File"][0].Dump()
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

func TestStateMapFromProcessedJson(t *testing.T) {
	stateMap, err := StateMapFromProcessedJson(processed)
	if err != nil {
		fmt.Println("ERROR:", err)
		t.Fail()
	}
	if len(stateMap.States["docker-engine"]) != 1 {
		fmt.Println("Did not load correct amount of states", stateMap.States)
	}
}

func TestStateValidation(t *testing.T) {
	stateMap, err := StateMapFromYaml(missing)
	if err != nil {
		fmt.Println("Failed to load YAML")
		t.Fail()
	}
	err = stateMap.Validate()
	if err == nil {
		fmt.Println("Failed to detect missing requirement")
		t.Fail()
	}
	stateMap, err = StateMapFromYaml(circular)
	if err != nil {
		fmt.Println("Failed to load YAML")
		t.Fail()
	}
	err = stateMap.Validate()
	if err == nil {
		fmt.Println("Failed to detect circular requirement")
		t.Fail()
	}
}
