package main

import (
	"fmt"
	"time"

	"github.com/tomekzakrzewski/lottery/types"
)

type NumberReceiver interface {
	CreateTicket(nums *types.UserNumbers) (*types.Ticket, error)
	NextDrawDate() types.DrawDate
	GetTicketByHash(hash string) (*types.Ticket, error)
	GetWinningTickets(winningNumbers types.WinningNumbers) ([]*types.Ticket, error)
}

type ReceiverService struct {
	ticketStore *MongoTicketStore
}

func NewNumberReceiver(ticketStore *MongoTicketStore) NumberReceiver {
	return &ReceiverService{
		ticketStore: ticketStore,
	}
}

func (n *ReceiverService) CreateTicket(nums *types.UserNumbers) (*types.Ticket, error) {
	// validacja
	// better error handling
	if !nums.ValidateNumbers() {
		return nil, fmt.Errorf("invalid numbers")
	}
	params := &types.CreateTicketParams{
		Numbers:  nums.Numbers,
		DrawDate: n.NextDrawDate().Date,
	}

	ticket := types.NewTicketFromParams(params)

	res, err := n.ticketStore.Insert(ticket)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// every saturday at 12:00
func (s *ReceiverService) NextDrawDate() types.DrawDate {
	currentTime := time.Now()

	// If it's Saturday and before noon, return today's date at draw time
	if currentTime.Weekday() == time.Saturday && currentTime.Hour() < 12 {
		drawDate := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 12, 0, 0, 0, currentTime.UTC().Location())
		return types.DrawDate{
			Date: drawDate,
		}
	}

	// Otherwise, find the next Saturday and return its date at draw time
	for currentTime.Weekday() != time.Saturday {
		currentTime = currentTime.AddDate(0, 0, 1)
	}
	drawDate := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 12, 0, 0, 0, currentTime.UTC().Location())
	return types.DrawDate{
		Date: drawDate,
	}
}

func (n *ReceiverService) GetTicketByHash(hash string) (*types.Ticket, error) {
	ticket, err := n.ticketStore.FindByHash(hash)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}

func (n *ReceiverService) GetWinningTickets(winningNumbers types.WinningNumbers) ([]*types.Ticket, error) {
	fmt.Println("IN THE SERVICE")
	fmt.Println(winningNumbers.Numbers)
	fmt.Println(winningNumbers.DrawDate)
	fmt.Println("---------------------")

	tickets, err := n.ticketStore.FindByWinningNumersAndDrawDate(winningNumbers)
	if err != nil {
		return []*types.Ticket{}, nil
	}
	return tickets, nil
}
