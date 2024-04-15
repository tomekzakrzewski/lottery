package client

import (
	"context"

	"github.com/tomekzakrzewski/lottery/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Endpoint string
	client   types.ReceiverClient
}

func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := types.NewReceiverClient(conn)

	return &GRPCClient{
		Endpoint: endpoint,
		client:   c,
	}, nil
}

func (c *GRPCClient) GetTicketByHash(ctx context.Context, hash string) (*types.Ticket, error) {
	ticket, err := c.client.GetTicketByHash(ctx, &types.TicketHashRequest{Hash: hash})
	if err != nil {
		return nil, err
	}

	numbers := make([]int, len(ticket.Numbers))
	for i, v := range ticket.Numbers {
		numbers[i] = int(v)
	}

	id, err := primitive.ObjectIDFromHex(ticket.Id)
	if err != nil {
		return nil, err
	}

	return &types.Ticket{
		ID:       id,
		Numbers:  numbers,
		DrawDate: ticket.DrawDate.AsTime(),
		Hash:     ticket.Hash,
	}, nil
}

func (c *GRPCClient) GetNextDrawDate() *types.DrawDate {
	drawDate, err := c.client.NextDrawDate(context.Background(), &types.Empty{})
	if err != nil {
		return nil
	}

	return &types.DrawDate{
		Date: drawDate.Date.AsTime(),
	}
}
