package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis"
	receiver "github.com/tomekzakrzewski/lottery/number_receiver/client"
	checker "github.com/tomekzakrzewski/lottery/result_checker/client"
	"github.com/tomekzakrzewski/lottery/types"
	"google.golang.org/grpc"
)

func main() {

	var (
		redisUri     = os.Getenv("REDIS_URI")
		annoucerGRPC = os.Getenv("ANNOUNCER_GRPC")
		annoucerHTTP = os.Getenv("ANNOUNCER_HTTP")
		checkerGRPC  = os.Getenv("CHECKER_GRPC")
		receiverGRPC = os.Getenv("RECEIVER_GRPC")
		redisClient  = makeRedis(redisUri)
		redis        = NewRedisStore(redisClient)
		r            = chi.NewRouter()
	)
	checkerGRPCClient, _ := checker.NewGRPCClient(checkerGRPC)
	receiverGRPCClient, _ := receiver.NewGRPCClient(receiverGRPC)
	svc := NewResultAnnoucerService(checkerGRPCClient, receiverGRPCClient, redis)
	m := NewLogMiddleware(svc)
	srv := NewHttpTransport(m)

	go func() {
		log.Fatal(makeGRPCTransport(annoucerGRPC, m))
	}()
	r.Get("/win/{hash}", srv.handleCheckResult)

	http.ListenAndServe(annoucerHTTP, r)

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

func makeRedis(uri string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     uri,
		Password: "",
		DB:       0,
	})
}
