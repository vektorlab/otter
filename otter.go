package main

import (
	"os"
	"fmt"
	"flag"
	"strconv"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/vektorlab/otter/state"
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

func dumpResults(results []state.Result) {

	td := make([][]string, len(results))

	for _, result := range results {
		c := boolToColor(result.Consistent).SprintfFunc()
		td = append(td, []string{
			c(result.Metadata.Name),
			c(result.Metadata.Type),
			c(result.Metadata.State),
			c(strconv.FormatBool(result.Consistent)),
			fmt.Sprint(result.Message),
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
	fmt.Println(" execute	Execute the state file against your operating system")
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
		err = stateLoader.Consistent()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		dumpResults(stateLoader.Results)
	case "execute":
		err = stateLoader.Execute()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}
}
