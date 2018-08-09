package model

type Transfer struct {
	FromAccount string `json:"from_account"`
	ToAccount   string `json:"to_account"`
	Quantity    int64  `json:"quantity"`
	Token       string `json:"token"`
	Memo        string `json:"memo"`
	RefBlockNum int64  `json:"ref_block_num"`
}
