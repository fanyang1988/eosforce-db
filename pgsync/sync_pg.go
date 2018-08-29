package pgsync

import (
	"time"

	"github.com/cihub/seelog"
	"github.com/fanyang1988/eos-go"
	"github.com/fanyang1988/eos-go/eosforce"
	"github.com/fanyang1988/eos-go/system"
	"github.com/fanyang1988/eosforce-db/model"
	"github.com/go-pg/pg"
)

type Sync2pgDB struct {
	pgOpt pg.Options
	pgDB  *pg.DB
}

var lastBlockCount int64

func (s *Sync2pgDB) OnBlock(blockID string, block *eos.SignedBlock) {
	curr := time.Now().UnixNano()
	seelog.Debugf("on block ...%s %d conf:%d trx:%d time:%v %v",
		blockID[len(blockID)-6:], block.BlockNumber(), block.Confirmed, len(block.Transactions),
		(curr-lastBlockCount)/1000000, block.Producer)
	lastBlockCount = curr

	b := &model.Block{
		BlockID:               blockID,
		BlockNum:              int64(block.BlockNumber()),
		ProduceTime:           block.Timestamp.Time,
		Producer:              string(block.Producer),
		PrevBlockID:           block.Previous.String(),
		TransactionMerkleRoot: block.TransactionMRoot.String(),
		ActionMerkleRoot:      block.ActionMRoot.String(),
		NumTransactions:       int(len(block.Transactions)),
		Confirmed:             int(block.Confirmed),
	}

	// TODO use auto gen set sql
	_, err := s.pgDB.Model(b).OnConflict("(block_num) DO UPDATE").
		Set(
			"produce_time = EXCLUDED.produce_time,confirmed=EXCLUDED.confirmed," +
				"block_id = EXCLUDED.block_id,prev_block_id = EXCLUDED.prev_block_id," +
				"producer = EXCLUDED.producer,transaction_merkle_root = EXCLUDED.transaction_merkle_root," +
				"action_merkle_root = EXCLUDED.action_merkle_root,num_transactions = EXCLUDED.num_transactions").
		Insert()

	if err != nil {
		seelog.Errorf("insert account err by %s", err.Error())
	}
}

func (s *Sync2pgDB) OnTrx(blockID string, blk *eos.SignedBlock, trx *eos.SignedTransaction) {
	//seelog.Tracef("on trx %v %v", trx.ID(), trx.String())

	t := &model.Transaction{
		Id:             trx.ID(),
		RefBlockNum:    int64(blk.BlockNumber()),
		BlockID:        blockID,
		RefBlockPrefix: int64(trx.RefBlockPrefix),
		Expiration:     trx.Expiration.Time,
		NumActions:     int64(len(trx.Actions)),
		DelaySec:       int64(trx.DelaySec),
		Fee:            int64(trx.Fee.Amount),
		Data:           *trx,
	}

	err := s.pgDB.Insert(t)

	if err != nil {
		seelog.Errorf("insert trx err by %s", err.Error())
		return
	}

}

func (s *Sync2pgDB) OnAction(blockID string, trx *eos.SignedTransaction, act *eos.Action) {
	data, _ := act.MarshalJSON()
	seelog.Tracef("on act %v %v %v", trx.ID(), trx.String(), string(data))

	a := &model.Action{
		Account:        string(act.Account),
		RefBlockNum:    int64(trx.RefBlockNum),
		BlockID:        blockID,
		RefBlockPrefix: int64(trx.RefBlockPrefix),
		TrxID:          trx.ID(),
		Name:           string(act.Name),
		Fee:            int64(trx.Fee.Amount),
		Data:           *act,
	}

	err := s.pgDB.Insert(a)

	if err != nil {
		seelog.Errorf("insert act err by %s", err.Error())
		return
	}
}

func (s *Sync2pgDB) OnNewAccount(blockID string, trx *eos.SignedTransaction, act *eos.Action, data *system.NewAccount) {
	seelog.Infof("on new account %v", *data)
	_, err := s.pgDB.Model(&model.Accounts{
		Name:        string(data.Name),
		CreateAt:    trx.Expiration.Time,
		UpdateAt:    trx.Expiration.Time,
		Creater:     string(data.Creator),
		RefBlockNum: int64(trx.RefBlockNum),
		Data:        *act,
	}).OnConflict("(name) DO UPDATE").
		Set("create_at = EXCLUDED.create_at,update_at = EXCLUDED.update_at," +
			"creater = EXCLUDED.creater,ref_block_num = EXCLUDED.ref_block_num,data = EXCLUDED.data").
		Insert()

	if err != nil {
		seelog.Errorf("insert account err by %s", err.Error())
	}

	for _, key := range data.Owner.Keys {
		err := s.pgDB.Insert(&model.AccountPermission{
			Account:    string(data.Name),
			Permission: "Owner",
			Pubkey:     string(key.PublicKey.String()),
		})
		if err != nil {
			seelog.Errorf("insert pubkey err by %s", err.Error())
		}
	}

	for _, key := range data.Active.Keys {
		err := s.pgDB.Insert(&model.AccountPermission{
			Account:    string(data.Name),
			Permission: "Active",
			Pubkey:     string(key.PublicKey.String()),
		})
		if err != nil {
			seelog.Errorf("insert pubkey err by %s", err.Error())
		}
	}
}

func (s *Sync2pgDB) OnTransfer(blockID string, trx *eos.SignedTransaction, act *eos.Action, data *eosforce.Transfer) {
	seelog.Infof("transfer asset %s : %s --> %s by %s", data.Quantity, data.From, data.To, data.Memo)
	a := &model.Transfer{
		FromAccount: string(data.From),
		ToAccount:   string(data.To),
		RefBlockNum: int64(trx.RefBlockNum),
		BlockID:     blockID,
		Quantity:    data.Quantity.Amount,
		TrxID:       trx.ID(),
		Token:       string(data.Quantity.Symbol.Symbol),
		Fee:         int64(trx.Fee.Amount),
		Memo:        data.Memo,
		Data:        *act,
	}

	err := s.pgDB.Insert(a)

	if err != nil {
		seelog.Errorf("insert transfer err by %s", err.Error())
		return
	}

}

// TODO
func (s *Sync2pgDB) OnVote(blockID string, trx *eos.SignedTransaction) {

}

func NewSyncPgDB(pgAddr string, userName string, passwd string, db string) *Sync2pgDB {
	res := &Sync2pgDB{}
	res.pgOpt = pg.Options{
		Addr:     pgAddr,
		User:     userName,
		Password: passwd,
		Database: db,
	}
	seelog.Infof("connect db %v", res.pgOpt)
	res.pgDB = pg.Connect(&res.pgOpt)
	return res
}
