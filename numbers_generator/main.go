package main

func main() {

	svc := NewGeneratorService()
	m := NewLogMiddleware(svc)

	m.generateWinningNumbers()
}
