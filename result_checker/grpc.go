package main

import (
	"context"

	"github.com/tomekzakrzewski/lottery/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCCheckerServer struct {
	types.UnimplementedCheckerServer
	svc ResultChecker
}

func NewGRPCCheckerServer(svc ResultChecker) *GRPCCheckerServer {
	return &GRPCCheckerServer{
		svc: svc,
	}
}

func (s *GRPCCheckerServer) CheckTicket(ctx context.Context, req *types.TicketTransport) (*types.ResultResp, error) {
	numbers := make([]int, len(req.Numbers))
	for i, v := range req.Numbers {
		numbers[i] = int(v)
	}

	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	ticket := &types.Ticket{
		ID:       id,
		Numbers:  numbers,
		DrawDate: req.DrawDate.AsTime(),
		Hash:     req.Hash,
	}

	ticketResp, err := s.svc.CheckTicketWin(ticket)
	if err != nil {
		return nil, err
	}

	winningNums := make([]int32, len(ticketResp.Numbers))
	for i, v := range ticketResp.Numbers {
		winningNums[i] = int32(v)
	}
	resultResp := &types.ResultResp{
		Hash:           ticketResp.Hash,
		Numbers:        req.Numbers,
		WinningNumbers: winningNums,
		Win:            ticketResp.Win,
		DrawDate:       timestamppb.New(ticketResp.DrawDate),
	}

	return resultResp, nil
}
