package model

import (
	"time"
)

type Block struct {
	BlockNum    int64     `json:"block_num"`
	ProduceTime time.Time `json:"produce_time"`
	BlockID     string    `json:"block_id"`
	PrevBlockID string    `json:"prev_block_id"`
	Producer    string    `json:"producer"`
}
