package main

import (
	"time"

	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next ResultChecker
}

func NewLogMiddleware(next ResultChecker) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) GetWinningTickets() error {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
		}).Info("GetWinningTickets")
	}(time.Now())
	return m.next.GetWinningTickets()
}
