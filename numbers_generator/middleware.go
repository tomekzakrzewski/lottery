package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tomekzakrzewski/lottery/types"
)

type LogMiddleware struct {
	next GeneratorServicer
}

func NewLogMiddleware(next GeneratorServicer) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) generateWinningNumbers() (winningNumbers *types.WinningNumbers) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":            time.Since(start),
			"winning numbers": winningNumbers,
			"date":            winningNumbers.DrawDate,
		}).Info("generate winning numbers")
	}(time.Now())
	winningNumbers = m.next.GenerateWinningNumbers()
	return
}
