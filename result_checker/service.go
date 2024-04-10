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
	receiver  receiver.HTTPClient
	generator generator.HTTPClient
	store     *MongoNumbersStore
}

func NewResultCheckerService(receiver receiver.HTTPClient, generator generator.HTTPClient, store MongoNumbersStore) *ResultCheckerService {
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

/*
// RAN BY SCHEDULER, czy wystawiac na to endpoint wgl?
func (r *ResultCheckerService) GetWinningTickets() error {
	winningNumbers := r.generator.GenerateWinningNumbers()
	winningNumbersMock := types.WinningNumbers{
		Numbers:  []int{1, 2, 3, 4, 5, 6},
		DrawDate: winningNumbers.DrawDate,
	}
	//	winningTickets, _ := r.receiver.GetWinningTickets(context.Background(), *winningNumbers)
	winningTickets, _ := r.receiver.GetWinningTickets(context.Background(), winningNumbersMock)
	fmt.Println(len(winningTickets))
	err := r.winningTicketStore.InsertWinningTickets(winningTickets)
	if err != nil {
		logrus.Info(err)
		return err
	}

	return nil
}
*/

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
