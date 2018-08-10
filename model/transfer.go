package model

import "github.com/fanyang1988/eos-go"

type Transfer struct {
	FromAccount string     `json:"from_account"`
	ToAccount   string     `json:"to_account"`
	Quantity    int64      `json:"quantity"`
	Token       string     `json:"token"`
	Memo        string     `json:"memo"`
	RefBlockNum int64      `json:"ref_block_num"`
	Data        eos.Action `json:"data"`
	TrxID       string     `json:"trx_id"`
	BlockID     string     `json:"block_id"`
	Fee         int64      `json:"fee"`
}
