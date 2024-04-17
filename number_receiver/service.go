package main

import (
	"fmt"
	"time"

	"github.com/tomekzakrzewski/lottery/types"
)

var drawTime = time.Date(0, 0, 0, 12, 0, 0, 0, time.UTC)

type NumberReceiver interface {
	CreateTicket(nums *types.UserNumbers) (*types.Ticket, error)
	NextDrawDate(currentTime time.Time) types.DrawDate
	GetTicketByHash(hash string) (*types.Ticket, error)
}

type ReceiverService struct {
	ticketStore TicketStore
}

func NewNumberReceiver(ticketStore TicketStore) NumberReceiver {
	return &ReceiverService{
		ticketStore: ticketStore,
	}
}

func (n *ReceiverService) CreateTicket(nums *types.UserNumbers) (*types.Ticket, error) {
	if !nums.ValidateNumbers() {
		return nil, fmt.Errorf("invalid numbers")
	}
	params := &types.CreateTicketParams{
		Numbers:  nums.Numbers,
		DrawDate: n.NextDrawDate(time.Now()).Date,
	}

	ticket := types.NewTicketFromParams(params)

	res, err := n.ticketStore.Insert(ticket)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *ReceiverService) NextDrawDate(currentTime time.Time) types.DrawDate {
	if s.isSaturdayAndBeforeNoon(currentTime) {
		time := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), drawTime.Hour(), drawTime.Minute(), drawTime.Second(), 0, currentTime.Location())
		return types.DrawDate{Date: time}
	}
	nextSaturday := currentTime.AddDate(0, 0, 1)
	for nextSaturday.Weekday() != time.Saturday {
		nextSaturday = nextSaturday.AddDate(0, 0, 1)
	}

	time := time.Date(nextSaturday.Year(), nextSaturday.Month(), nextSaturday.Day(), drawTime.Hour(), drawTime.Minute(), drawTime.Second(), 0, nextSaturday.Location())
	return types.DrawDate{Date: time}
}

func (s *ReceiverService) isSaturdayAndBeforeNoon(currentDateTime time.Time) bool {
	return currentDateTime.Weekday() == time.Saturday && currentDateTime.Hour() < drawTime.Hour()
}

func (n *ReceiverService) GetTicketByHash(hash string) (*types.Ticket, error) {
	ticket, err := n.ticketStore.FindByHash(hash)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}
