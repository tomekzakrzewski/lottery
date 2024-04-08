package main

import (
	"time"

	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next ResultAnnoucer
}

func NewLogMiddleware(next ResultAnnoucer) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) CheckTicketWin(hash string) bool {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"hash": hash,
		}).Info("GetWinningTickets")
	}(time.Now())
	return m.next.CheckTicketWin(hash)
}
