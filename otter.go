package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/vektorlab/otter/daemon"
	"github.com/vektorlab/otter/state"
	"os"
	"strconv"
	"strings"
)

var (
	command  string
	file     string
	dumpJson bool
	etcdUrl  string
)

func boolToColor(b bool) *color.Color {
	if b {
		return color.New(color.FgGreen)
	} else {
		return color.New(color.FgHiRed)
	}
}

func getLoader() *state.Loader {
	loader, err := state.FromPathToYaml(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return loader
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

func dumpHosts(hosts []string) {
	td := make([][]string, len(hosts))

	for _, host := range hosts {
		c := boolToColor(true).SprintfFunc()
		td = append(td, []string{c(host)})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"host"})

	for _, v := range td {
		table.Append(v)
	}
	table.Render()
}

func usage(options []string) {
	fmt.Println("Otter is an opinionated configuration management framework for servers that run containers\n")
	fmt.Println("Usage: otter [OPTIONS] [load, ls, state, execute, daemon] \n")
	fmt.Println("Options:")
	for i := 0; i < len(options); i++ {
		f := flag.Lookup(options[i])
		if f != nil {
			fmt.Printf(" --%s		%s [%s]\n", f.Name, f.Usage, f.DefValue)
		}
	}
	fmt.Println("\nCommands:")
	fmt.Println(" load	Load a state configuration into etcd")
	fmt.Println(" ls	List all hosts registered in an etcd cluster")
	fmt.Println(" state	Show the state of your operating system")
	fmt.Println(" execute	Execute the state file against your operating system")
	fmt.Println(" daemon	Run Otter in daemon mode")
}

func main() {

	options := []string{"c", "e"}
	flag.NewFlagSet("Otter", flag.ExitOnError)
	flag.StringVar(&command, "Command", "", "Otter command [ls, state, execute, daemon]")
	flag.StringVar(&file, "c", "otter.yml", "The path to an Otter state file")
	flag.StringVar(&etcdUrl, "e", "http://127.0.0.1:2379", "URL to etcd hosts")
	flag.BoolVar(&dumpJson, "json", false, "Dump state output to JSON")

	flag.Parse()
	flag.Usage = func() { usage(options) }

	switch flag.Arg(0) {
	case "load":
		loader := getLoader()
		out, err := loader.Dump()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		daemon, err := daemon.NewDaemon(strings.Split(etcdUrl, ","))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = daemon.LoadState(string(out))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "ls":
		daemon, err := daemon.NewDaemon(strings.Split(etcdUrl, ","))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		hosts, err := daemon.ListHosts()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		dumpHosts(hosts)
	case "state":
		loader := getLoader()
		err := loader.Consistent()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		dumpResults(loader.Results)
	case "execute":
		loader := getLoader()
		err := loader.Execute()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "daemon":
		daemon, err := daemon.NewDaemon(strings.Split(etcdUrl, ","))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		daemon.Run()
	default:
		flag.Usage()
		os.Exit(1)
	}
}
