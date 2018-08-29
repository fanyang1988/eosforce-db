package cmd

import (
	"github.com/fanyang1988/eosforce-db/logsync"

	"github.com/fanyang1988/eosforce-db/p2p-node"

	"github.com/cihub/seelog"
	"github.com/spf13/cobra"
)

var listenCmd = &cobra.Command{
	Use:   "listen [apiUrl] [p2pAddress]",
	Short: "listen eosforce data to log",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		defer seelog.Flush()

		seelog.Infof("start init watch to %v %v", args[0], args[1])

		client, err := p2p_node.NewP2PSyncClient(args[0], args[1])
		if err != nil {
			seelog.Errorf("new p2p sync client err by %v", err.Error())
			return
		}

		client.WithHandler(logsync.NewSync2Log())

		client.StartListen()
		err = <-client.StopChann()
		if err != nil {
			seelog.Errorf("sync err by %v", err.Error())
			return
		}
		seelog.Warnf("sync watch stop!")
	},
}

func init() {
	RootCmd.AddCommand(listenCmd)
}
