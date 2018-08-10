package main

import (
	"github.com/cihub/seelog"
	"github.com/fanyang1988/eosforce-db/eosc2db/cmd"
)

func main() {
	defer seelog.Flush()
	cmd.Execute()
}
