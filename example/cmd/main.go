package main

import (
	"flag"
	"log"

	"encoding/json"
	"fmt"

	"github.com/cihub/seelog"
	"github.com/fanyang1988/eos-go"
	"github.com/fanyang1988/eos-go/eosforce"
	"github.com/fanyang1988/eos-go/p2p"
	"github.com/fanyang1988/eos-go/system"
	"github.com/fanyang1988/eosforce-db/model"
	"github.com/go-pg/pg"
)

var p2pAddr = flag.String("p2p-addr", "localhost:9001", "P2P socket connection")
var apiUrl = flag.String("api-url", "http://127.0.0.1:8888", "API Url Address")
var networkVersion = flag.Int("network-version", 1, "Network version")

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

	flag.Parse()

	done := make(chan bool)

	api := eos.New(*apiUrl)
	info, err := api.GetInfo()
	if err != nil {
		log.Fatal("Error getting info: ", err)
	}
	cID := info.ChainID

	client := p2p.NewClient(*p2pAddr, cID, uint16(*networkVersion))
	if err != nil {
		log.Fatal(err)
	}
	client.WithLogger(seelog.Current)

	pgDB := pg.Connect(&pg.Options{
		Addr:     "127.0.0.1:5432",
		User:     "pgfy",
		Password: "123456",
		Database: "test4",
	})

	client.RegisterHandler(p2p.HandlerFunc(func(msg p2p.Message) {
		data, error := json.Marshal(msg)
		if error != nil {
			fmt.Println("logger plugin err: ", error)
			return
		}

		switch msg.Envelope.Type {
		case eos.SignedBlockType:
			{
				signedBlockMsg, ok := msg.Envelope.P2PMessage.(*eos.SignedBlock)
				if !ok {
					seelog.Errorf("typ error by signedBlockMsg")
					return
				}

				blockID, err := signedBlockMsg.BlockID()
				if err != nil {
					seelog.Errorf("BlockID error by signedBlockMsg %v", err.Error())
					return
				}

				seelog.Infof("on block %d %s %v", signedBlockMsg.BlockNumber(), signedBlockMsg.Producer, blockID)

				pgDB.Insert(&model.Block{
					BlockID:     blockID.String(),
					BlockNum:    int64(signedBlockMsg.BlockNumber()),
					ProduceTime: signedBlockMsg.Timestamp.Time,
					Producer:    string(signedBlockMsg.Producer),
					PrevBlockID: signedBlockMsg.Previous.String(),
				})

				tt, _ := signedBlockMsg.Timestamp.MarshalJSON()
				seelog.Tracef("block %v %s", signedBlockMsg.Timestamp.Time, string(tt))

				for tidx, tr := range signedBlockMsg.Transactions {
					seelog.Infof("tidx %d : %v %v", tidx, tr.Status)
					trx, err := tr.Transaction.Packed.Unpack()
					if err != nil {
						seelog.Errorf("transaction unpack err by %s", err.Error())
						continue
					}

					seelog.Infof("trx %v %v", trx.Fee.String(), trx.String())

					for aidx, action := range trx.Actions {
						seelog.Infof("on action %d %d %v %v --> %v", tidx, aidx, action.Account, action.Name, action.ActionData)
						switch action.Account {
						case "eosio":
							{
								switch action.Name {
								case "transfer":
									{
										transferAct, ok := action.ActionData.Data.(*eosforce.Transfer)
										if !ok {
											seelog.Errorf("transfer act data err")
											continue
										}

										err := pgDB.Insert(&model.Transfer{
											FromAccount: string(transferAct.From),
											ToAccount:   string(transferAct.To),
											Quantity:    transferAct.Quantity.Amount,
											Token:       transferAct.Quantity.Symbol.Symbol,
											Memo:        transferAct.Memo,
											RefBlockNum: int64(trx.RefBlockNum),
										})

										if err != nil {
											seelog.Errorf("create transfer err by %s", err.Error())
											continue
										}

										seelog.Infof("transfer %v", *transferAct)
									}
								case "newaccount":
									{
										newAccountAct, ok := action.ActionData.Data.(*system.NewAccount)
										if !ok {
											seelog.Errorf("newAccountAct act data err")
											continue
										}

										err := pgDB.Insert(&model.Accounts{
											Name:        string(newAccountAct.Name),
											CreateAt:    signedBlockMsg.Timestamp.Time,
											UpdateAt:    signedBlockMsg.Timestamp.Time,
											Creater:     string(newAccountAct.Creator),
											Data:        *action,
											RefBlockNum: int64(trx.RefBlockNum),
										})

										if err != nil {
											seelog.Errorf("create account err by %s", err.Error())
											continue
										}

										seelog.Infof("newAccountAct %v", *newAccountAct)
									}
								}
							}
						}
					}
				}

				return
			}
		}

		seelog.Infof("recv msg from %s --> %s", msg.Route.From, string(data))

	}))
	//time.Sleep(5 * time.Second)

	err = client.ConnectSyncAll()
	//err = client.ConnectRecent()
	if err != nil {
		log.Fatal(err)
	}

	<-done

}
