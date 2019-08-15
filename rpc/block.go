package rpc

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/assetsadapterstore/velas-adapter/crypto"
	"github.com/go-errors/errors"
	"gopkg.in/resty.v1"
)

type Block struct {
	bk *BaseClient
}

func newBlockClient(bk *BaseClient) *Block {
	return &Block{
		bk: bk,
	}
}

type BlockResponse struct {
	Header       *Header      `json:"header"` // block headero
	Transactions []*crypto.Tx `json:"txns"`   // Block transactions, in format of "tx" command
}

type Header struct {
	Type        uint32 `json:"type"`         // Block type
	Hash        string `json:"hash"`         // Hash
	Height      uint32 `json:"height"`       // Height
	Size        uint64 `json:"size"`         // Size
	Version     uint32 `json:"version"`      // Block version information (note, this is signed)
	PrevBlock   string `json:"prev_block"`   // The hash value of the previous block this particular block references
	MerkleRoot  string `json:"merkle_root"`  // The reference to a Merkle tree collection which is a hash of all transactions related to this block
	Timestamp   uint32 `json:"timestamp"`    // A timestamp recording when this block was created (Will overflow in 2106[2])
	Bits        uint32 `json:"bits"`         // Not used
	Nonce       uint32 `json:"nonce"`        // The nonce used to generate this block… to allow variations of the header and compute different hashes
	Seed        string `json:"seed"`         // The random seed, not used
	TxnCount    uint32 `json:"txn_count"`    // Transaction count
	AdviceCount uint32 `json:"advice_count"` // Advise list count
	Script      string `json:"script"`       // The node's (block owner) signature
}

func (blk *Block) GetByHash(hash string) (*BlockResponse, error) {
	resp, err := resty.
		R().
		Get(blk.bk.baseAddress + "/api/v1/blocks/" + hash)
	if err != nil {
		return nil, err
	}
	body, err := blk.bk.ReadResponse(resp)
	if err != nil {
		return nil, err
	}
	blockResponse := BlockResponse{}
	if err := json.Unmarshal(body, &blockResponse); err != nil {
		return nil, errors.New(err)
	}
	return &blockResponse, nil
}

func (blk *Block) GetByHeight(height uint32) (*BlockResponse, error) {

	resp, err := resty.
		R().
		Get(blk.bk.baseAddress + "/api/v1/blocks?limit=1")
	if err != nil {
		return nil, err
	}
	body, err := blk.bk.ReadResponse(resp)
	if err != nil {
		return nil, err
	}
	headers := []*Header{}
	if err := json.Unmarshal(body, &headers); err != nil {
		return nil, errors.New(err)
	}
	if len(headers) == 0 {
		return nil, fmt.Errorf("cannot get headers")
	}

	offset := strconv.FormatInt(int64((headers[0].Height - height)), 10)
	resp, err = resty.
		R().
		Get(blk.bk.baseAddress + "/api/v1/blocks?limit=20&offset=" + offset)
	if err != nil {
		return nil, err
	}
	body, err = blk.bk.ReadResponse(resp)
	if err != nil {
		return nil, err
	}
	blocks := []*Header{}
	if err := json.Unmarshal(body, &blocks); err != nil {
		return nil, errors.New(err)
	}

	for _, block := range blocks {
		if block.Height == height {
			return blk.GetByHash(block.Hash)
		}
	}

	return nil, fmt.Errorf("cannot fetch block [%v]", height)
}
