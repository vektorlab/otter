package state

import (
	"fmt"
	"testing"
)

var simpleFile = []byte(`
{"path": "/tmp/no-exist", "mode": 644, "source": "git:///git@github.com:vektorlab/otter/README.md"}
`)

var simpleFileMeta = Metadata{
	Name: "Simple File",
	Type: "file",
	State: "rendered",
}

func TestStateFactory(t *testing.T) {
	state, err := StateFactory(simpleFileMeta, simpleFile)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	name := state.Meta().Name
	if name != "Simple File" {
		fmt.Printf("Did not load metadata correctly: %s", name)
		t.Fail()
	}
}
