package main

import (
	"context"

	receiver "github.com/tomekzakrzewski/lottery/number_receiver/client"
	generator "github.com/tomekzakrzewski/lottery/numbers_generator/client"
)

type ResultChecker interface {
	GetWinningTickets(hash string) (bool, error)
}

type ResultCheckerService struct {
	receiver           receiver.HTTPClient
	generator          generator.HTTPClient
	winningTicketStore *MongoWinningTicketStore
}

func NewResultCheckerService(receiver receiver.HTTPClient, generator generator.HTTPClient, store MongoWinningTicketStore) *ResultCheckerService {
	return &ResultCheckerService{
		receiver:           receiver,
		generator:          generator,
		winningTicketStore: &store,
	}
}

func (r *ResultCheckerService) GetWinningTickets() error {
	winningNumbers := r.generator.GenerateWinningNumbers()
	winningTickets, _ := r.receiver.GetWinningTickets(context.Background(), *winningNumbers)

	err := r.winningTicketStore.InsertWinningTickets(winningTickets)
	if err != nil {
		return err
	}

	return nil
}
