package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/vektorlab/otter/executors"
	"github.com/vektorlab/otter/state"
	"os"
	"strconv"
)

var (
	command  string
	file     string
	dumpJson bool
)

func boolToColor(b bool) *color.Color {
	if b {
		return color.New(color.FgGreen)
	} else {
		return color.New(color.FgHiRed)
	}
}

func dumpResults(executioner *executors.Executioner) {

	td := make([][]string, len(executioner.Executors))

	for _, result := range executioner.Results {
		c := boolToColor(result.Consistent).SprintfFunc()
		td = append(td, []string{
			c(result.Metadata.Name),
			c(result.Metadata.Type),
			c(result.Metadata.State),
			c(strconv.FormatBool(result.Consistent)),
			fmt.Sprint(result.Result),
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Type", "State", "Consistent", "Result"})

	for _, v := range td {
		table.Append(v)
	}
	table.Render()
}

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
	fmt.Println(" state	Show the state of your operating system")
}

func main() {

	options := []string{"c"}
	flag.NewFlagSet("Otter", flag.ExitOnError)
	flag.StringVar(&command, "Command", "", "Otter command [ls, state]")
	flag.StringVar(&file, "c", "otter.yml", "The path to an Otter state file")
	flag.BoolVar(&dumpJson, "json", false, "Dump state output to JSON")

	flag.Parse()
	flag.Usage = func() { usage(options) }

	stateLoader, err := state.FromPath(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch flag.Arg(0) {
	case "ls":
		out, err := stateLoader.Dump()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(out))
	case "state":
		executioner, err := executors.FromStateLoader(stateLoader)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = executioner.Run()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		dumpResults(executioner)

	default:
		flag.Usage()
		os.Exit(1)
	}
}
