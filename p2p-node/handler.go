package p2p_node

import (
	"github.com/fanyang1988/eos-go"
	"github.com/fanyang1988/eos-go/eosforce"
	"github.com/fanyang1988/eos-go/system"
)

// DataHandlerInterface data handler interface, when sync data will call it
type DataHandlerInterface interface {
	OnBlock(block *eos.SignedBlock)
	OnTrx(trx *eos.SignedTransaction)
	OnAction(trx *eos.SignedTransaction, act *eos.Action)

	OnNewAccount(trx *eos.SignedTransaction, data *system.NewAccount)
	OnTransfer(trx *eos.SignedTransaction, data *eosforce.Transfer)

	// TODO
	OnVote(trx *eos.SignedTransaction)
}
