package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	checker "github.com/tomekzakrzewski/lottery/result_checker/client"
)

func main() {
	checkerClient := checker.NewHTTPClient("http://localhost:5000")
	svc := NewResultAnnoucerService(*checkerClient)
	m := NewLogMiddleware(svc)
	srv := NewHttpTransport(m)

	r := chi.NewRouter()
	r.Get("/win/{hash}", srv.handleCheckIsTicketWinning)

	http.ListenAndServe(":6000", r)

}
