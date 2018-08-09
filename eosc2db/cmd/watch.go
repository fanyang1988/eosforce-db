package cmd

import (
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "watch eosforce then sync data to db",
}

var watchInitCmd = &cobra.Command{
	Use:   "init",
	Short: "sync all data from eosforce then watch new change",
	Long: `sync all data from eosforce then watch new change 
it may cost a large time.
`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var watchSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync data from eosforce then watch new change",
	Long: `sync data from eosforce then watch new change,
first will get the last lib block from db, then sync from it.
it may cost a large time.
`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	watchCmd.AddCommand(watchInitCmd)
	watchCmd.AddCommand(watchSyncCmd)
	RootCmd.AddCommand(watchCmd)
}
