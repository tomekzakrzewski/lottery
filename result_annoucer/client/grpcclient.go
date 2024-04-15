package client

import (
	"context"

	"github.com/tomekzakrzewski/lottery/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Endpoint string
	client   types.AnnoucerClient
}

func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := types.NewAnnoucerClient(conn)

	return &GRPCClient{
		Endpoint: endpoint,
		client:   c,
	}, nil
}

func (c *GRPCClient) CheckResult(hash string) (*types.ResultResponse, error) {
	resultResp, err := c.client.CheckResult(context.Background(), &types.TicketHashRequest{Hash: hash})
	if err != nil {
		return nil, err
	}

	winningNums := make([]int, len(resultResp.WinningNumbers))
	for i, v := range resultResp.WinningNumbers {
		winningNums[i] = int(v)
	}

	numbers := make([]int, len(resultResp.Numbers))
	for i, v := range resultResp.Numbers {
		numbers[i] = int(v)
	}

	return &types.ResultResponse{
		Hash:           resultResp.Hash,
		Numbers:        numbers,
		WinningNumbers: winningNums,
		Win:            resultResp.Win,
		DrawDate:       resultResp.DrawDate.AsTime(),
	}, nil
}
