package rpc

import (
	"encoding/hex"
	"encoding/json"
	"strconv"

	"github.com/assetsadapterstore/velas-adapter/crypto"
	"github.com/go-errors/errors"
	"gopkg.in/resty.v1"
)

type TxResponse struct {
	Size               uint32 `json:"size"`
	Block              string `json:"block"` // block hash
	Confirmed          uint32 `json:"confirmed"`
	ConfirmedTimestamp uint32 `json:"confirmed_timestamp"`
	Total              int    `json:"total,omitempty"`
	*crypto.Tx
}

func (txr *TxResponse) UnmarshalJSON(data []byte) error {
	type Alias TxResponse
	aux := &struct {
		Hash string `json:"hash"`
		*Alias
	}{
		Alias: (*Alias)(txr),
	}
	if aux.Tx == nil {
		aux.Alias.Tx = &crypto.Tx{}
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	pHash, err := hex.DecodeString(aux.Hash)
	if err != nil {
		return err
	}

	var hash [32]byte
	if len(pHash) == 32 {
		copy(hash[:], pHash[:32])
	}
	txr.Hash = hash
	return nil
}

type Tx struct {
	bk *BaseClient
}

func newTxClient(bk *BaseClient) *Tx {
	return &Tx{
		bk: bk,
	}
}

func (tx *Tx) GetHashListByAddress(address string) ([]string, error) {
	resp, err := resty.
		R().
		Get(tx.bk.baseAddress + "/api/v1/wallet/txs/" + address)
	if err != nil {
		return nil, errors.New(err)
	}
	body, err := tx.bk.ReadResponse(resp)
	response := make([]string, 0)
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.New(err)
	}
	return response, nil
}

func (tx *Tx) GetHashListByHeight(height int) ([]string, error) {
	resp, err := resty.
		R().
		Get(tx.bk.baseAddress + "/api/v1/txs/height/" + strconv.Itoa(height))
	if err != nil {
		return nil, errors.New(err)
	}
	body, err := tx.bk.ReadResponse(resp)
	if err != nil {
		return nil, err
	}
	response := make([]string, 0)
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.New(err)
	}
	return response, nil
}

func (tx *Tx) GetByHashList(hashes ...string) ([]TxResponse, error) {
	arg := struct {
		Hashes []string `json:"hashes"`
	}{Hashes: hashes}
	resp, err := resty.
		R().
		SetBody(arg).
		Post(tx.bk.baseAddress + "/api/v1/txs")
	if err != nil {
		return nil, errors.New(err)
	}
	body, err := tx.bk.ReadResponse(resp)
	if err != nil {
		return nil, err
	}
	response := make([]TxResponse, 0)
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.New(err)
	}
	return response, nil
}

type TxPublishResponse struct {
	Result string `json:"result"`
}

func (tx *Tx) Validate(txData crypto.Tx) error {
	resp, err := resty.
		R().
		SetBody(&txData).
		Post(tx.bk.baseAddress + "/api/v1/txs/validate")
	if err != nil {
		return errors.New(err)
	}
	body, err := tx.bk.ReadResponse(resp)
	if err != nil {
		return err
	}
	response := TxPublishResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return errors.New(err)
	}
	return nil
}

func (tx *Tx) Publish(txData crypto.Tx) error {
	resp, err := resty.
		R().
		SetBody(&txData).
		Post(tx.bk.baseAddress + "/api/v1/txs/publish")
	if err != nil {
		return errors.New(err)
	}
	body, err := tx.bk.ReadResponse(resp)
	if err != nil {
		return err
	}
	response := TxPublishResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return errors.New(err)
	}
	return nil
}
