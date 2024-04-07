package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	receiver "github.com/tomekzakrzewski/lottery/number_receiver/client"
	generator "github.com/tomekzakrzewski/lottery/numbers_generator/client"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	uri := "mongodb://localhost:27017"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	receiverHTTPClient := receiver.NewHTTPClient("http://localhost:3000")
	generatorHTTPClient := generator.NewHTTPClient("http://localhost:3001")
	store := NewWinningTicketStore(client)
	svc := NewResultCheckerService(*receiverHTTPClient, *generatorHTTPClient, *store)
	m := NewLogMiddleware(svc)
	srv := NewHttpTransport(m)
	r := chi.NewRouter()
	m.GetWinningTickets()

	r.Get("/win/{hash}", srv.handleCheckIsTicketWinning)

	http.ListenAndServe(":5000", r)
}
