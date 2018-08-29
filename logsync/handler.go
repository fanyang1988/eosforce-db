package logsync

import (
	"time"

	"github.com/cihub/seelog"
	"github.com/fanyang1988/eos-go"
	"github.com/fanyang1988/eos-go/eosforce"
	"github.com/fanyang1988/eos-go/system"
)

// Sync2Log just sync data to log
type Sync2Log struct {
}

var lastBlockCount = time.Now().UnixNano()

func (s *Sync2Log) OnBlock(blockID string, block *eos.SignedBlock) {
	curr := time.Now().UnixNano()

	seelog.Debugf("on block ...%s %d conf:%d trx:%d time:%v %v",
		blockID[len(blockID)-6:], block.BlockNumber(), block.Confirmed, len(block.Transactions),
		(curr-lastBlockCount)/1000000, block.Producer)

	lastBlockCount = curr
}

func (s *Sync2Log) OnTrx(blockID string, blk *eos.SignedBlock, trx *eos.SignedTransaction) {
	seelog.Tracef("on trx %v %v", trx.ID(), trx.String())
}

func (s *Sync2Log) OnAction(blockID string, trx *eos.SignedTransaction, act *eos.Action) {
	data, _ := act.MarshalJSON()
	seelog.Tracef("on act %v %v %v", trx.ID(), trx.String(), string(data))
}

func (s *Sync2Log) OnNewAccount(blockID string, trx *eos.SignedTransaction, act *eos.Action, data *system.NewAccount) {
	seelog.Infof("on new account %v", *data)
}

func (s *Sync2Log) OnTransfer(blockID string, trx *eos.SignedTransaction, act *eos.Action, data *eosforce.Transfer) {
	seelog.Infof("transfer asset %s : %s --> %s by %s", data.Quantity, data.From, data.To, data.Memo)
}

// TODO
func (s *Sync2Log) OnVote(blockID string, trx *eos.SignedTransaction) {

}

// NewSync2Log new sync to log
func NewSync2Log() *Sync2Log {
	res := &Sync2Log{}
	return res
}
