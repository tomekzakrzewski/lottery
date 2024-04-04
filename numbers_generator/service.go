package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/tomekzakrzewski/lottery/number_receiver/client"
	"github.com/tomekzakrzewski/lottery/types"
)

type Service struct {
	client client.HTTPClient
}

type GeneratorServicer interface {
	GenerateWinningNumbers() *types.WinningNumbers
}

func NewGeneratorService(client client.HTTPClient) *Service {
	return &Service{
		client: client,
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
	drawDate, err := s.client.GetNextDrawDate(context.Background())
	if err != nil {
		fmt.Println("failed to fetch draw date")
	}

	return &types.WinningNumbers{
		Numbers:  numbers,
		DrawDate: drawDate.Date,
	}
}
