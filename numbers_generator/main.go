package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tomekzakrzewski/lottery/number_receiver/client"
)

func main() {

	var (
		receiverEndpoint = "http://localhost:3000"
	)

	httpClient := client.NewHTTPClient(receiverEndpoint)
	svc := NewGeneratorService(*httpClient)
	m := NewLogMiddleware(svc)

	srv := NewHttpTransport(m)

	r := chi.NewRouter()
	r.Get("/winningNumbers", srv.handleGetWinningNumbers)
	http.ListenAndServe(":3001", r)
}
