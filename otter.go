package main

import (
	"flag"
	"fmt"
	"github.com/vektorlab/otter/state"
	"os"
)

var (
	command string
	file    string
)

func usage(options []string) {
	fmt.Println("Otter is an opinionated configuration management framework for servers that run containers\n")
	fmt.Println("Usage: otter [OPTIONS] [ls] \n")
	fmt.Println("Options:")
	for i := 0; i < len(options); i++ {
		f := flag.Lookup(options[i])
		if f != nil {
			fmt.Printf(" --%s		%s [%s]\n", f.Name, f.Usage, f.DefValue)
		}
	}
	fmt.Println("\nCommands:")
	fmt.Println(" ls	Output the state from the Otter configuration file")
}

func main() {

	options := []string{"state"}
	flag.NewFlagSet("Otter", flag.ExitOnError)
	flag.StringVar(&command, "Command", "", "Otter command [ls]")
	flag.StringVar(&file, "state", "otter.yml", "The path to an Otter state file")

	flag.Parse()
	flag.Usage = func() { usage(options) }

	switch flag.Arg(0) {
	case "ls":
		loader, err := state.FromPath(file)
		if err != nil {
			fmt.Println("ERROR:", err)
			os.Exit(1)
		}
		state, err := loader.Dump()
		if err != nil {
			fmt.Println("Problem dumping state")
		}
		fmt.Println(string(state))

	default:
		flag.Usage()
		os.Exit(1)
	}
}
