package pgsync

import (
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

func (s *Sync2pgDB) OnBlock(block *eos.SignedBlock) {
	seelog.Tracef("on block %v", block.BlockNumber())

	blockID, err := block.BlockID()
	if err != nil {
		seelog.Errorf("BlockID error by signedBlockMsg %v", err.Error())
		return
	}

	b := &model.Block{
		BlockID:               blockID.String(),
		BlockNum:              int64(block.BlockNumber()),
		ProduceTime:           block.Timestamp.Time,
		Producer:              string(block.Producer),
		PrevBlockID:           block.Previous.String(),
		TransactionMerkleRoot: block.TransactionMRoot.String(),
		ActionMerkleRoot:      block.ActionMRoot.String(),
		NumTransactions:       int(len(block.Transactions)),
		Confirmed:             int(block.Confirmed),
	}

	_, err = s.pgDB.Model(b).OnConflict("(block_num) DO UPDATE").
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

func (s *Sync2pgDB) OnTrx(trx *eos.SignedTransaction) {

}

func (s *Sync2pgDB) OnAction(trx *eos.SignedTransaction, act *eos.Action) {
	//seelog.Infof("onAction")

}

func (s *Sync2pgDB) OnNewAccount(trx *eos.SignedTransaction, data *system.NewAccount) {
	seelog.Infof("on new account %v", *data)
	_, err := s.pgDB.Model(&model.Accounts{
		Name:        string(data.Name),
		CreateAt:    trx.Expiration.Time,
		UpdateAt:    trx.Expiration.Time,
		Creater:     string(data.Creator),
		RefBlockNum: int64(trx.RefBlockNum),
		Data:        *data,
	}).OnConflict("(name) DO UPDATE").
		Set("create_at = EXCLUDED.create_at,update_at = EXCLUDED.update_at,creater = EXCLUDED.creater,ref_block_num = EXCLUDED.ref_block_num,data = EXCLUDED.data").
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

func (s *Sync2pgDB) OnTransfer(trx *eos.SignedTransaction, data *eosforce.Transfer) {

}

// TODO
func (s *Sync2pgDB) OnVote(trx *eos.SignedTransaction) {

}

func NewSyncPgDB(pgAddr string, userName string, passwd string, db string) *Sync2pgDB {
	res := &Sync2pgDB{}
	res.pgOpt = pg.Options{
		Addr:     pgAddr,
		User:     userName,
		Password: passwd,
		Database: db,
	}
	res.pgDB = pg.Connect(&res.pgOpt)
	return res
}
