package handlers

import (
	"DeNet/proto"
	"DeNet/server/utils"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"io"
	"strconv"
)

type Server struct {
	proto.UnimplementedEthereumServiceServer
	Client *ethclient.Client
}

func (s *Server) GetAccount(ctx context.Context, req *proto.AccountRequest) (*proto.AccountResponse, error) {
	addres := common.HexToAddress(req.EthereumAddress)

	if !utils.VerifySignature(addres, req.CryptoSignature) {
		return nil, fmt.Errorf("invalid signature")
	}

	balance, nonce, err := utils.GetBalanceAndNonce(s.Client, addres)
	if err != nil {
		return nil, err
	}

	return &proto.AccountResponse{
		GastokenBalance: strconv.FormatUint(balance, 10),
		WalletNonce:     strconv.FormatUint(nonce, 10),
	}, nil
}
func (s *Server) GetAccounts(stream proto.EthereumService_GetAccountsServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		for _, ethAddress := range req.EthereumAddresses {
			address := common.HexToAddress(ethAddress)
			tokenAddress := common.HexToAddress(req.Erc20TokenAddress)

			balance, err := utils.GetERC20Balance(s.Client, address, tokenAddress)
			if err != nil {
				return err
			}
			if err := stream.Send(&proto.AccountsResponse{
				EthereumAddress: ethAddress,
				Erc20Balance:    balance,
			}); err != nil {
				return err
			}
		}
	}
}
