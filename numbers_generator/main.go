package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/tomekzakrzewski/lottery/number_receiver/client"
	"github.com/tomekzakrzewski/lottery/types"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	var (
		generatorHTTP = os.Getenv("GENERATOR_HTTP")
		generatorGRPC = os.Getenv("GENERATOR_GRPC")
		receiverGRPC  = os.Getenv("RECEIVER_GRPC")
		r             = chi.NewRouter()
	)
	grpcClient, err := client.NewGRPCClient(receiverGRPC)
	if err != nil {
		panic(err)
	}
	svc := NewGeneratorService(grpcClient)
	m := NewLogMiddleware(svc)

	go func() {
		log.Fatal(makeGRPCTransport(generatorGRPC, m))
	}()
	srv := NewHttpTransport(m)

	r.Get("/winningNumbers", srv.handleGetWinningNumbers)
	http.ListenAndServe(generatorHTTP, r)
}

func makeGRPCTransport(listenAddr string, svc GeneratorServicer) error {
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
	types.RegisterGeneratorServer(grpcServer, NewGeneratorGRPCServer(svc))
	return grpcServer.Serve(ln)
}
