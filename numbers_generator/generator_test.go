package main

import (
	"testing"

	receiver "github.com/tomekzakrzewski/lottery/number_receiver/client"
)

func TestGenerateWinningNumbers(t *testing.T) {
	c, err := receiver.NewGRPCClient("localhost:3006")

	if err != nil {
		t.Error(err)
	}

	svc := NewGeneratorService(c)

	ticket := svc.GenerateWinningNumbers()
	if ticket == nil {
		t.Error("ticket is nil")
	}

	t.Log(ticket)
}
