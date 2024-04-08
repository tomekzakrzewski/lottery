package main

import (
	"context"

	checker "github.com/tomekzakrzewski/lottery/result_checker/client"
)

type ResultAnnoucer interface {
	CheckTicketWin(hash string) bool
}

type ResultAnnoucerService struct {
	checker checker.HTTPClient
}

func NewResultAnnoucerService(checker checker.HTTPClient) *ResultAnnoucerService {
	return &ResultAnnoucerService{
		checker: checker,
	}
}

func (s *ResultAnnoucerService) CheckTicketWin(hash string) bool {
	winning := s.checker.IsTicketWinning(context.Background(), hash)
	return *winning
}
