package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	uri := "mongodb://localhost:27017"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	store := NewTicketStore(client)
	svc := NewNumberReceiver(store)
	m := NewLogMiddleware(svc)

	srv := NewHttpTransport(m)

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	r := chi.NewRouter()
	r.Get("/drawDate", srv.handleGetNextDrawDate)
	r.Get("/ticket/{hash}", srv.handleFindByHash)
	r.Post("/ticket", srv.handlePostTicket)
	r.Post("/winningTickets", srv.handleGetWinningTickets)
	http.ListenAndServe(":3000", r)
}
