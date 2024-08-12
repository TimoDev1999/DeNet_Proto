package test

import (
	"DeNet/client/utils"
	"DeNet/client/wallet"
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"DeNet/proto"
)

var tokens = []string{
	"0xdAC17F958D2ee523a2206206994597C13D831ec7", // USDT
	"0xA0b86991c6218b36c1d19D4a2e9eb0CE3606EB48", // USDC
	"0x6B175474E89094C44Da98b954EedeAC495271d0F", // DAI
}

func TestGetAccount(client proto.EthereumServiceClient, wallet *wallet.Wallet) {
	start := time.Now()

	req := &proto.AccountRequest{
		EthereumAddress: wallet.Address,
		CryptoSignature: wallet.Signature,
	}

	resp, err := client.GetAccount(context.Background(), req)
	if err != nil {
		log.Fatalf("could not get account: %v", err)
	}

	fmt.Printf("GetAccount Response: Gastoken Balance: %s, Wallet Nonce: %s\n", resp.GastokenBalance, resp.WalletNonce)

	elapsed := time.Since(start)
	fmt.Printf("GetAccount took %s\n", elapsed)
}

func TestGetAccounts(client proto.EthereumServiceClient, totalRequests int) {
	fmt.Printf("\nRunning test with %d total requests...\n", totalRequests)
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Minute)
	defer cancel()

	stream, err := client.GetAccounts(ctx)
	if err != nil {
		log.Fatalf("could not start stream: %v", err)
	}

	addressChunks := utils.ChunkAddresses(utils.Addresses[:totalRequests], totalRequests/len(tokens))

	var wg sync.WaitGroup
	sem := make(chan struct{}, 5) //

	for _, token := range tokens {
		for _, chunk := range addressChunks {
			wg.Add(1)
			sem <- struct{}{}
			go func(chunk []string, token string) {
				defer wg.Done()
				defer func() { <-sem }()
				req := &proto.AccountsRequest{
					EthereumAddresses: chunk,
					Erc20TokenAddress: token,
				}
				if err := stream.Send(req); err != nil {
					log.Printf("could not send request for token %s: %v", token, err)
					return
				}
			}(chunk, token)
		}
	}

	go func() {
		wg.Wait()
		stream.CloseSend()
	}()

	var resultWg sync.WaitGroup
	resultChan := make(chan *proto.AccountsResponse, totalRequests)

	resultWg.Add(1)
	go func() {
		defer resultWg.Done()
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("could not receive response: %v", err)
				break
			}
			resultChan <- resp
		}
		close(resultChan)
	}()

	go func() {
		for resp := range resultChan {
			fmt.Printf("Address: %s, Balance: %d\n", resp.EthereumAddress, resp.Erc20Balance)
		}
	}()

	resultWg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("Test with %d total requests took %s\n", totalRequests, elapsed)
}
