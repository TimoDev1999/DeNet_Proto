package main

import (
	"DeNet/proto"
	"DeNet/server/handlers"
	"DeNet/server/utils"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("../config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %v", err)
	}

	InfuraURL := viper.GetString("infura_url")
	serverPort := viper.GetString("server_port")

	client := utils.InitEthereumClient(InfuraURL)

	listener, err := net.Listen("tcp", serverPort)
	if err != nil {
		log.Fatalf("Error listening on port %s", serverPort)
	}

	s := grpc.NewServer()
	proto.RegisterEthereumServiceServer(s, &handlers.Server{Client: client})

	log.Printf("Starting server on port %s", serverPort)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Error serving on port %s", serverPort)
	}

}
