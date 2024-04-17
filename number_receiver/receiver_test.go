package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/tomekzakrzewski/lottery/types"
)

func TestCreateTicketValidNumbers(t *testing.T) {
	store := NewInMemoryTicketStore()
	svc := NewNumberReceiver(store)

	numbers := &types.UserNumbers{
		Numbers: []int{1, 2, 3, 4, 5, 6},
	}

	ticket, err := svc.CreateTicket(numbers)
	if err != nil {
		t.Error(err)
	}

	if ticket == nil {
		t.Error("ticket is nil, not created")
	}
}

func TestCreateTicketInvalidNumbers(t *testing.T) {
	store := NewInMemoryTicketStore()
	svc := NewNumberReceiver(store)

	fiveNumbers := &types.UserNumbers{
		Numbers: []int{2, 3, 4, 5, 6},
	}
	_, err := svc.CreateTicket(fiveNumbers)
	if err == nil {
		t.Error("should not be able to create ticket with 5 numbers")
	}

	negativeNumbers := &types.UserNumbers{
		Numbers: []int{-1, 2, 3, 4, 5, 6},
	}
	_, err = svc.CreateTicket(negativeNumbers)
	if err == nil {
		t.Error("should not be able to create ticket with negative number")
	}

	tooBigNumbers := &types.UserNumbers{
		Numbers: []int{100, 2, 3, 4, 5, 6},
	}
	_, err = svc.CreateTicket(tooBigNumbers)
	if err == nil {
		t.Error("should not be able to create ticket with number bigger than 99")
	}

	tooSmallNumbers := &types.UserNumbers{
		Numbers: []int{0, 2, 3, 4, 5, 6},
	}
	_, err = svc.CreateTicket(tooSmallNumbers)
	if err == nil {
		t.Error("should not be able to create ticket with number small than 1")
	}

	nonUniqueNumbers := &types.UserNumbers{
		Numbers: []int{1, 1, 2, 3, 4, 5},
	}
	_, err = svc.CreateTicket(nonUniqueNumbers)
	if err == nil {
		t.Error("should not be able to create ticket with number small than 1")
	}
}

func TestGetTicketByHashValid(t *testing.T) {
	store := NewInMemoryTicketStore()
	svc := NewNumberReceiver(store)

	numbers := &types.UserNumbers{
		Numbers: []int{1, 2, 3, 4, 5, 6},
	}
	insertedTicket, err := svc.CreateTicket(numbers)

	if err != nil {
		t.Error(err)
	}

	ticket, err := svc.GetTicketByHash(insertedTicket.Hash)
	if err != nil {
		t.Error(err)
	}
	if ticket == nil {
		t.Error("ticket is nil, not created")
	}
}

func TestGetTicketByHashInvalid(t *testing.T) {
	store := NewInMemoryTicketStore()
	svc := NewNumberReceiver(store)

	_, err := svc.GetTicketByHash("invalid-hash")
	if err == nil {
		t.Error("should not be able to get ticket with invalid hash")
	}
}

func TestNextDrawDateSaturdayAfterNoon(t *testing.T) {
	store := NewInMemoryTicketStore()
	svc := NewNumberReceiver(store)
	currentTime := time.Date(2024, time.April, 20, 12, 0, 0, 0, time.UTC).UTC()
	expectedDrawDate := time.Date(2024, time.April, 27, 12, 0, 0, 0, time.UTC).UTC()
	drawDate := svc.NextDrawDate(currentTime)

	if drawDate.Date != expectedDrawDate {
		t.Errorf("Expected draw date: %v, got: %v", expectedDrawDate, drawDate.Date)
	}
}

func TestNextDrawDateNotSaturday(t *testing.T) {
	store := NewInMemoryTicketStore()
	svc := NewNumberReceiver(store)

	currentTime := time.Date(2024, time.April, 19, 13, 0, 0, 0, time.UTC)
	expectedDrawDate := time.Date(2024, time.April, 20, 12, 0, 0, 0, time.UTC)
	drawDate := svc.NextDrawDate(currentTime)
	if drawDate.Date != expectedDrawDate {
		t.Errorf("Expected draw date: %v, got: %v", expectedDrawDate, drawDate.Date)
	}
}

type InMemoryTicketStore struct {
	tickets []*types.Ticket
}

func NewInMemoryTicketStore() *InMemoryTicketStore {
	return &InMemoryTicketStore{
		tickets: []*types.Ticket{},
	}
}

func (s *InMemoryTicketStore) Insert(ticket *types.Ticket) (*types.Ticket, error) {
	s.tickets = append(s.tickets, ticket)
	return ticket, nil
}

func (s *InMemoryTicketStore) FindByHash(hash string) (*types.Ticket, error) {
	for _, ticket := range s.tickets {
		if ticket.Hash == hash {
			return ticket, nil
		}
	}
	return nil, fmt.Errorf("ticket not found")
}
