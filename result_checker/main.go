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

	uri := "mongodb://localhost:27017"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	var (
		checkerGRPC = "localhost:3009"
	)
	receiverHTTPClient := receiver.NewHTTPClient("http://localhost:3000")
	//generatorHTTPClient := generator.NewHTTPClient("http://localhost:3001")
	generatorGRCPClient, err := generator.NewGRPCClient("localhost:3005")

	store := NewNumbersStore(client)
	svc := NewResultCheckerService(*receiverHTTPClient, generatorGRCPClient, *store)
	m := NewLogMiddleware(svc)
	srv := NewHttpTransport(m)

	go func() {
		log.Fatal(makeGRPCTransport(checkerGRPC, m))
	}()

	s := gocron.NewScheduler(time.UTC)
	//_, err = s.Every(1).Saturday().At("11:55").Do(m.GetWinningTickets)
	_, err = s.Every(1).Minutes().Do(m.GetWinningNumbers)
	if err != nil {
		fmt.Println(err)
	}
	s.StartAsync()

	r := chi.NewRouter()
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
