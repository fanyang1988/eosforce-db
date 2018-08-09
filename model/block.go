package model

import (
	"time"
)

type Block struct {
	BlockNum              int64     `json:"block_num"`
	ProduceTime           time.Time `json:"produce_time"`
	BlockID               string    `json:"block_id"`
	PrevBlockID           string    `json:"prev_block_id"`
	Producer              string    `json:"producer"`
	TransactionMerkleRoot string    `json:"transaction_merkle_root"`
	ActionMerkleRoot      string    `json:"action_merkle_root"`
	NumTransactions       int       `json:"num_transactions"`
	Confirmed             int       `json:"confirmed"`
}
