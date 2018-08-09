package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "eosc2db",
	Short: "eosc2db is tool to dump data from eosforce to db",
	Long: `eosc2db is tool to dump data from eosforce to db
Source code is available at: https://github.com/fanyang1988/eosforce-db
`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
