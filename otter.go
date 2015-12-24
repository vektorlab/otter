package main

import (
	"fmt"
	"github.com/vektorlab/otter/cli"
	"os"
)

func main() {

	command, path, urls := cli.Parse()

	cli, err := cli.NewOtterCLI(command, path, urls)

	if err != nil {
		fmt.Println("ERROR: ", err)
		os.Exit(1)
	}

	err = cli.Run()

	if err != nil {
		fmt.Println("ERROR: ", err)
	}
}
