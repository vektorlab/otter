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

func TestFromBytes(t *testing.T) {

	loader, err := FromBytes(simple)
	if err != nil {
		fmt.Println("ERROR:", err)
		t.Fail()
	} else {
		state, err := loader.Dump()
		if err != nil {
			fmt.Println("failed to dump state")
			t.Fail()
		}
		fmt.Println(string(state))
		rcf, _ := loader.State["Really Cool File"][0].Dump()
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

func TestMissingRequirement(t *testing.T) {

	_, err := FromBytes(missing)

	if err == nil {
		fmt.Println("Failed to detect missing requirement")
		t.Fail()
	}
}

func TestCircularRequirement(t *testing.T) {
	_, err := FromBytes(circular)

	if err == nil {
		fmt.Println("Failed to detect circular requirement")
		t.Fail()
	}
}
