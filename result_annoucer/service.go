package main

import (
	"context"
	"fmt"

	receiver "github.com/tomekzakrzewski/lottery/number_receiver/client"
	checker "github.com/tomekzakrzewski/lottery/result_checker/client"
	"github.com/tomekzakrzewski/lottery/types"
)

type ResultAnnoucer interface {
	CheckResult(hash string) (*types.ResultRespone, error)
}

type ResultAnnoucerService struct {
	checker  checker.HTTPClient
	receiver receiver.HTTPClient
	redis    *RedisStore
}

func NewResultAnnoucerService(checker checker.HTTPClient, receiver receiver.HTTPClient, redis *RedisStore) *ResultAnnoucerService {
	return &ResultAnnoucerService{
		checker:  checker,
		receiver: receiver,
		redis:    redis,
	}
}

func (s *ResultAnnoucerService) CheckResult(hash string) (*types.ResultRespone, error) {
	// sprawdzic czy hash jest w redis
	result, _ := s.redis.Find(hash)
	if result != nil {
		return result, nil
	}
	// brak w redisie, pobranie z receivera
	ticket, err := s.receiver.GetTicketByHash(context.Background(), hash)
	if err != nil {
		return nil, fmt.Errorf("ticket with hash %s not found", hash)
	}

	// sprawdzic czy ticket jest wygrany
	ticketWin := s.checker.IsTicketWinning(context.Background(), hash)

	// stworzenie result z ticketa i ticket win
	resultTicket := types.ResultRespone{
		Hash:     ticket.Hash,
		Numbers:  ticket.Numbers,
		Win:      *ticketWin,
		DrawDate: ticket.DrawDate,
	}

	// zapisanie result do redisa
	err = s.redis.Insert(&resultTicket)
	if err != nil {
		return nil, err
	}

	return &resultTicket, nil
}
