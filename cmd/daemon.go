package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vektorlab/otter/daemon"
)

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run the Otter client in daemon mode",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		daemon, err := daemon.NewDaemon(GetEtcdUrls(cmd.Flag("etcd")))
		if err != nil {
			return err
		}
		daemon.Run()
		return nil
	},
}

func init() {
	RootCmd.AddCommand(daemonCmd)
}
