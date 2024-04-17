package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/tomekzakrzewski/lottery/number_receiver/client"
	"github.com/tomekzakrzewski/lottery/types"
)

type Service struct {
	receiver client.Client
}

type GeneratorServicer interface {
	GenerateWinningNumbers() *types.WinningNumbers
}

func NewGeneratorService(receiver client.Client) *Service {
	return &Service{
		receiver: receiver,
	}
}

func (s *Service) GenerateWinningNumbers() *types.WinningNumbers {
	rand.Seed(time.Now().UnixNano())
	uniqueNumbers := make(map[int]bool)

	for len(uniqueNumbers) < 6 {
		randomNumber := rand.Intn(99) + 1
		uniqueNumbers[randomNumber] = true
	}

	var numbers []int
	for key := range uniqueNumbers {
		numbers = append(numbers, key)
	}
	drawDate := s.receiver.GetNextDrawDate(context.Background(), time.Now())

	return &types.WinningNumbers{
		Numbers:  numbers,
		DrawDate: drawDate.Date,
	}
}
