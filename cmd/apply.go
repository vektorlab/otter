package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vektorlab/otter/client"
	"github.com/vektorlab/otter/state"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the loaded state to remote Otter daemons",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		stateMap, err := state.StateMapFromYamlPath(GetStatePath(cmd.Flag("state")))
		if err != nil {
			return err
		}
		if cmd.Flag("local").Changed {
			DumpResults(stateMap.Apply())
			return nil
		} else {
			client, err := clients.NewOtterClient(GetEtcdUrls(cmd.Flag("etcd")))
			if err != nil {
				return err
			}
			data, err := stateMap.ToJson()
			if err != nil {
				return err
			}
			client.SubmitState(string(data))
			resultMap, err := client.SubmitCommands("*", "apply")
			if err != nil {
				return err
			}
			DumpResults(resultMap)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(applyCmd)
}
