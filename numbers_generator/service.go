package main

import (
	"math/rand"
	"time"

	"github.com/tomekzakrzewski/lottery/types"
)

type Service struct{}

type GeneratorServicer interface {
	GenerateWinningNumbers() *types.WinningNumbers
}

func NewGeneratorService() *Service {
	return &Service{}
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
	return &types.WinningNumbers{
		Numbers:  numbers,
		DrawDate: time.Now(), //fetch from number receiver
	}
}
