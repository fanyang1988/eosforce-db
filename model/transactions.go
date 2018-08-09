package model

import "time"

type Transaction struct {
	RefBlockNum    int64     `json:"ref_block_num"`
	RefblockPrefix int64     `json:"ref_block_prefix"`
	BlockID        string    `json:"block_id"`
	Expiration     int       `json:"expiration"`
	Pending        int       `json:"pending"`
	CreatedAt      time.Time `json:"created_at"`
	NumActions     int64     `json:"num_actions"`
	UpdatedAt      time.Time `json:"updated_at"`
	Irreversible   int64     `json:"irreversible"`
}
