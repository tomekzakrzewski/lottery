package main

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	receiver "github.com/tomekzakrzewski/lottery/number_receiver/client"
	generator "github.com/tomekzakrzewski/lottery/numbers_generator/client"
	"github.com/tomekzakrzewski/lottery/types"
)

type ResultChecker interface {
	GetWinningTickets() error
	IsTicketWinning(hash string) bool
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

func (r *ResultCheckerService) IsTicketWinning(hash string) bool {
	return r.winningTicketStore.CheckIfTicketIsWinning(hash)
}
