package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tomekzakrzewski/lottery/types"
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

	p := &types.UserNumbers{
		Numbers: []int{1, 2, 3, 4, 5, 6},
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	r := chi.NewRouter()
	m.CreateTicket(p)
	http.ListenAndServe(":3000", r)
}
