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

import "testing"

func TestGetBlockHeight(t *testing.T) {
	height, _ := tw.GetBlockHeight()
	t.Logf("GetBlockHeight height = %d \n", height)
}

func TestGetBlock(t *testing.T) {
	block, err := tw.WalletClient.Block.GetByHeight(46044)
	if err != nil {
		t.Errorf("GetBlockHash failed unexpected error: %v\n", err)
		return
	}
	t.Logf("GetByHeight = %v \n", block)
}

// func TestGetBlock(t *testing.T) {
// 	raw, err := tw.GetBlock("2E643447A46CC033A3D4576858C0FF664A84F2F7BE79B3D63EBE34C18AD4E1C4")
// 	if err != nil {
// 		t.Errorf("GetBlock failed unexpected error: %v\n", err)
// 		return
// 	}
// 	t.Logf("GetBlock = %v \n", raw)
// }

// func TestGetTransaction(t *testing.T) {
// 	raw, err := tw.GetTransaction("67D60DD67995B06574F9634DC5C50C96BE5B479A07AF3894E91001F5A607051C")
// 	if err != nil {
// 		t.Errorf("GetTransaction failed unexpected error: %v\n", err)
// 		return
// 	}

// 	t.Logf("BlockHash = %v \n", raw.BlockHash)
// 	t.Logf("BlockHeight = %v \n", raw.BlockHeight)
// 	t.Logf("Blocktime = %v \n", raw.Blocktime)
// 	t.Logf("Fees = %v \n", raw.Fees)

// 	t.Logf("========= vins ========= \n")

// 	for i, vin := range raw.Vins {
// 		t.Logf("TxID[%d] = %v \n", i, vin.TxID)
// 		t.Logf("Vout[%d] = %v \n", i, vin.Vout)
// 		t.Logf("Addr[%d] = %v \n", i, vin.Addr)
// 		t.Logf("Value[%d] = %v \n", i, vin.Amount)
// 	}

// 	t.Logf("========= vouts ========= \n")

// 	for i, out := range raw.Vouts {
// 		t.Logf("Vout[%d] = %v \n", i, out.Vout)
// 		t.Logf("ScriptPubKey[%d] = %v \n", i, out.LockScript)
// 		t.Logf("Addr[%d] = %v \n", i, out.Addr)
// 		t.Logf("Value[%d] = %v \n", i, out.Amount)
// 	}
// }
