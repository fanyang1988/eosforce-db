package model

import "github.com/fanyang1988/eos-go"

type Action struct {
	Account        string     `json:"account"`
	Name           string     `json:"name"`
	TrxID          string     `json:"trx_id"`
	BlockID        string     `json:"block_id"`
	RefBlockNum    int64      `json:"ref_block_num"`
	RefBlockPrefix int64      `json:"ref_block_prefix"`
	Fee            int64      `json:"fee"`
	Data           eos.Action `json:"data"`
}
