package executors

import (
	"fmt"
	"testing"
)

var simpleService = []byte(`{"Name": "docker", "running": true, "Metadata": {"Name": "docker", "Type": "service", "State": "running"}}`)

func TestService(t *testing.T) {
	service, err := ServiceFromJson(simpleService)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if service.state.Running != true {
		fmt.Println("Unable to load service")
		t.Fail()
	}
}

func TestServiceConsistent(t *testing.T) {
	service, err := ServiceFromJson(simpleService)
	if err != nil {
		fmt.Println("Unable to load service: ", err)
		t.Fail()
	}
	fmt.Println(service.Consistent())
}
