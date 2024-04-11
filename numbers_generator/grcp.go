package main

import (
	"context"
	"fmt"

	"github.com/tomekzakrzewski/lottery/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRCPGeneratorServer struct {
	types.UnimplementedGeneratorServer
	svc GeneratorServicer
}

func NewGeneratorGRPCServer(svc GeneratorServicer) *GRCPGeneratorServer {
	return &GRCPGeneratorServer{
		svc: svc,
	}
}

func (s *GRCPGeneratorServer) GenerateWinningNumbers(ctx context.Context, req *types.Empty) (*types.WinningNums, error) {
	winningNumbers := s.svc.GenerateWinningNumbers()
	fmt.Println(winningNumbers.DrawDate, winningNumbers.DrawDate)
	numbers := make([]int32, len(winningNumbers.Numbers))
	for i, v := range winningNumbers.Numbers {
		numbers[i] = int32(v)
	}

	date := timestamppb.New(winningNumbers.DrawDate)

	return &types.WinningNums{
		Numbers:  numbers,
		DrawDate: date,
	}, nil
}
