package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum/crypto"
)

type Wallet struct {
	Address   string
	Signature string
}

func CreateAndSignWallet() (*Wallet, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	signature, err := SignData(privateKey, []byte(address))
	if err != nil {
		return nil, err
	}

	return &Wallet{
		Address:   address,
		Signature: hex.EncodeToString(signature),
	}, nil
}

func SignData(privateKey *ecdsa.PrivateKey, data []byte) ([]byte, error) {
	hash := crypto.Keccak256Hash(data)
	return crypto.Sign(hash.Bytes(), privateKey)
}
