package crypto

import (
	"github.com/go-errors/errors"
)

type HD struct {
	publicKey  []byte
	privateKey []byte
}

// func GenerateHD() (*HD, error) {
// 	privateKey, publicKey, err := cryptosign.CryptoSignKeyPair()
// 	if err != 0 {
// 		return nil, errors.Errorf("can't generate keys, libsodium error code %d", err)
// 	}
// 	return &HD{
// 		publicKey:  publicKey,
// 		privateKey: privateKey,
// 	}, nil
// }

// func HDFromPrivateKeyHex(sk string) (*HD, error) {
// 	bytesSK, err := hex.DecodeString(sk)
// 	if err != nil {
// 		return nil, errors.New(err)
// 	}
// 	hd := HDFromPrivateKey(bytesSK)
// 	return &hd, nil
// }

// func HDFromPrivateKey(sk []byte) HD {
// 	ssk := sodium.SignSecretKey{Bytes: sk}
// 	spk := ssk.PublicKey()

// 	return HD{
// 		publicKey:  spk.Bytes,
// 		privateKey: ssk.Bytes,
// 	}
// }

// TODO need implement it
func HDFromSeed(seed string) (*HD, error) {
	return nil, errors.Errorf("Not implemented yet")
	/*bytesSeed, err := hex.DecodeString(seed)
	if err != nil {
		return nil, errors.New(err)
	}
	sk, pk, errCode := cryptosign.CryptoSignSeedKeyPair(bytesSeed)
	if errCode != 0 {
		return nil, errors.Errorf("can't generate keys, libsodium error code %d", errCode)
	}
	return &HD{
		publicKey:  pk,
		privateKey: sk,
	}, nil*/
}

func (hd *HD) ToWallet() (*Wallet, error) {
	return CreateWallet(hd.publicKey)
}
