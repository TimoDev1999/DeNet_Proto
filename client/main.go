package main

import (
	"DeNet/client/config"
	"DeNet/client/test"
	"DeNet/client/wallet"
	"DeNet/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {

	config.LoadConfig()
	serverAddress := viper.GetString("server_port")

	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewEthereumServiceClient(conn)

	wallet, err := wallet.CreateAndSignWallet()
	if err != nil {
		log.Fatalf("failed to create and sign wallet: %v", err)
	}

	test.TestGetAccount(client, wallet)

	test.TestGetAccounts(client, 100)
	test.TestGetAccounts(client, 1000)
	test.TestGetAccounts(client, 10000)
}
