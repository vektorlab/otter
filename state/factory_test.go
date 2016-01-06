package state

import (
	"fmt"
	"testing"
)

var simpleFile = []byte(`
{"path": "/tmp/no-exist", "mode": "644", "source": "git:///git@github.com:vektorlab/otter/README.md"}
`)

var simpleFileMeta = Metadata{
	Name:  "non-existant-file",
	Type:  "file",
	State: "rendered",
}

var simplePackage = []byte(`
{"version": "1.1.1"}
`)

var simplePackageMeta = Metadata{
	Name:  "non-existant-package",
	Type:  "package",
	State: "installed",
}

var simpleService = []byte(`{}`)

var simpleServiceMeta = Metadata{
	Name:  "non-existant-service",
	Type:  "service",
	State: "running",
}

func stateSetup(metadata Metadata, stateBytes []byte, t *testing.T) State {
	state, err := StateFactory(metadata, stateBytes)
	if err != nil {
		fmt.Println("Unable to load file: ", err)
		t.Fail()
	}
	err = state.Initialize()
	if err != nil {
		fmt.Println("File failed to initialize: ", err)
	}
	return state
}

func TestStateFactoryFile(t *testing.T) {
	file := stateSetup(simpleFileMeta, simpleFile, t)
	name := file.Meta().Name
	if name != "non-existant-file" {
		fmt.Printf("Failed to setup file: %s", name)
		t.Fail()
	}
}

func TestStateFactoryPackage(t *testing.T) {
	pkg := stateSetup(simplePackageMeta, simplePackage, t)
	name := pkg.Meta().Name
	if name != "non-existant-package" {
		fmt.Printf("Failed to setup package: %s", name)
	}
}

func TestStateFactoryService(t *testing.T) {
	service := stateSetup(simpleServiceMeta, simpleService, t)
	name := service.Meta().Name
	if name != "non-existant-service" {
		fmt.Printf("Failed to setup service: %s", name)
		t.Fail()
	}
}
