package model

import (
	"time"

	"github.com/fanyang1988/eos-go"
)

type Transaction struct {
	Id             string                `json:"id"`
	RefBlockNum    int64                 `json:"ref_block_num"`
	RefBlockPrefix int64                 `json:"ref_block_prefix"`
	BlockID        string                `json:"block_id"`
	Expiration     time.Time             `json:"expiration"`
	NumActions     int64                 `json:"num_actions"`
	DelaySec       int64                 `json:"delay_sec"`
	Fee            int64                 `json:"fee"`
	Data           eos.SignedTransaction `json:"data"`
}
