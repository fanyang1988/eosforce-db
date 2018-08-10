package main

import (
	"github.com/cihub/seelog"
	"github.com/fanyang1988/eosforce-db/eosc2db/cmd"
)

func init() {
	defaultConfig := `
<seelog>
    <outputs>
        <filter levels="trace">
            <console formatid="common"/>
        </filter>
        <filter levels="debug">
            <console formatid="coloredmagenta"/>
        </filter>
        <filter levels="info">
            <console formatid="coloredyellow"/>
        </filter>
        <filter levels="warn">
            <console formatid="coloredblue"/>
        </filter>
        <filter levels="error,critical">
            <splitter formatid="coloredred">
                <console/>
                <file path="./log/gamex-err.log"/>
            </splitter>
        </filter>
        <file formatid="common" path="./log/gamex.log"/>
    </outputs>
    <formats>
        <format id="coloredblue"  format="[%Date %Time] %EscM(34)[%LEV] [%File(%Line)] [%Func] %Msg%EscM(39)%n%EscM(0)"/>
        <format id="coloredred"  format="[%Date %Time] %EscM(31)[%LEV] [%File(%Line)] [%Func] %Msg%EscM(39)%n%EscM(0)"/>
        <format id="coloredgreen"  format="[%Date %Time] %EscM(32)[%LEV] [%File(%Line)] [%Func] %Msg%EscM(39)%n%EscM(0)"/>
        <format id="coloredyellow"  format="[%Date %Time] %EscM(33)[%LEV] [%File(%Line)] [%Func] %Msg%EscM(39)%n%EscM(0)"/>
        <format id="coloredcyan"  format="[%Date %Time] %EscM(36)[%LEV] [%File(%Line)] [%Func] %Msg%EscM(39)%n%EscM(0)"/>
        <format id="coloredmagenta"  format="[%Date %Time] %EscM(35)[%LEV] [%File(%Line)] [%Func] %Msg%EscM(39)%n%EscM(0)"/>
        <format id="common"  format="[%Date %Time] [%LEV] [%File(%Line)] [%Func] %Msg%n"/>
        <format id="sentry"  format="%Msg%n"/>
    </formats>
</seelog>
	`

	logger, err := seelog.LoggerFromConfigAsBytes([]byte(defaultConfig))
	if err != nil {
		panic(err)
		return
	}
	seelog.ReplaceLogger(logger)
	/*
		_logger = seelog.Current
		SetStackTraceDepth(defaultStackTraceDepth)

		signalhandler.SignalReloadFunc(func() {
			LoadLogConfig(lastCfg)
			fmt.Printf("Got A SIGUSR2 Signal! Now Reloading Conf....\n")
		}) */
}

func main() {
	defer seelog.Flush()
	cmd.Execute()
}
