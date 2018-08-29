package p2p_node

import (
	"time"

	"github.com/fanyang1988/eos-go/eosforce/token"

	"github.com/cihub/seelog"
	"github.com/fanyang1988/eos-go"
	"github.com/fanyang1988/eos-go/eosforce"
	"github.com/fanyang1988/eos-go/p2p"
	"github.com/fanyang1988/eos-go/system"
)

type p2pSyncClient struct {
	client    *p2p.Client
	handler   DataHandlerInterface
	stopChann chan error
}

// NewP2PSyncClient create a sync client to eosforce
func NewP2PSyncClient(apiUrl, p2pAddr string) (*p2pSyncClient, error) {
	res := &p2pSyncClient{
		stopChann: make(chan error),
	}

	// get chain info from url
	api := eos.New(apiUrl)
	info, err := api.GetInfo()
	if err != nil {
		return nil, err
	}

	cID := info.ChainID
	client := p2p.NewClient(p2pAddr, cID, 0)
	if err != nil {
		return nil, err
	}

	client.WithLogger(seelog.Current)
	res.client = client

	return res, nil
}

func (p *p2pSyncClient) WithHandler(h DataHandlerInterface) {
	p.handler = h
}

func (p *p2pSyncClient) onMsg(msg p2p.Message) {
	switch msg.Envelope.Type {
	case eos.SignedBlockType:
		{
			signedBlockMsg, ok := msg.Envelope.P2PMessage.(*eos.SignedBlock)
			if !ok {
				seelog.Errorf("typ error by signedBlockMsg")
				return
			}

			blockID, err := signedBlockMsg.BlockID()
			if err != nil {
				seelog.Errorf("block id get err by %s", err.Error())
				return
			}
			blockIDStr := blockID.String()

			p.handler.OnBlock(blockIDStr, signedBlockMsg)

			for _, tr := range signedBlockMsg.Transactions {
				trx, err := tr.Transaction.Packed.Unpack()
				if err != nil {
					seelog.Errorf("transaction unpack err by %s", err.Error())
					continue
				}

				p.handler.OnTrx(blockIDStr, signedBlockMsg, trx)

				for _, action := range trx.Actions {
					p.handler.OnAction(blockIDStr, trx, action)

					switch action.Account {
					case "eosio":
						{
							switch action.Name {
							case "transfer":
								{
									transferAct, ok := action.ActionData.Data.(*eosforce.Transfer)
									if !ok {
										seelog.Errorf("transfer act data err")
										continue
									}

									p.handler.OnTransfer(blockIDStr, trx, action, transferAct)
								}
							case "newaccount":
								{
									newAccountAct, ok := action.ActionData.Data.(*system.NewAccount)
									if !ok {
										seelog.Errorf("newAccountAct act data err")
										continue
									}

									p.handler.OnNewAccount(blockIDStr, trx, action, newAccountAct)
								}
							case "claim":
								{
									act, ok := action.ActionData.Data.(*eosforce.Claim)
									if !ok {
										seelog.Errorf("Claim act data err")
										continue
									}

									p.handler.OnClaim(blockIDStr, trx, action, act)
								}
							case "vote":
								{
									act, ok := action.ActionData.Data.(*eosforce.Vote)
									if !ok {
										seelog.Errorf("Claim act data err")
										continue
									}

									p.handler.OnVote(blockIDStr, trx, action, act)
								}
							}
						}
					case "eosio.token":
						{
							switch action.Name {
							case "transfer":
								{
									transferAct, ok := action.ActionData.Data.(*token.Transfer)
									if !ok {
										seelog.Errorf("token transfer act data err")
										continue
									}

									p.handler.OnTokenTransfer(blockIDStr, trx, action, transferAct)
								}
							case "issue":
								{
									act, ok := action.ActionData.Data.(*token.Issue)
									if !ok {
										seelog.Errorf("token.Issue act data err")
										continue
									}

									p.handler.OnTokenIssue(blockIDStr, trx, action, act)
								}
							case "create":
								{
									act, ok := action.ActionData.Data.(*token.Create)
									if !ok {
										seelog.Errorf("token.Create act data err")
										continue
									}

									p.handler.OnTokenCreate(blockIDStr, trx, action, act)
								}
							}
						}
					}
				}
			}

			return
		}
	}

}

// StartListen just listen new blocks
func (p *p2pSyncClient) StartListen() {
	p.client.RegisterHandler(p2p.HandlerFunc(p.onMsg))

	go func() {
		p.stopChann <- p.client.ConnectRecent()
	}()

	return
}

// StartSync from a block to listen
func (p *p2pSyncClient) StartSync(headBlock uint32, headBlockID eos.SHA256Bytes, headBlockTime time.Time, lib uint32, libID eos.SHA256Bytes) {
	p.client.RegisterHandler(p2p.HandlerFunc(p.onMsg))

	go func() {
		p.stopChann <- p.client.ConnectAndSync(headBlock, headBlockID, headBlockTime, lib, libID)
	}()

	return
}

func (p *p2pSyncClient) StopChann() <-chan error {
	return p.stopChann
}
