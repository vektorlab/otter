package cli

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/vektorlab/otter/state"
	"os"
	"strconv"
	"strings"
)

func usage(options []string) {
	fmt.Println("Otter is an opinionated configuration management framework for servers that run containers\n")
	fmt.Println("Usage: otter [OPTIONS] [load, ls, state, execute, daemon] \n")
	fmt.Println("Flags:")
	for i := 0; i < len(options); i++ {
		f := flag.Lookup(options[i])
		if f != nil {
			fmt.Printf(" -%s		%s [%s]\n", f.Name, f.Usage, f.DefValue)
		}
	}
	fmt.Println("\nCommands:")
	fmt.Println(" daemon	Run Otter in daemon mode")
	fmt.Println(" execute	Execute the state file against remote hosts")
	fmt.Println(" ls	List all hosts registered to the cluster")
	fmt.Println(" load	Load a state configuration into the cluster")
	fmt.Println(" state	Show the state of remote hosts in the cluster")
}

func Parse() (string, string, []string) {

	var (
		command string
		path    string
		urls    string
	)

	options := []string{"c", "e"}

	flag.NewFlagSet("Otter", flag.ExitOnError)
	flag.StringVar(&command, "Command", "", "Otter command [ls, state, execute, daemon]")
	flag.StringVar(&path, "c", "otter.yml", "The path to an Otter state file")
	flag.StringVar(&urls, "e", "http://127.0.0.1:2379", "URL to etcd hosts")

	flag.Parse()
	flag.Usage = func() { usage(options) }

	etcdUrls := strings.Split(urls, ",")

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	command = flag.Args()[0]

	return command, path, etcdUrls

}

func boolToColor(b bool) *color.Color {
	if b {
		return color.New(color.FgGreen)
	} else {
		return color.New(color.FgHiRed)
	}
}

func DumpResults(results []state.Result) {

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

func DumpHosts(hosts []string) {
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
