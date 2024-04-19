package client

import (
	"context"
	"fmt"

	"github.com/tomekzakrzewski/lottery/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCClient struct {
	Endpoint string
	client   types.CheckerClient
}

func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	endpointUrl := fmt.Sprintf("checker%s", endpoint)
	conn, err := grpc.Dial(endpointUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := types.NewCheckerClient(conn)

	return &GRPCClient{
		Endpoint: endpoint,
		client:   c,
	}, nil
}

func (c *GRPCClient) IsTicketWinning(ctx context.Context, ticket *types.Ticket) (*types.ResultResponse, error) {

	numbers := make([]int32, len(ticket.Numbers))
	for i, v := range ticket.Numbers {
		numbers[i] = int32(v)
	}

	ticketTransport := &types.TicketTransport{
		Id:       ticket.ID.Hex(),
		Numbers:  numbers,
		DrawDate: timestamppb.New(ticket.DrawDate),
		Hash:     ticket.Hash,
	}

	resultResp, err := c.client.CheckTicket(ctx, ticketTransport)
	if err != nil {
		return nil, err
	}

	winningNums := make([]int, len(resultResp.Numbers))
	for i, v := range resultResp.Numbers {
		winningNums[i] = int(v)
	}

	resultResponse := &types.ResultResponse{
		Hash:           resultResp.Hash,
		Numbers:        ticket.Numbers,
		WinningNumbers: winningNums,
		Win:            resultResp.Win,
		DrawDate:       resultResp.DrawDate.AsTime(),
	}

	return resultResponse, nil
}
