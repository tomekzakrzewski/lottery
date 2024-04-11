package main

import (
	"context"

	"github.com/tomekzakrzewski/lottery/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCReceiverServer struct {
	types.UnimplementedReceiverServer
	svc NumberReceiver
}

func NewGRPCReceiverServer(svc NumberReceiver) *GRPCReceiverServer {
	return &GRPCReceiverServer{
		svc: svc,
	}
}

func (s *GRPCReceiverServer) GetTicketByHash(ctx context.Context, req *types.TicketHashRequest) (*types.TicketTransport, error) {
	hash := req.GetHash()
	ticket, err := s.svc.GetTicketByHash(hash)
	if err != nil {
		return nil, err
	}

	numbers := make([]int32, len(ticket.Numbers))
	for i, v := range ticket.Numbers {
		numbers[i] = int32(v)
	}

	return &types.TicketTransport{
		Id:       ticket.ID.Hex(),
		Numbers:  numbers,
		DrawDate: timestamppb.New(ticket.DrawDate),
		Hash:     ticket.Hash,
	}, nil
}

func (s *GRPCReceiverServer) NextDrawDate(ctx context.Context, req *types.Empty) (*types.NextDate, error) {
	drawDate := s.svc.NextDrawDate()
	return &types.NextDate{
		Date: timestamppb.New(drawDate.Date),
	}, nil
}
