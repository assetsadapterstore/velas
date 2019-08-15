package rpc

const Url = "https://testnet.velas.website"

const Pk = "89d5bd2d31889df63cb1c895e4c6f16772e7b06a8c71228bb59d4c9a0c434fc1f6e586d5d051065a580969d15f48f88251ed24b9c77422410bc39a0e7247e53a"

// func TestGetWalletBalance(t *testing.T) {
// 	client := NewClient(Url)
// 	hd, err := crypto.HDFromPrivateKeyHex(Pk)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	wallet, err := hd.ToWallet()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	balance, err := client.Wallet.GetBalance(wallet.Base58Address)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	t.Log(balance)
// }

// func TestGetWalletUnspents(t *testing.T) {
// 	client := NewClient(Url)
// 	hd, err := crypto.HDFromPrivateKeyHex(Pk)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	wallet, err := hd.ToWallet()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	unspents, err := client.Wallet.GetUnspent(wallet.Base58Address)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	t.Log(unspents)
// }
