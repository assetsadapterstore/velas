package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/assetsadapterstore/velas-adapter/crypto/helpers"
	"github.com/btcsuite/btcutil/base58"
	"github.com/go-errors/errors"
)

type Tx struct {
	Hash     [32]byte            `json:"hash"`
	Version  uint32              `json:"version"`
	LockTime uint32              `json:"lock_time"`
	Inputs   []TransactionInput  `json:"tx_in"`
	Outputs  []TransactionOutput `json:"tx_out"`
}

func NewTransaction(unspents []*TransactionInputOutpoint, toAddresses map[string]uint64, changeAddress string, commission uint64) (*Tx, error) {
	totalin := uint64(0)
	totalout := uint64(0)

	for _, previousOutput := range unspents {
		totalin += previousOutput.Value
	}

	index := uint32(0)

	txIns := make([]TransactionInput, 0)

	txOuts := make([]TransactionOutput, 0)

	txOuts = append(txOuts, TransactionOutput{
		Index: index,
		Value: commission,
	})

	for to, amount := range toAddresses {
		index++
		totalout += amount
		txOuts = append(txOuts, TransactionOutput{
			Index:         index,
			Script:        base58.Decode(to),
			Value:         amount,
			WalletAddress: base58.Decode(to),
		})
	}

	change := totalin - totalout - commission

	if change < 0 {
		return nil, errors.Errorf("Insufficient funds, total amount %d, commission %d, send amount %d", totalin, commission, totalout)
	} else if change > 0 {
		index++
		txOuts = append(txOuts, TransactionOutput{
			Index:         index,
			Script:        base58.Decode(changeAddress),
			Value:         uint64(change),
			WalletAddress: base58.Decode(changeAddress),
		})
	}

	tx := Tx{
		Version:  1,
		LockTime: 0,
		Outputs:  txOuts,
	}

	for _, previousOutput := range unspents {
		txIns = append(txIns, TransactionInput{
			Sequence:       1,
			PreviousOutput: *previousOutput,
			WalletAddress:  base58.Decode(previousOutput.Address),
		})
	}
	tx.Inputs = txIns
	return &tx, nil
}

// MsgForSign return msg for sign
func (tx *Tx) MsgForSign(hash [32]byte, index uint32) []byte {
	txOutSlices := make([][]byte, 0)
	for _, txOut := range tx.Outputs {
		txOutSlices = append(txOutSlices, txOut.msgForSign())
	}
	txOutSlice := helpers.ConcatByteArray(txOutSlices)

	txSlices := [][]byte{
		hash[:],                            // 32
		helpers.UInt32ToBytes(index),       // 4 bytes
		helpers.UInt32ToBytes(tx.Version),  // 4 bytes
		helpers.UInt32ToBytes(tx.LockTime), // 4 bytes
		txOutSlice,
	}

	return helpers.ConcatByteArray(txSlices)
}

// GenerateHash return generated hash
func (tx *Tx) GenerateHash() [32]byte {
	txInSlices := make([][]byte, 0)
	for _, txIn := range tx.Inputs {
		txInSlices = append(txInSlices, txIn.forBlkHash())
	}
	txInSlice := helpers.ConcatByteArray(txInSlices)

	txOutSlices := make([][]byte, 0)
	for _, txOut := range tx.Outputs {
		txOutSlices = append(txOutSlices, txOut.forBlkHash())
	}
	txOutSlice := helpers.ConcatByteArray(txOutSlices)

	txSlices := [][]byte{
		helpers.UInt32ToBytes(tx.Version),  // 4 bytes
		helpers.UInt32ToBytes(tx.LockTime), // 4 bytes
		txInSlice,
		txOutSlice,
	}

	msg := helpers.ConcatByteArray(txSlices)
	return DHASH(msg)
}

func DHASH(data []byte) [32]byte {
	sum := sha256.Sum256(data)
	sum = sha256.Sum256(sum[:])
	return sum
}

func (tx *Tx) MarshalJSON() ([]byte, error) {
	type Alias Tx
	return json.Marshal(&struct {
		Hash string `json:"hash"`
		*Alias
	}{
		Hash:  hex.EncodeToString(tx.Hash[:]),
		Alias: (*Alias)(tx),
	})
}

// UnmarshalJSON custom json convert
func (tx *Tx) UnmarshalJSON(data []byte) error {
	type Alias Tx
	aux := &struct {
		Hash string `json:"hash"`
		*Alias
	}{
		Alias: (*Alias)(tx),
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
	tx.Hash = hash
	return nil
}
