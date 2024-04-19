package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-co-op/gocron"
	receiver "github.com/tomekzakrzewski/lottery/number_receiver/client"
	generator "github.com/tomekzakrzewski/lottery/numbers_generator/client"
	"github.com/tomekzakrzewski/lottery/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://mongodb:27017"))
	if err != nil {
		log.Fatal(err)
	}
	receiverGRPCClient, _ := receiver.NewGRPCClient(":3009")
	generatorGRCPClient, _ := generator.NewGRPCClient(":3005")
	store := NewNumbersStore(client)
	r := chi.NewRouter()
	svc := NewResultCheckerService(receiverGRPCClient, generatorGRCPClient, *store)
	m := NewLogMiddleware(svc)
	srv := NewHttpTransport(m)

	go func() {
		log.Fatal(makeGRPCTransport(":3009", m))
	}()

	s := gocron.NewScheduler(time.Local)
	//_, err = s.Every(1).Saturday().At("11:55").Do(m.GetWinningNumbers)
	_, err = s.Every(1).Minutes().Do(m.GetWinningNumbers)
	if err != nil {
		log.Fatal(err)
	}
	s.StartAsync()

	r.Post("/win", srv.handleCheckIsTicketWinning)
	http.ListenAndServe(":5000", r)
}

func makeGRPCTransport(listenAddr string, svc ResultChecker) error {
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
	types.RegisterCheckerServer(grpcServer, NewGRPCCheckerServer(svc))
	return grpcServer.Serve(ln)
}
