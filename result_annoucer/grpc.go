package main

import (
	"context"

	"github.com/tomekzakrzewski/lottery/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCAnnoucerServer struct {
	types.UnimplementedAnnoucerServer
	svc ResultAnnoucer
}

func NewAnnoucerGRPCServer(svc ResultAnnoucer) *GRPCAnnoucerServer {
	return &GRPCAnnoucerServer{
		svc: svc,
	}
}

func (s *GRPCAnnoucerServer) CheckResult(ctx context.Context, req *types.TicketHashRequest) (*types.ResultResp, error) {
	hash := req.GetHash()
	ticket, err := s.svc.CheckResult(hash)
	if err != nil {
		return nil, err
	}

	winningNums := make([]int32, len(ticket.WinningNumbers))
	for i, v := range ticket.WinningNumbers {
		winningNums[i] = int32(v)
	}

	numbers := make([]int32, len(ticket.Numbers))
	for i, v := range ticket.Numbers {
		numbers[i] = int32(v)
	}

	return &types.ResultResp{
		Hash:           ticket.Hash,
		Numbers:        numbers,
		WinningNumbers: winningNums,
		Win:            ticket.Win,
		DrawDate:       timestamppb.New(ticket.DrawDate),
	}, nil
}
