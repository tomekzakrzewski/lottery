package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tomekzakrzewski/lottery/number_receiver/client"
	"github.com/tomekzakrzewski/lottery/types"
	"google.golang.org/grpc"
)

func main() {

	var (
		//receiverEndpoint = "http://localhost:3000"
		generatorGRCP = "localhost:3005"
		receiverGRPC  = "localhost:3006"
	)

	//httpClient := client.NewHTTPClient(receiverEndpoint)
	grpcClient, err := client.NewGRPCClient(receiverGRPC)
	if err != nil {
		panic(err)
	}
	svc := NewGeneratorService(*grpcClient)
	m := NewLogMiddleware(svc)

	go func() {
		log.Fatal(makeGRPCTransport(generatorGRCP, m))
	}()
	srv := NewHttpTransport(m)

	r := chi.NewRouter()
	r.Get("/winningNumbers", srv.handleGetWinningNumbers)
	http.ListenAndServe(":3001", r)
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
