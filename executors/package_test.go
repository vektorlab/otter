package executors

import (
	"fmt"
	"testing"
)

var simplePkg = []byte(`{"Name": "docker-engine", "version": "1.9.1", "Metadata": {"Name": "docker", "Type": "package", "State": "installed"}}`)

func TestPackage(t *testing.T) {
	pkg := Package{}
	pkg.Load(simplePkg)
	if pkg.state.Version != "1.9.1" {
		fmt.Println("Unable to load package")
		t.Fail()
	}
}

func TestPackageConsistent(t *testing.T) {
	pkg := Package{}
	err := pkg.Load(simplePkg)
	if err != nil {
		fmt.Println("Unable to load package: ", err)
		t.Fail()
	}
	fmt.Println(pkg.Consistent())
}
