package model

import (
	"github.com/fanyang1988/eos-go/system"
	"time"
)

type Accounts struct {
	Name        string            `json:"name"`
	CreateAt    time.Time         `json:"create_at"`
	UpdateAt    time.Time         `json:"update_at"`
	Creater     string            `json:"creater"`
	RefBlockNum int64             `json:"ref_block_num"`
	Data        system.NewAccount `json:"data"`
}
