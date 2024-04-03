package main

import "time"

func main() {

	svc := NewGeneratorService()
	m := NewLogMiddleware(svc)

	m.generateWinningNumbers()
	time.Sleep(5 * time.Minute)
}
