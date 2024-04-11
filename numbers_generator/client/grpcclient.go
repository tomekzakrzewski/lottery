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

func (c *GRPCClient) GenerateWinningNumbers() *types.WinningNumbers {
	winningNumbers, err := c.client.GenerateWinningNumbers(context.Background(), &types.Empty{})
	if err != nil {
		return nil
	}

	numbers := make([]int, len(winningNumbers.Numbers))
	for i, v := range winningNumbers.Numbers {
		numbers[i] = int(v)
	}
	drawDate := winningNumbers.DrawDate.AsTime()

	return &types.WinningNumbers{
		Numbers:  numbers,
		DrawDate: drawDate,
	}
}
