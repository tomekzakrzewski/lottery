package main

import (
	"fmt"
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

	var (
		//redisUri     = os.Getenv("REDIS_URI")
		annoucerGRPC = ":6006"
		checkerGRPC  = ":3009"
		receiverGRPC = ":3006"
		r            = chi.NewRouter()
	)

	redis := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	store := NewRedisStore(redis)

	checkerClient, err := checker.NewGRPCClient(checkerGRPC)
	if err != nil {
		panic(err)
	}

	receiverClient, err := receiver.NewGRPCClient(receiverGRPC)
	if err != nil {
		panic(err)
	}
	svc := NewResultAnnoucerService(checkerClient, receiverClient, store)
	m := NewLogMiddleware(svc)

	go func() {
		panic(makeGRPCTransport(annoucerGRPC, m))
	}()
	srv := NewHttpTransport(m)

	r.Get("/win/{hash}", srv.handleCheckResult)
	http.ListenAndServe(":6001", r)
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
