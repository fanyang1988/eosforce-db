package cmd

import (
	"fmt"
	"os"

	"github.com/cihub/seelog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringP("log-cfg", "", "./log_cfg.xml", "seelog cfg file")

	RootCmd.PersistentFlags().StringP("db-address", "", "127.0.0.1:5432", "db address")
	RootCmd.PersistentFlags().StringP("db-user", "", "pgfy", "db user")
	RootCmd.PersistentFlags().StringP("db-passwd", "", "123456", "db password")
	RootCmd.PersistentFlags().StringP("db", "", "test1", "db")

	for _, flag := range []string{"db-address", "db-user", "db-passwd", "db", "log-cfg"} {
		if err := viper.BindPFlag(flag, RootCmd.PersistentFlags().Lookup(flag)); err != nil {
			panic(err)
		}
	}

	configFilePath := viper.GetString("log-cfg")
	if configFilePath != "" && IsFileExists(configFilePath) {
		logger, err := seelog.LoggerFromConfigAsFile(configFilePath)
		if err != nil {
			panic(err)
		}
		seelog.ReplaceLogger(logger)
	}
}

func initConfig() {
}

// Exists reports whether the named file or directory exists.
func IsFileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
