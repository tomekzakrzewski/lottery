package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	receiver "github.com/tomekzakrzewski/lottery/number_receiver/client"
	generator "github.com/tomekzakrzewski/lottery/numbers_generator/client"
	"github.com/tomekzakrzewski/lottery/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	mongoClient := makeClient(os.Getenv("MONGO_URI"))
	receiverGRPCClient, _ := receiver.NewGRPCClient(os.Getenv("RECEIVER_GRPC"))
	generatorGRCPClient, _ := generator.NewGRPCClient("localhost:3005")
	var (
		checkerGRPC = os.Getenv("CHECKER_GRPC")
		checkerHTTP = os.Getenv("CHECKER_HTTP")
		store       = NewNumbersStore(mongoClient)
		r           = chi.NewRouter()
		svc         = NewResultCheckerService(receiverGRPCClient, generatorGRCPClient, *store)
		m           = NewLogMiddleware(svc)
		srv         = NewHttpTransport(m)
	)

	go func() {
		log.Fatal(makeGRPCTransport(checkerGRPC, m))
	}()

	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(1).Saturday().At("11:55").Do(m.GetWinningNumbers)
	//_, err = s.Every(1).Minutes().Do(m.GetWinningNumbers)
	if err != nil {
		fmt.Println(err)
	}
	s.StartAsync()

	r.Post("/win", srv.handleCheckIsTicketWinning)

	http.ListenAndServe(checkerHTTP, r)
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

func makeClient(uri string) *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	return client
}
