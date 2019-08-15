package rpc

import (
	"encoding/json"

	"github.com/assetsadapterstore/velas-adapter/rpc/response"

	"github.com/go-errors/errors"
	"gopkg.in/resty.v1"
)

type Client struct {
	baseAddress string
	Wallet      *Wallet
	Tx          *Tx
	Block       *Block
	bk          *BaseClient
}

func NewClient(baseAddress string) *Client {
	bk := newBaseClient(baseAddress)
	return &Client{
		baseAddress: baseAddress,
		bk:          newBaseClient(baseAddress),
		Wallet:      newWalletClient(bk),
		Tx:          newTxClient(bk),
		Block:       newBlockClient(bk),
	}
}

func (cl *Client) NodeInfo() (*response.Node, error) {
	resp, err := resty.
		R().
		Get(cl.baseAddress + "/api/v1/info")
	if err != nil {
		return nil, errors.New(err)
	}
	body, err := cl.bk.ReadResponse(resp)
	response := response.Node{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.New(err)
	}
	return &response, nil
}
