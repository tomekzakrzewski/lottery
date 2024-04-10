package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tomekzakrzewski/lottery/types"
)

type LogMiddleware struct {
	next ResultAnnoucer
}

func NewLogMiddleware(next ResultAnnoucer) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) CheckResult(hash string) (result *types.ResultRespone, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":    time.Since(start),
			"error":   err,
			"hash":    hash,
			"numbers": result.Numbers,
			"win":     result.Win,
			"date":    result.DrawDate,
		}).Info("CheckResult")
	}(time.Now())
	return m.next.CheckResult(hash)
}
