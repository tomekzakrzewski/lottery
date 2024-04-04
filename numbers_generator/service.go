package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
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
		DrawDate: fetchNextDrawDate().Date,
	}
}

func fetchNextDrawDate() *types.DrawDate {
	apiUrl := "http://localhost:3000/drawDate"
	response, err := http.Get(apiUrl)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	var date types.DrawDate
	err = json.NewDecoder(response.Body).Decode(&date)
	if err != nil {
		panic(err)
	}

	return &date
}
