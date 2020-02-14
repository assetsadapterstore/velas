package txsigner

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/assetsadapterstore/velas-adapter/crypto"
	owcrypt "github.com/blocktree/go-owcrypt"
)

var Default = &TransactionSigner{}

type TransactionSigner struct {
}

type SigPub struct {
	Signature []byte
	Pubkey    []byte
}

// SignTransactionHash 交易哈希签名算法
// required
func (singer *TransactionSigner) SignTransactionHash(message []byte, privateKey []byte, eccType uint32) ([]byte, error) {

	if len(message) == 0 {
		return nil, fmt.Errorf("No message to sign")
	}

	if privateKey == nil || len(privateKey) != 32 {
		return nil, fmt.Errorf("Invalid private key")
	}

	signature, _, retCode := owcrypt.Signature(privateKey, nil, message, eccType)

	if retCode != owcrypt.SUCCESS {
		return nil, fmt.Errorf("Failed to sign message")
	}

	return signature, nil
}

// VerifyAndCombineTransaction verify signature
// required
func (singer *TransactionSigner) VerifyAndCombineTransaction(emptyTrans string, sigPub []SigPub) (bool, string, error) {
	trx := crypto.Tx{}

	err := json.Unmarshal([]byte(emptyTrans), &trx)

	if err != nil {
		return false, "", errors.New("Invalid empty transaction data")
	}

	if sigPub == nil || len(sigPub) == 0 || len(sigPub) != len(trx.Inputs) {
		return false, "", errors.New("Signatures are not enough to unlock transaction")
	}

	for i, s := range sigPub {
		utxo := trx.Inputs[i]
		msg := trx.MsgForSign(utxo.PreviousOutput.Hash, utxo.PreviousOutput.Index)

		fmt.Println("msg:", hex.EncodeToString(msg))
		fmt.Println("sig:", hex.EncodeToString(s.Signature))
		fmt.Println("pub:", hex.EncodeToString(s.Pubkey))
		if owcrypt.SUCCESS != owcrypt.Verify(s.Pubkey, nil, msg, s.Signature, owcrypt.ECC_CURVE_ED25519) {
			return false, "", errors.New("Signature verify failed")
		}
		utxo.Script = s.Signature
		utxo.PublicKey = s.Pubkey
		trx.Inputs[i] = utxo
	}

	txHash := trx.GenerateHash()
	trx.Hash = txHash

	txBytes, err := trx.MarshalJSON()
	if err != nil {
		return false, "", errors.New("Failed to marshal transaction")
	}

	return true, hex.EncodeToString(txBytes), nil
}
