package client

import (
	"context"

	"github.com/tomekzakrzewski/lottery/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Endpoint string
	client   types.GeneratorClient
}

func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := types.NewGeneratorClient(conn)

	return &GRPCClient{
		Endpoint: endpoint,
		client:   c,
	}, nil
}

func (c *GRPCClient) GenerateWinningNumbers() *types.WinningNums {
	winningNumbers, err := c.client.GenerateWinningNumbers(context.Background(), &types.Empty{})
	if err != nil {
		return nil
	}
	return winningNumbers
}
