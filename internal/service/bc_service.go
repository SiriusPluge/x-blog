package service

import "github.com/ethereum/go-ethereum/ethclient"

type EthService struct {
	client *ethclient.Client
}

func NewEthService(client *ethclient.Client) EthApi {
	return &EthService{
		client: client,
	}
}

