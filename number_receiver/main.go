package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

	srv := NewHttpTransport(m)
	// test
	//dodaj(svc, store)

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

func dodaj(svc NumberReceiver, s *MongoTicketStore) {
	ticket := types.Ticket{
		Numbers:  []int{1, 2, 3, 4, 5, 6},
		DrawDate: svc.NextDrawDate().Date,
		Hash:     uuid.New().String(),
	}
	ticketAdded, _ := s.Insert(&ticket)
	fmt.Println("Ticket added: ", ticketAdded.Hash, ticketAdded.Numbers)

	ticket = types.Ticket{
		Numbers:  []int{2, 3, 4, 5, 6, 1},
		DrawDate: svc.NextDrawDate().Date,
		Hash:     uuid.New().String(),
	}
	ticketAdded, _ = s.Insert(&ticket)
	fmt.Println("Ticket added: ", ticketAdded.Hash, ticketAdded.Numbers)

	ticket = types.Ticket{
		Numbers:  []int{1, 2, 4, 5, 6, 3},
		DrawDate: svc.NextDrawDate().Date,
		Hash:     uuid.New().String(),
	}
	ticketAdded, _ = s.Insert(&ticket)
	fmt.Println("Ticket added: ", ticketAdded.Hash, ticketAdded.Numbers)
}
