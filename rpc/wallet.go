package rpc

import (
	"encoding/json"

	"github.com/assetsadapterstore/velas-adapter/crypto"
	"github.com/go-errors/errors"
	"gopkg.in/resty.v1"
)

type Wallet struct {
	bk *BaseClient
}

type Balance struct {
	Amount uint64 `json:"amount"`
}

func newWalletClient(bk *BaseClient) *Wallet {
	return &Wallet{
		bk: bk,
	}
}

func (w *Wallet) GetBalance(address string) (uint64, error) {
	resp, err := resty.
		R().
		Get(w.bk.baseAddress + "/api/v1/wallet/balance/" + address)
	if err != nil {
		return 0, err
	}
	body, err := w.bk.ReadResponse(resp)
	if err != nil {
		return 0, err
	}
	balanceResponse := Balance{}
	if err := json.Unmarshal(body, &balanceResponse); err != nil {
		return 0, errors.New(err)
	}
	return balanceResponse.Amount, nil
}

func (w *Wallet) GetUnspent(address string) ([]crypto.TransactionInputOutpoint, error) {
	resp, err := resty.
		R().
		Get(w.bk.baseAddress + "/api/v1/wallet/unspent/" + address)
	if err != nil {
		return nil, err
	}
	body, err := w.bk.ReadResponse(resp)
	if err != nil {
		return nil, err
	}
	unspents := make([]crypto.TransactionInputOutpoint, 0)
	if err := json.Unmarshal(body, &unspents); err != nil {
		return nil, errors.New(err)
	}
	return unspents, nil
}
