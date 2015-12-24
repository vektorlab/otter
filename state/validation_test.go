package state

import (
	"fmt"
	"testing"
)

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

func TestValidate(t *testing.T) {

	states, err := StatesFromYaml(missing)

	if err != nil {
		fmt.Println("Failed to load YAML")
		t.Fail()
	}
	err = Validate(states)

	if err == nil {
		fmt.Println("Failed to detect missing requirement")
		t.Fail()
	}

	states, err = StatesFromYaml(circular)

	if err != nil {
		fmt.Println("Failed to load YAML")
		t.Fail()
	}

	err = Validate(states)

	if err == nil {
		fmt.Println("Failed to detect circular requirement")
		t.Fail()
	}
}
