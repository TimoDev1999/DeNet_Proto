package utils

import (
	"context"
	"encoding/hex"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func InitEthereumClient(InfuraURL string) *ethclient.Client {
	client, err := ethclient.Dial(InfuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	return client
}

func VerifySignature(address common.Address, signature string) bool {
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		log.Printf("Failed to decode signature: %v", err)
		return false
	}

	hash := crypto.Keccak256Hash([]byte(address.Hex()))

	recoveredPubKey, err := crypto.SigToPub(hash.Bytes(), sigBytes)
	if err != nil {
		log.Printf("Failed to recover public key from signature: %v", err)
		return false
	}

	recoveredAddress := crypto.PubkeyToAddress(*recoveredPubKey)
	return recoveredAddress == address
}

func GetBalanceAndNonce(client *ethclient.Client, address common.Address) (uint64, uint64, error) {
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return 0, 0, err
	}

	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		return 0, 0, err
	}

	return balance.Uint64(), nonce, nil
}

const erc20ABI = `[{"constant":true,"inputs":[{"name":"_owner","type":"address"}],
"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,
"stateMutability":"view","type":"function"}]`

func GetERC20Balance(client *ethclient.Client, address, tokenAddress common.Address) (uint64, error) {
	tokenABI, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		return 0, err
	}

	data, err := tokenABI.Pack("balanceOf", address)
	if err != nil {
		return 0, err
	}

	msg := ethereum.CallMsg{
		To:   &tokenAddress,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return 0, err
	}

	var balance *big.Int
	err = tokenABI.UnpackIntoInterface(&balance, "balanceOf", result)
	if err != nil {
		return 0, err
	}

	return balance.Uint64(), nil
}
