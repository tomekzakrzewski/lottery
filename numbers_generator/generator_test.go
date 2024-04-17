package main

import (
	"testing"

	receiver "github.com/tomekzakrzewski/lottery/number_receiver/client"
)

func TestGenerateWinningNumbers(t *testing.T) {

	c := receiver.NewHTTPClient("localhost:3000")
	svc := NewGeneratorService(c)

	ticket := svc.GenerateWinningNumbers()

	t.Log(ticket)
}
