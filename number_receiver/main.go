package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/tomekzakrzewski/lottery/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	mongoClient := makeClient("mongodb://localhost:27017")

	var (
		receiverGRPC = "localhost:3006"
		receiverHTTP = ":3000"
		store        = NewTicketStore(mongoClient)
		svc          = NewNumberReceiver(store)
		m            = NewLogMiddleware(svc)
	)

	srv := NewHttpTransport(m)
	defer func() {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
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
	http.ListenAndServe(receiverHTTP, r)
}

func makeClient(uri string) *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	return client
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
