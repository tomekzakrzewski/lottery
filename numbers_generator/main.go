package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tomekzakrzewski/lottery/number_receiver/client"
)

func main() {

	var (
		receiverEndpoint = "http://localhost:3000"
	)

	svc := NewGeneratorService()
	m := NewLogMiddleware(svc)

	srv := NewHttpTransport(m)

	httpClient := client.NewHTTPClient(receiverEndpoint)
	date, err := httpClient.GetNextDrawDate(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(date)
	r := chi.NewRouter()
	r.Get("/winningNumbers", srv.handleGetWinningNumbers)
	http.ListenAndServe(":3001", r)

}
