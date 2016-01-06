package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vektorlab/otter/client"
	"github.com/vektorlab/otter/state"
)

// stateCmd represents the state command
var stateCmd = &cobra.Command{
	Use:   "state",
	Short: "Show the state of remote hosts in the cluster",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		stateMap, err := state.StateMapFromYamlPath(GetStatePath(cmd.Flag("state")))
		if err != nil {
			return err
		}
		if cmd.Flag("local").Changed {
			DumpResults(stateMap.Consistent())
			return nil
		} else {
			client, err := clients.NewOtterClient(GetEtcdUrls(cmd.Flag("etcd")))
			data, err := stateMap.ToJson()
			if err != nil {
				return err
			}
			client.SubmitState(string(data))
			if err != nil {
				return err
			}
			resultMap, err := client.SubmitCommands("*", "state")
			if err != nil {
				return err
			}
			DumpResults(resultMap)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(stateCmd)
}
