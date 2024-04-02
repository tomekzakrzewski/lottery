package main

import (
	"fmt"
	"time"

	"github.com/tomekzakrzewski/lottery/types"
)

type NumberReceiver interface {
	CreateTicket(nums *types.UserNumbers) (*types.Ticket, error)
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
	if !nums.ValidateNumbers() {
		return nil, fmt.Errorf("invalid numbers")
	}
	// DRAW DATE FROM ANOTHER SERVICE
	params := &types.CreateTicketParams{
		Numbers:  nums.Numbers,
		DrawDate: time.Now().Add(24 * time.Hour), // FETCH FROM ANOTHER SERVICE
	}

	ticket := types.NewTicketFromParams(params)

	res, err := n.ticketStore.Insert(ticket)
	if err != nil {
		return nil, err
	}
	return res, nil
}
