package main

import (
	"context"
	"fmt"
	"time"

	receiver "github.com/tomekzakrzewski/lottery/number_receiver/client"
	checker "github.com/tomekzakrzewski/lottery/result_checker/client"
	"github.com/tomekzakrzewski/lottery/types"
)

type ResultAnnoucer interface {
	CheckResult(hash string) (*types.ResultResponse, error)
}

type ResultAnnoucerService struct {
	checker  checker.GRPCClient
	receiver receiver.GRPCClient
	redis    *RedisStore
}

func NewResultAnnoucerService(checker checker.GRPCClient, receiver receiver.GRPCClient, redis *RedisStore) *ResultAnnoucerService {
	return &ResultAnnoucerService{
		checker:  checker,
		receiver: receiver,
		redis:    redis,
	}
}

func (s *ResultAnnoucerService) CheckResult(hash string) (*types.ResultResponse, error) {
	// sprawdzic czy hash jest w redis
	result, _ := s.redis.Find(hash)
	fmt.Println(result)
	if result != nil {
		return result, nil
	}
	// brak w redisie, pobranie z receivera
	ticket, err := s.receiver.GetTicketByHash(context.Background(), hash)
	fmt.Println(ticket)
	if err != nil {
		return nil, fmt.Errorf("ticket with hash %s not found", hash)
	}

	// sprawdzic czy ticket jest wygrany
	resultTicket, err := s.checker.IsTicketWinning(context.Background(), ticket)
	fmt.Println(resultTicket)
	if err != nil {
		return nil, err
	}

	// zapisanie result do redisa
	err = s.redis.Insert(resultTicket)
	if err != nil {
		return nil, err
	}

	if !isAfterAnnoucement(resultTicket.DrawDate) {
		return nil, fmt.Errorf("result will be annouced after %s", resultTicket.DrawDate)
	}
	return resultTicket, nil
}

func isAfterAnnoucement(date time.Time) bool {
	now := time.Now()
	return now.After(date)
}
