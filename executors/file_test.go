package executors

import (
	"fmt"
	"testing"
)

var simpleFile = []byte(`{"Name": "/tmp/testfile", "mode": 644, "Metadata": {"Name": "Test File", "Type": "file", "State": "rendered"}}`)

func TestFile(t *testing.T) {
	file := File{}
	file.Load(simpleFile)
	if file.state.Mode != 644 {
		fmt.Println("Unable to load file")
		t.Fail()
	}
}

func TestFileConsistent(t *testing.T) {
	file := File{}
	err := file.Load(simpleFile)
	if err != nil {
		fmt.Println("Unable to load package: ", err)
		t.Fail()
	}
	fmt.Println(file.Consistent())
}
