package p2p_node

import (
	"github.com/fanyang1988/eos-go"
	"github.com/fanyang1988/eos-go/eosforce"
	"github.com/fanyang1988/eos-go/system"
)

// DataHandlerInterface data handler interface, when sync data will call it
type DataHandlerInterface interface {
	OnBlock(blockID string, block *eos.SignedBlock)
	OnTrx(blockID string, blk *eos.SignedBlock, trx *eos.SignedTransaction)
	OnAction(blockID string, trx *eos.SignedTransaction, act *eos.Action)

	OnNewAccount(blockID string, trx *eos.SignedTransaction, act *eos.Action, data *system.NewAccount)
	OnTransfer(blockID string, trx *eos.SignedTransaction, act *eos.Action, data *eosforce.Transfer)

	// TODO
	OnVote(blockID string, trx *eos.SignedTransaction)
}
