package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tomekzakrzewski/lottery/types"
)

type LogMiddleware struct {
	next ResultChecker
}

func NewLogMiddleware(next ResultChecker) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) GetWinningNumbers() error {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
		}).Info("GetWinningNumbers")
	}(time.Now())
	return m.next.GetWinningNumbers()
}

func (m *LogMiddleware) CheckTicketWin(ticket *types.Ticket) (result *types.ResultResponse, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":    time.Since(start),
			"error":   err,
			"hash":    ticket.Hash,
			"numbers": ticket.Numbers,
			"win":     result.Win,
			"date":    result.DrawDate,
		}).Info("CheckTicketWin")
	}(time.Now())
	return m.next.CheckTicketWin(ticket)
}
