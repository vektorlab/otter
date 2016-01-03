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
	fmt.Println("Usage: otter [OPTIONS] [apply, load, ls, daemon] \n")
	fmt.Println("Flags:")
	for i := 0; i < len(options); i++ {
		f := flag.Lookup(options[i])
		if f != nil {
			fmt.Printf(" -%s		%s [%s]\n", f.Name, f.Usage, f.DefValue)
		}
	}
	fmt.Println("\nCommands:")
	fmt.Println(" apply	Execute the state file against remote hosts")
	fmt.Println(" daemon	Run Otter in daemon mode")
	fmt.Println(" ls	List all hosts registered to the cluster")
	fmt.Println(" load	Load a state configuration into the cluster")
}

func Parse() (string, string, []string) {

	var (
		command string
		path    string
		urls    string
	)

	options := []string{"c", "e"}

	flag.NewFlagSet("Otter", flag.ExitOnError)
	flag.StringVar(&command, "Command", "", "Otter command [ls, load, execute, daemon]")
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

func isConsistent(results []state.Result) bool {
	for _, result := range results {
		if !result.Consistent {
			return false
		}
	}
	return true
}

func boolToColor(b bool) *color.Color {
	if b {
		return color.New(color.FgGreen)
	} else {
		return color.New(color.FgHiRed)
	}
}

func DumpResults(resultMap *state.ResultMap) {
	table := tablewriter.NewWriter(os.Stdout)
	tableData := make([][]string, len(resultMap.Results))
	for host, results := range resultMap.Results {
		for _, result := range results {
			c := boolToColor(result.Consistent).SprintfFunc()
			tableData = append(tableData, []string{
				c(host),
				c(result.Metadata.Name),
				c(result.Metadata.Type),
				c(result.Metadata.State),
				c(strconv.FormatBool(result.Consistent)),
				fmt.Sprint(result.Message),
			})
		}
	}
	for _, v := range tableData {
		table.Append(v)
	}
	table.SetHeader([]string{"Host", "Name", "Type", "State", "Consistent", "Result"})
	table.Render()
}

func DumpHosts(hosts map[string]bool) {
	td := make([][]string, len(hosts))

	for host, consistent := range hosts {
		c := boolToColor(consistent).SprintfFunc()
		td = append(td, []string{
			c(host),
			c(strconv.FormatBool(consistent)),
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Host", "Consistent"})

	for _, v := range td {
		table.Append(v)
	}
	table.Render()
}
