package client

import (
	"context"
	"time"

	"github.com/tomekzakrzewski/lottery/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (c *GRPCClient) GetNextDrawDate(ctx context.Context, currentTime time.Time) *types.DrawDate {
	drawDate, err := c.client.NextDrawDate(context.Background(), &types.NextDateRequest{
		Date: timestamppb.New(currentTime)})
	if err != nil {
		return nil
	}

	return &types.DrawDate{
		Date: drawDate.Date.AsTime(),
	}
}

func (c *GRPCClient) CreateTicket(ctx context.Context, nums *types.UserNumbers) (*types.Ticket, error) {
	numbers := make([]int32, len(nums.Numbers))
	for i, v := range nums.Numbers {
		numbers[i] = int32(v)
	}
	ticket, err := c.client.CreateTicket(ctx, &types.UserNumbersTransport{Numbers: numbers})
	if err != nil {
		return nil, err
	}
	id, err := primitive.ObjectIDFromHex(ticket.Id)
	if err != nil {
		return nil, err
	}
	return &types.Ticket{
		ID:       id,
		Numbers:  nums.Numbers,
		DrawDate: ticket.DrawDate.AsTime(),
		Hash:     ticket.Hash,
	}, nil
}
