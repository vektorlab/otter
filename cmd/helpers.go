package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/pflag"
	"github.com/vektorlab/otter/state"
	"os"
	"strconv"
	"strings"
)

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

func GetEtcdUrls(flag *pflag.Flag) []string {
	if flag.Changed {
		return strings.Split(flag.Value.String(), ",")
	} else {
		return strings.Split(flag.DefValue, ",")
	}
}

func GetStatePath(flag *pflag.Flag) string {
	if flag.Changed {
		return flag.Value.String()
	} else {
		return flag.DefValue
	}
}
