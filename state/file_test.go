package state

import (
	"fmt"
	"io/ioutil"
	"testing"
)

var fileHttpRender = []byte(`
{"path": "/tmp/otter-fileHttpRender", "mode": "0644", "source": "http://otter.vektor.nyc/tests/file-render-test.txt"}
`)

var fileHttpRenderMeta = Metadata{
	Name:  "Simple File From HTTP",
	Type:  "file",
	State: "rendered",
}

func compareFile(path, other string, t *testing.T) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Failed to read file: %s", err.Error())
		t.Fail()
	}
	if string(data) != other {
		fmt.Println("%s != %s", string(data), other)
		t.Fail()
	}
}

func TestFileConsistent(t *testing.T) {
	state := stateSetup(simpleFileMeta, simpleFile, t)
	result := state.Consistent()
	if result.Consistent != false {
		fmt.Println("Detected non-existant file: ", result.Metadata.Name)
		t.Fail()
	}
}

func TestRenderFromHTTP(t *testing.T) {
	state := stateSetup(fileHttpRenderMeta, fileHttpRender, t)
	result := state.Execute()
	if result.Consistent != true {
		fmt.Println("Failed to render file: ", result.Message)
		t.Fail()
	}
	compareFile("/tmp/otter-fileHttpRender", "Otter Test!\n", t)
}
