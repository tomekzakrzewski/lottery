package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tomekzakrzewski/lottery/types"
)

type LogMiddleware struct {
	next NumberReceiver
}

func NewLogMiddleware(next NumberReceiver) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) CreateTicket(nums *types.UserNumbers) (ticket *types.Ticket, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":        time.Since(start),
			"err":         err,
			"ticketID":    ticket.ID.Hex(),
			"ticket hash": ticket.Hash,
			"numbers":     nums,
			"draw date":   ticket.DrawDate,
		}).Info("CreateTicket")
	}(time.Now())
	ticket, err = m.next.CreateTicket(nums)
	return
}

func (m *LogMiddleware) NextDrawDate(currentTime time.Time) (nextDraw types.DrawDate) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":        time.Since(start),
			"date":        nextDraw,
			"currentTime": currentTime,
		}).Info("NextDrawDate")
	}(time.Now())
	nextDraw = m.next.NextDrawDate(currentTime)
	return
}

func (m *LogMiddleware) GetTicketByHash(hash string) (ticket *types.Ticket, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
			"hash": hash,
			"date": ticket.DrawDate,
		}).Info("GetTicketByHash")
	}(time.Now())
	ticket, err = m.next.GetTicketByHash(hash)
	return
}
