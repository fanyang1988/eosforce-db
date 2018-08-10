package model

import (
	"time"

	eos "github.com/fanyang1988/eos-go"
)

type Accounts struct {
	Name        string     `json:"name"`
	CreateAt    time.Time  `json:"create_at"`
	UpdateAt    time.Time  `json:"update_at"`
	Creater     string     `json:"creater"`
	RefBlockNum int64      `json:"ref_block_num"`
	Data        eos.Action `json:"data"`
}

type AccountPermission struct {
	Account    string `json:"account"`
	Permission string `json:"permission"`
	Pubkey     string `json:"pubkey"`
}
