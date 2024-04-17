package main

import (
	"reflect"

	receiver "github.com/tomekzakrzewski/lottery/number_receiver/client"
	generator "github.com/tomekzakrzewski/lottery/numbers_generator/client"
	"github.com/tomekzakrzewski/lottery/types"
)

type ResultChecker interface {
	GetWinningNumbers() error
	CheckTicketWin(ticket *types.Ticket) (*types.ResultResponse, error)
}

type ResultCheckerService struct {
	receiver  receiver.Client
	generator generator.Client
	store     *MongoNumbersStore
}

func NewResultCheckerService(receiver receiver.Client, generator generator.Client, store MongoNumbersStore) *ResultCheckerService {
	return &ResultCheckerService{
		receiver:  receiver,
		generator: generator,
		store:     &store,
	}
}

// RAN BY SCHEDULER
func (r *ResultCheckerService) GetWinningNumbers() error {
	winningNumbers := r.generator.GenerateWinningNumbers()

	_, err := r.store.InsertWinningNumbers(winningNumbers)
	if err != nil {
		return err
	}
	return nil
}

func (r *ResultCheckerService) CheckTicketWin(ticket *types.Ticket) (*types.ResultResponse, error) {
	winningNumbers, err := r.store.FindWinningNumbersByDate(ticket.DrawDate)
	if err != nil {
		return nil, err
	}

	win := reflect.DeepEqual(ticket.Numbers, winningNumbers.Numbers)
	return &types.ResultResponse{
		Hash:           ticket.Hash,
		Numbers:        ticket.Numbers,
		Win:            win,
		WinningNumbers: winningNumbers.Numbers,
		DrawDate:       ticket.DrawDate,
	}, nil
}
