package cmd

import (
	"time"

	"github.com/cihub/seelog"
	"github.com/fanyang1988/eosforce-db/p2p-node"
	"github.com/fanyang1988/eosforce-db/pgsync"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "watch eosforce then sync data to db",
}

var watchInitCmd = &cobra.Command{
	Use:   "init [apiUrl] [p2pAddress]",
	Short: "sync all data from eosforce then watch new change",
	Long: `sync all data from eosforce then watch new change 
it may cost a large time.
`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		defer seelog.Flush()

		seelog.Infof("start init watch to %v %v", args[0], args[1])

		client, err := p2p_node.NewP2PSyncClient(args[0], args[1])
		if err != nil {
			seelog.Errorf("new p2p sync client err by %v", err.Error())
			return
		}

		client.WithHandler(pgsync.NewSyncPgDB(
			viper.GetString("db-address"),
			viper.GetString("db-user"),
			viper.GetString("db-passwd"),
			viper.GetString("db")))

		client.StartSync(0, make([]byte, 32), time.Now(), 0, make([]byte, 32))
		err = <-client.StopChann()
		if err != nil {
			seelog.Errorf("sync err by %v", err.Error())
			return
		}
		seelog.Warnf("sync watch stop!")
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
