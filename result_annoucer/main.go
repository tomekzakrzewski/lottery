package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis"
	receiver "github.com/tomekzakrzewski/lottery/number_receiver/client"
	checker "github.com/tomekzakrzewski/lottery/result_checker/client"
	"github.com/tomekzakrzewski/lottery/types"
	"google.golang.org/grpc"
)

func main() {
	annoucerGRPC := "localhost:6006"
	//checkerClient := checker.NewHTTPClient("http://localhost:5000")
	checkerGRPCClient, _ := checker.NewGRPCClient("localhost:3009")
	//receiverClient := receiver.NewHTTPClient("http://localhost:3000")
	receiverGRPCClient, _ := receiver.NewGRPCClient("localhost:3006")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	redis := NewRedisStore(redisClient)
	svc := NewResultAnnoucerService(*checkerGRPCClient, *receiverGRPCClient, redis)
	m := NewLogMiddleware(svc)
	srv := NewHttpTransport(m)

	go func() {
		log.Fatal(makeGRPCTransport(annoucerGRPC, m))
	}()
	r := chi.NewRouter()
	r.Get("/win/{hash}", srv.handleCheckResult)

	http.ListenAndServe(":6000", r)

}

func makeGRPCTransport(listenAddr string, svc ResultAnnoucer) error {
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
	types.RegisterAnnoucerServer(grpcServer, NewAnnoucerGRPCServer(svc))
	return grpcServer.Serve(ln)
}
