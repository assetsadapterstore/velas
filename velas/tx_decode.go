/*
 * Copyright 2018 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package velas

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/assetsadapterstore/velas-adapter/crypto"
	"github.com/assetsadapterstore/velas-adapter/txsigner"
	"github.com/blocktree/openwallet/common"
	"github.com/blocktree/openwallet/openwallet"
	"github.com/shopspring/decimal"
)

type TransactionDecoder struct {
	openwallet.TransactionDecoderBase
	wm *WalletManager //钱包管理者
}

//NewTransactionDecoder 交易单解析器
func NewTransactionDecoder(wm *WalletManager) *TransactionDecoder {
	decoder := TransactionDecoder{}
	decoder.wm = wm
	return &decoder
}

//CreateRawTransaction 创建交易单
func (decoder *TransactionDecoder) CreateRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {
	if rawTx.Coin.IsContract {
		return fmt.Errorf("do not support token transaction")
	} else {
		return decoder.CreateVLXRawTransaction(wrapper, rawTx)
	}
}

//SignRawTransaction 签名交易单
func (decoder *TransactionDecoder) SignRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {
	if rawTx.Coin.IsContract {
		return fmt.Errorf("do not support token transaction")
	} else {
		return decoder.SignVLXRawTransaction(wrapper, rawTx)
	}
}

//VerifyRawTransaction 验证交易单，验证交易单并返回加入签名后的交易单
func (decoder *TransactionDecoder) VerifyRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {
	if rawTx.Coin.IsContract {
		return fmt.Errorf("do not support token transaction")
	} else {
		return decoder.VerifyVLXRawTransaction(wrapper, rawTx)
	}
}

// CreateSummaryRawTransactionWithError 创建汇总交易，返回能原始交易单数组（包含带错误的原始交易单）
func (decoder *TransactionDecoder) CreateSummaryRawTransactionWithError(wrapper openwallet.WalletDAI, sumRawTx *openwallet.SummaryRawTransaction) ([]*openwallet.RawTransactionWithError, error) {
	if sumRawTx.Coin.IsContract {
		return nil, fmt.Errorf("do not support token transaction")
	} else {
		return decoder.CreateVLXSummaryRawTransaction(wrapper, sumRawTx)
	}
}

//CreateSummaryRawTransaction 创建汇总交易，返回原始交易单数组
func (decoder *TransactionDecoder) CreateSummaryRawTransaction(wrapper openwallet.WalletDAI, sumRawTx *openwallet.SummaryRawTransaction) ([]*openwallet.RawTransaction, error) {
	var (
		rawTxWithErrArray []*openwallet.RawTransactionWithError
		rawTxArray        = make([]*openwallet.RawTransaction, 0)
		err               error
	)
	if sumRawTx.Coin.IsContract {
		return nil, fmt.Errorf("do not support token transaction")
	} else {
		rawTxWithErrArray, err = decoder.CreateVLXSummaryRawTransaction(wrapper, sumRawTx)
	}
	if err != nil {
		return nil, err
	}
	for _, rawTxWithErr := range rawTxWithErrArray {
		if rawTxWithErr.Error != nil {
			continue
		}
		rawTxArray = append(rawTxArray, rawTxWithErr.RawTx)
	}
	return rawTxArray, nil
}

//SubmitRawTransaction 广播交易单
func (decoder *TransactionDecoder) SubmitRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) (*openwallet.Transaction, error) {
	var trx crypto.Tx

	if len(rawTx.RawHex) == 0 {
		return nil, fmt.Errorf("transaction hex is empty")
	}

	if !rawTx.IsCompleted {
		return nil, fmt.Errorf("transaction is not completed validation")
	}

	rawHex, err := hex.DecodeString(rawTx.RawHex)
	if err != nil {
		return nil, openwallet.ConvertError(err)
	}

	err = json.Unmarshal(rawHex, &trx)
	if err != nil {
		return nil, openwallet.ConvertError(err)
	}

	err = decoder.wm.WalletClient.Tx.Validate(trx)
	if err != nil {
		return nil, err
	}

	_, err = decoder.wm.WalletClient.Tx.Publish(trx)
	if err != nil {
		return nil, err
	}

	rawTx.TxID = hex.EncodeToString(trx.Hash[:])
	rawTx.IsSubmit = true

	decimals := int32(0)
	fees := "0"
	if rawTx.Coin.IsContract {
		decimals = int32(rawTx.Coin.Contract.Decimals)
		fees = "0"
	} else {
		decimals = int32(decoder.wm.Decimal())
		fees = rawTx.Fees
	}

	//记录一个交易单
	tx := &openwallet.Transaction{
		From:       rawTx.TxFrom,
		To:         rawTx.TxTo,
		Amount:     rawTx.TxAmount,
		Coin:       rawTx.Coin,
		TxID:       rawTx.TxID,
		Decimal:    decimals,
		AccountID:  rawTx.Account.AccountID,
		Fees:       fees,
		SubmitTime: time.Now().Unix(),
	}

	tx.WxID = openwallet.GenTransactionWxID(tx)

	return tx, nil
}

////////////////////////// VLX implement //////////////////////////

//CreateVLXRawTransaction 创建交易单
func (decoder *TransactionDecoder) CreateVLXRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {

	var (
		unspents    []*crypto.TransactionInputOutpoint
		affordUTXO  []*crypto.TransactionInputOutpoint
		outputAddrs = make(map[string]decimal.Decimal)
		balance     = decimal.New(0, 0)
		totalSend   = decimal.New(0, 0)
		fixFees     = decimal.New(0, 0)
		accountID   = rawTx.Account.AccountID
		targets     = make([]string, 0)
		limit       = 2000
	)

	if len(rawTx.To) == 0 {
		return errors.New("Receiver address is empty")
	}

	address, err := wrapper.GetAddressList(0, limit, "AccountID", rawTx.Account.AccountID)
	if err != nil {
		return err
	}

	if len(address) == 0 {
		return openwallet.Errorf(openwallet.ErrAccountNotAddress, "[%s] have not address", accountID)
	}

	searchAddrs := make([]string, 0)
	for _, address := range address {
		searchAddrs = append(searchAddrs, address.Address)
		outputs, err := decoder.wm.WalletClient.Wallet.GetUnspent(address.Address)
		if err != nil {
			return err
		}
		unspents = append(unspents, outputs...)
	}

	if len(unspents) == 0 {
		return fmt.Errorf("[%s] balance is not enough", accountID)
	}

	//计算总发送金额
	for addr, amount := range rawTx.To {
		amt, _ := decimal.NewFromString(amount)
		totalSend = totalSend.Add(amt)
		targets = append(targets, addr)
	}

	//获取utxo，按小到大排序
	sort.Slice(unspents, func(a, b int) bool {
		return unspents[a].Value < unspents[b].Value
	})

	if len(rawTx.FeeRate) == 0 {
		fixFees, err = decimal.NewFromString(decoder.wm.Config.FixFees)
		if err != nil {
			return err
		}
	} else {
		fixFees, _ = decimal.NewFromString(rawTx.FeeRate)
	}

	decoder.wm.Log.Info("Calculating wallet unspent record to build transaction...")
	computeTotalSend := totalSend.Add(fixFees)

	//计算一个可用于支付的余额
	for _, u := range unspents {
		v := common.IntToDecimals(int64(u.Value), decoder.wm.Decimal())
		balance = balance.Add(v)
		affordUTXO = append(affordUTXO, u)
		if balance.GreaterThanOrEqual(computeTotalSend) {
			break
		}
	}

	//判断余额是否足够支付发送数额+手续费
	if balance.LessThan(computeTotalSend) {
		return fmt.Errorf("The balance: %s is not enough! ", balance.StringFixed(decoder.wm.Decimal()))
	}

	//取账户最后一个地址
	changeAddress := affordUTXO[0].Address

	changeAmount := balance.Sub(computeTotalSend)
	rawTx.FeeRate = fixFees.StringFixed(decoder.wm.Decimal())
	rawTx.Fees = fixFees.StringFixed(decoder.wm.Decimal())

	decoder.wm.Log.Std.Notice("-----------------------------------------------")
	decoder.wm.Log.Std.Notice("From Account: %s", accountID)
	decoder.wm.Log.Std.Notice("To Address: %s", strings.Join(targets, ", "))
	decoder.wm.Log.Std.Notice("Balance: %v", balance.String())
	decoder.wm.Log.Std.Notice("Fees: %v", fixFees.String())
	decoder.wm.Log.Std.Notice("Receive: %v", computeTotalSend.String())
	decoder.wm.Log.Std.Notice("Change: %v", changeAmount.String())
	decoder.wm.Log.Std.Notice("Change Address: %v", changeAddress)
	decoder.wm.Log.Std.Notice("-----------------------------------------------")

	//装配输出
	for to, amount := range rawTx.To {
		decamount, _ := decimal.NewFromString(amount)
		outputAddrs = appendOutput(outputAddrs, to, decamount)
	}

	if changeAmount.GreaterThan(decimal.New(0, 0)) {
		outputAddrs = appendOutput(outputAddrs, changeAddress, changeAmount)
	}

	err = decoder.createVLXRawTransaction(wrapper, rawTx, affordUTXO, outputAddrs, changeAddress, fixFees)
	if err != nil {
		return err
	}

	return nil
}

//SignVLXRawTransaction 签名交易单
func (decoder *TransactionDecoder) SignVLXRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {

	if rawTx.Signatures == nil || len(rawTx.Signatures) == 0 {
		return fmt.Errorf("transaction signature is empty")
	}

	key, err := wrapper.HDKey()
	if err != nil {
		return err
	}

	keySignatures := rawTx.Signatures[rawTx.Account.AccountID]
	if keySignatures != nil {
		for _, keySignature := range keySignatures {

			childKey, err := key.DerivedKeyWithPath(keySignature.Address.HDPath, keySignature.EccType)
			keyBytes, err := childKey.GetPrivateKeyBytes()
			if err != nil {
				return err
			}
			txHash := keySignature.Message
			decoder.wm.Log.Debug("hash:", txHash)

			data, err := hex.DecodeString(txHash)
			if err != nil {
				return fmt.Errorf("Invalid message to sign")
			}

			//签名交易
			/////////交易单哈希签名
			signature, err := txsigner.Default.SignTransactionHash(data, keyBytes, keySignature.EccType)
			if err != nil {
				return fmt.Errorf("transaction hash sign failed, unexpected error: %v", err)
			}

			keySignature.Signature = hex.EncodeToString(signature)
		}
	}

	decoder.wm.Log.Info("transaction hash sign success")

	rawTx.Signatures[rawTx.Account.AccountID] = keySignatures

	return nil
}

//VerifyVLXRawTransaction 验证交易单，验证交易单并返回加入签名后的交易单
func (decoder *TransactionDecoder) VerifyVLXRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {

	var (
		sigPub = make([]txsigner.SigPub, 0)
	)

	rawHex, err := hex.DecodeString(rawTx.RawHex)
	if err != nil {
		return err
	}

	emptyTrans := string(rawHex)

	if rawTx.Signatures == nil || len(rawTx.Signatures) == 0 {
		return fmt.Errorf("transaction signature is empty")
	}

	for accountID, keySignatures := range rawTx.Signatures {
		decoder.wm.Log.Debug("accountID Signatures:", accountID)
		for _, keySignature := range keySignatures {

			signature, _ := hex.DecodeString(keySignature.Signature)
			pubkey, _ := hex.DecodeString(keySignature.Address.PublicKey)

			signaturePubkey := txsigner.SigPub{
				Signature: signature,
				Pubkey:    pubkey,
			}

			sigPub = append(sigPub, signaturePubkey)

			decoder.wm.Log.Debug("Signature:", keySignature.Signature)
			decoder.wm.Log.Debug("PublicKey:", keySignature.Address.PublicKey)
		}
	}

	/////////验证交易单
	//验证时，对于公钥哈希地址，需要将对应的锁定脚本传入TxUnlock结构体
	pass, signedTrans, err := txsigner.Default.VerifyAndCombineTransaction(emptyTrans, sigPub)
	if pass {
		decoder.wm.Log.Debug("transaction verify passed")
		rawTx.IsCompleted = true
		rawTx.RawHex = signedTrans
	} else {
		decoder.wm.Log.Errorf("transaction verify failed, unexpected error: %v", err)
		rawTx.IsCompleted = false
	}

	return nil
}

//GetRawTransactionFeeRate 获取交易单的费率
func (decoder *TransactionDecoder) GetRawTransactionFeeRate() (feeRate string, unit string, err error) {
	return decoder.wm.Config.FixFees, "TX", nil
}

//CreateVLXSummaryRawTransaction 创建BTC汇总交易
func (decoder *TransactionDecoder) CreateVLXSummaryRawTransaction(wrapper openwallet.WalletDAI, sumRawTx *openwallet.SummaryRawTransaction) ([]*openwallet.RawTransactionWithError, error) {

	var (
		accountID          = sumRawTx.Account.AccountID
		minTransfer, _     = decimal.NewFromString(sumRawTx.MinTransfer)
		retainedBalance, _ = decimal.NewFromString(sumRawTx.RetainedBalance)
		sumAddresses       = make([]string, 0)
		rawTxArray         = make([]*openwallet.RawTransactionWithError, 0)
		outputAddrs        map[string]decimal.Decimal
		totalInputAmount   decimal.Decimal
		unspents           []*crypto.TransactionInputOutpoint
		sumUnspents        []*crypto.TransactionInputOutpoint
		fixFees            = decimal.New(0, 0)
	)

	if minTransfer.LessThan(retainedBalance) {
		return nil, fmt.Errorf("mini transfer amount must be greater than address retained balance")
	}

	address, err := wrapper.GetAddressList(sumRawTx.AddressStartIndex, sumRawTx.AddressLimit, "AccountID", sumRawTx.Account.AccountID)
	if err != nil {
		return nil, err
	}

	if len(address) == 0 {
		return nil, fmt.Errorf("[%s] have not addresses", accountID)
	}

	searchAddrs := make([]string, 0)
	for _, address := range address {
		searchAddrs = append(searchAddrs, address.Address)
	}

	addrBalanceArray, err := decoder.wm.Blockscanner.GetBalanceByAddress(searchAddrs...)
	if err != nil {
		return nil, err
	}

	for _, addrBalance := range addrBalanceArray {
		decoder.wm.Log.Debugf("addrBalance: %+v", addrBalance)
		//检查余额是否超过最低转账
		addrBalanceDec, _ := decimal.NewFromString(addrBalance.Balance)
		if addrBalanceDec.GreaterThanOrEqual(minTransfer) {
			//添加到转账地址数组
			sumAddresses = append(sumAddresses, addrBalance.Address)
		}
	}

	if len(sumAddresses) == 0 {
		return nil, nil
	}

	//取得费率
	if len(sumRawTx.FeeRate) == 0 {
		fixFees, err = decimal.NewFromString(decoder.wm.Config.FixFees)
		if err != nil {
			return nil, err
		}
	} else {
		fixFees, _ = decimal.NewFromString(sumRawTx.FeeRate)
	}

	unspents = make([]*crypto.TransactionInputOutpoint, 0)
	outputAddrs = make(map[string]decimal.Decimal, 0)
	totalInputAmount = decimal.Zero

	for i, addr := range sumAddresses {

		outputs := make([]*crypto.TransactionInputOutpoint, 0)
		outputs, err := decoder.wm.WalletClient.Wallet.GetUnspent(addr)
		if err != nil {
			return nil, err
		}

		if len(outputs) > 0 {
			unspents = append(unspents, outputs...)
			sumUnspents = append(sumUnspents, unspents...)
			if retainedBalance.GreaterThan(decimal.Zero) {
				outputAddrs = appendOutput(outputAddrs, addr, retainedBalance)
			}
		}

		//如果遍历地址完结，就可以进行构建交易单
		if i == len(sumAddresses)-1 {

			//计算这笔交易单的汇总数量
			for _, u := range sumUnspents {
				ua := common.IntToDecimals(int64(u.Value), decoder.wm.Decimal())
				totalInputAmount = totalInputAmount.Add(ua)
			}

			/*

				汇总数量计算：

				1. 输入总数量 = 合计账户地址的所有utxo
				2. 账户地址输出总数量 = 账户地址保留余额 * 地址数
				3. 汇总数量 = 输入总数量 - 账户地址输出总数量 - 手续费
			*/
			retainedBalanceTotal := retainedBalance.Mul(decimal.New(int64(len(outputAddrs)), 0))
			sumAmount := totalInputAmount.Sub(retainedBalanceTotal).Sub(fixFees)

			decoder.wm.Log.Debugf("totalInputAmount: %v", totalInputAmount)
			decoder.wm.Log.Debugf("retainedBalanceTotal: %v", retainedBalanceTotal)
			decoder.wm.Log.Debugf("fees: %v", fixFees)
			decoder.wm.Log.Debugf("sumAmount: %v", sumAmount)

			//最后填充汇总地址及汇总数量
			outputAddrs = appendOutput(outputAddrs, sumRawTx.SummaryAddress, sumAmount)

			raxTxTo := make(map[string]string, 0)
			for a, m := range outputAddrs {
				raxTxTo[a] = m.StringFixed(decoder.wm.Decimal())
			}

			//创建一笔交易单
			rawTx := &openwallet.RawTransaction{
				Coin:     sumRawTx.Coin,
				Account:  sumRawTx.Account,
				FeeRate:  sumRawTx.FeeRate,
				To:       raxTxTo,
				Fees:     fixFees.StringFixed(decoder.wm.Decimal()),
				Required: 1,
			}

			createErr := decoder.createVLXRawTransaction(wrapper, rawTx, sumUnspents, outputAddrs, "", fixFees)
			rawTxWithErr := &openwallet.RawTransactionWithError{
				RawTx: rawTx,
				Error: openwallet.ConvertError(createErr),
			}

			//创建成功，添加到队列
			rawTxArray = append(rawTxArray, rawTxWithErr)

			//清空临时变量
			unspents = make([]*crypto.TransactionInputOutpoint, 0)
			outputAddrs = make(map[string]decimal.Decimal, 0)
			totalInputAmount = decimal.Zero

		}
	}

	return rawTxArray, nil
}

//createVLXRawTransaction 创建VLX原始交易单
func (decoder *TransactionDecoder) createVLXRawTransaction(
	wrapper openwallet.WalletDAI,
	rawTx *openwallet.RawTransaction,
	affordUTXO []*crypto.TransactionInputOutpoint,
	to map[string]decimal.Decimal,
	changeAddress string,
	fees decimal.Decimal,
) error {

	var (
		err              error
		totalSend        = decimal.New(0, 0)
		targets          = make([]string, 0)
		accountTotalSent = decimal.Zero
		txFrom           = make([]string, 0)
		txTo             = make([]string, 0)
		accountID        = rawTx.Account.AccountID
		vouts            = make(map[string]uint64)
	)

	if len(affordUTXO) == 0 {
		return fmt.Errorf("utxo is empty")
	}

	if len(to) == 0 {
		return fmt.Errorf("Receiver addresses is empty! ")
	}

	//计算总发送金额
	for addr, amount := range to {
		//deamount, _ := decimal.NewFromString(amount)
		totalSend = totalSend.Add(amount)
		targets = append(targets, addr)
		//计算账户的实际转账amount
		addresses, findErr := wrapper.GetAddressList(0, -1, "AccountID", accountID, "Address", addr)
		if findErr != nil || len(addresses) == 0 {
			//amountDec, _ := decimal.NewFromString(amount)
			accountTotalSent = accountTotalSent.Add(amount)
		}
	}

	//装配输入
	for _, utxo := range affordUTXO {
		amount := common.IntToDecimals(int64(utxo.Value), decoder.wm.Decimal())
		txFrom = append(txFrom, fmt.Sprintf("%s:%s", utxo.Address, amount))
	}

	//装配输入
	for to, amount := range to {
		txTo = append(txTo, fmt.Sprintf("%s:%s", to, amount.String()))
		amount = amount.Shift(decoder.wm.Decimal())
		intAmount := uint64(amount.IntPart())
		if origin, ok := vouts[to]; ok {
			origin += intAmount
			vouts[to] = origin
		} else {
			vouts[to] = intAmount
		}
	}

	commission := uint64(fees.Shift(decoder.wm.Decimal()).IntPart())

	trx, err := crypto.NewTransaction(affordUTXO, vouts, changeAddress, commission)
	if err != nil {
		return openwallet.Errorf(openwallet.ErrCreateRawTransactionFailed, "create transaction failed, unexpected error: %v", err)
	}

	json, err := trx.MarshalJSON()
	if err != nil {
		return openwallet.Errorf(openwallet.ErrCreateRawTransactionFailed, "marshal transaction failed, unexpected error: %v", err)
	}
	rawTx.RawHex = hex.EncodeToString(json)

	if rawTx.Signatures == nil {
		rawTx.Signatures = make(map[string][]*openwallet.KeySignature)
	}

	//装配签名
	keySigs := make([]*openwallet.KeySignature, 0)

	for i, utxo := range affordUTXO {

		sigMsg := trx.MsgForSign(utxo.Hash, utxo.Index)
		beSignHex := hex.EncodeToString(sigMsg)

		decoder.wm.Log.Std.Debug("txHash[%d]: %s", i, beSignHex)

		addr, err := wrapper.GetAddress(utxo.Address)
		if err != nil {
			return err
		}

		signature := openwallet.KeySignature{
			EccType: decoder.wm.Config.CurveType,
			Nonce:   "",
			Address: addr,
			Message: beSignHex,
		}

		keySigs = append(keySigs, &signature)

	}

	feesDec, _ := decimal.NewFromString(rawTx.Fees)
	accountTotalSent = accountTotalSent.Add(feesDec)
	accountTotalSent = decimal.Zero.Sub(accountTotalSent)

	//TODO:多重签名要使用owner的公钥填充

	rawTx.Signatures[rawTx.Account.AccountID] = keySigs
	rawTx.IsBuilt = true
	rawTx.TxAmount = accountTotalSent.StringFixed(decoder.wm.Decimal())
	rawTx.TxFrom = txFrom
	rawTx.TxTo = txTo

	return nil
}

func appendOutput(output map[string]decimal.Decimal, address string, amount decimal.Decimal) map[string]decimal.Decimal {
	if origin, ok := output[address]; ok {
		origin = origin.Add(amount)
		output[address] = origin
	} else {
		output[address] = amount
	}
	return output
}
