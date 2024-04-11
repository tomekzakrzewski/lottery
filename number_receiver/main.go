package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/tomekzakrzewski/lottery/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	uri := "mongodb://localhost:27017"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	var (
		receiverGRPC = "localhost:3006"
	)

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

	go func() {
		log.Fatal(makeGRPCTransport(receiverGRPC, m))
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

func makeGRPCTransport(listenAddr string, svc NumberReceiver) error {
	fmt.Println("GRPC running on ", listenAddr)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer func() {
		fmt.Println("stopping GRPC transport")
		ln.Close()
	}()
	grpcServer := grpc.NewServer([]grpc.ServerOption{}...)
	types.RegisterReceiverServer(grpcServer, NewGRPCReceiverServer(svc))
	return grpcServer.Serve(ln)
}
